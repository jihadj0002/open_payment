package testutil

import (
	"net/http"
	"net/http/httptest"
)

func NewRequest(method, target string, body []byte) *http.Request {
	req := httptest.NewRequest(method, target, nil)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func ExecuteRequest(req *http.Request, handler http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}
