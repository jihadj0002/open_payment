//go:build unit

package payment

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestProcessorServer(handler http.HandlerFunc) (*httptest.Server, *ProcessorClient) {
	server := httptest.NewServer(handler)
	client := NewProcessorClient(server.URL)
	return server, client
}

func TestProcessCard_Success(t *testing.T) {
	t.Parallel()

	server, client := newTestProcessorServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/process/card", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var req CardRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, int64(1000), req.Amount)
		assert.Equal(t, "USD", req.Currency)

		resp := ProcessorResponse{
			Success:      true,
			ProcessorRef: "proc_ref_123",
			Status:       "completed",
			Fee:          30,
			ProcessedAt:  "2025-01-01T00:00:00Z",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	req := CardRequest{
		Amount:   1000,
		Currency: "USD",
	}

	resp, err := client.ProcessCard(req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, "proc_ref_123", resp.ProcessorRef)
	assert.Equal(t, "completed", resp.Status)
	assert.Equal(t, int64(30), resp.Fee)
}

func TestProcessCard_Failure(t *testing.T) {
	t.Parallel()

	server, client := newTestProcessorServer(func(w http.ResponseWriter, r *http.Request) {
		resp := ProcessorResponse{
			Success: false,
			Status:  "failed",
			Message: "card declined by issuer",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	req := CardRequest{Amount: 500, Currency: "USD"}
	resp, err := client.ProcessCard(req)
	require.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "card declined by issuer", resp.Message)
}

func TestProcessCard_ServerError(t *testing.T) {
	t.Parallel()

	server, client := newTestProcessorServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal error"}`))
	})
	defer server.Close()

	req := CardRequest{Amount: 100, Currency: "USD"}
	resp, err := client.ProcessCard(req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Success)
}

func TestProcessCard_UnreachableServer(t *testing.T) {
	client := NewProcessorClient("http://127.0.0.1:1")
	req := CardRequest{Amount: 100, Currency: "USD"}

	_, err := client.ProcessCard(req)
	assert.Error(t, err)
}

func TestProcessCard_DefaultsApplied(t *testing.T) {
	t.Parallel()

	var capturedReq CardRequest
	server, client := newTestProcessorServer(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&capturedReq)
		resp := ProcessorResponse{Success: true, ProcessorRef: "ref_1", Status: "completed"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	req := CardRequest{Amount: 2000, Currency: "USD"}
	_, err := client.ProcessCard(req)
	require.NoError(t, err)
	assert.Equal(t, "4111111111111111", capturedReq.CardNumber)
	assert.Equal(t, "12", capturedReq.ExpiryMonth)
	assert.Equal(t, "2030", capturedReq.ExpiryYear)
	assert.Equal(t, "123", capturedReq.CVV)
}

func TestProcessCard_InvalidJSONResponse(t *testing.T) {
	t.Parallel()

	server, client := newTestProcessorServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{invalid json`))
	})
	defer server.Close()

	req := CardRequest{Amount: 100, Currency: "USD"}
	_, err := client.ProcessCard(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "decode response")
}
