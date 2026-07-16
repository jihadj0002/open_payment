package settlement

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/openpayment/gateway/internal/api"
	"github.com/openpayment/gateway/internal/service/auth"
)

func HandleTriggerSettlement(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		var req SettlementRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "invalid_body", "invalid request body")
			return
		}

		settlement, err := svc.TriggerSettlement(r.Context(), claims.MerchantID, req)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, "settlement_error", "settlement_failed", err.Error())
			return
		}

		api.RespondJSON(w, http.StatusCreated, settlement)
	}
}

func HandleGetSettlement(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			api.RespondError(w, http.StatusBadRequest, "invalid_request", "missing_id", "settlement ID is required")
			return
		}

		settlement, err := svc.GetSettlement(r.Context(), id, claims.MerchantID)
		if err != nil {
			api.RespondError(w, http.StatusNotFound, "not_found", "settlement_not_found", "settlement not found")
			return
		}

		api.RespondJSON(w, http.StatusOK, settlement)
	}
}

func HandleListSettlements(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}

		perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
		if perPage < 1 || perPage > 100 {
			perPage = 20
		}

		settlements, total, err := svc.ListSettlements(r.Context(), claims.MerchantID, page, perPage)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "an unexpected error occurred")
			return
		}

		data := make([]interface{}, len(settlements))
		for i, s := range settlements {
			data[i] = s
		}

		api.RespondPaginated(w, data, total, page, perPage)
	}
}

func HandleSettlementReport(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r.Context())
		if claims == nil {
			api.RespondError(w, http.StatusUnauthorized, "auth_error", "unauthorized", "not authenticated")
			return
		}

		fromStr := r.URL.Query().Get("from")
		toStr := r.URL.Query().Get("to")
		currency := r.URL.Query().Get("currency")

		from := time.Now().AddDate(0, -1, 0)
		to := time.Now()

		if fromStr != "" {
			if t, err := time.Parse("2006-01-02", fromStr); err == nil {
				from = t
			}
		}
		if toStr != "" {
			if t, err := time.Parse("2006-01-02", toStr); err == nil {
				to = t
			}
		}

		merchantID := claims.MerchantID
		if claims.Role == "admin" {
			if mid := r.URL.Query().Get("merchant_id"); mid != "" {
				merchantID = mid
			}
		}

		report, err := svc.GetSettlementReport(r.Context(), merchantID, currency, from, to)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "server_error", "internal_error", "failed to generate report")
			return
		}

		accept := r.Header.Get("Accept")
		format := r.URL.Query().Get("format")

		if accept == "text/csv" || format == "csv" {
			w.Header().Set("Content-Type", "text/csv")
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=settlement-report-%s.csv", time.Now().Format("2006-01-02")))
			writer := csv.NewWriter(w)
			writer.Write([]string{"Currency", "Total Volume", "Total Fees", "Total Net", "Transaction Count"})
			for _, item := range report.Items {
				writer.Write([]string{
					item.Currency,
					strconv.FormatInt(item.TotalVolume, 10),
					strconv.FormatInt(item.TotalFees, 10),
					strconv.FormatInt(item.TotalNet, 10),
					strconv.Itoa(item.TransactionCount),
				})
			}
			writer.Write([]string{
				"TOTAL",
				strconv.FormatInt(report.Totals.TotalVolume, 10),
				strconv.FormatInt(report.Totals.TotalFees, 10),
				strconv.FormatInt(report.Totals.TotalNet, 10),
				strconv.Itoa(report.Totals.TransactionCount),
			})
			writer.Flush()
			return
		}

		api.RespondJSON(w, http.StatusOK, report)
	}
}
