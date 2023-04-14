package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/turngeek/myaktion-go-2023/src/myaktion/handler"
)

func TestHealth(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	handler := http.HandlerFunc(handler.Health)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
