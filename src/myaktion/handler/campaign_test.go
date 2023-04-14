package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/turngeek/myaktion-go-2023/src/myaktion/db"
	"github.com/turngeek/myaktion-go-2023/src/myaktion/handler"
	"github.com/turngeek/myaktion-go-2023/src/myaktion/model"
)

func TestCreateCampaign(t *testing.T) {
	cleanUpDB := db.SetupTestDB(t)
	defer cleanUpDB()
	rr := httptest.NewRecorder()
	jsonData := `{"name":"Covid","organizerName":"Martin","donationMinimum":2,"targetAmount":100,"account":{"name":"Martin","bankName":"DKB","number":"123456"}}`
	req := httptest.NewRequest(http.MethodPost, "/dummy-url", bytes.NewBufferString(jsonData))
	req.Header.Set("Content-Type", "application/json")

	handler := http.HandlerFunc(handler.CreateCampaign)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var campaign model.Campaign
	err := json.NewDecoder(rr.Body).Decode(&campaign)
	if err != nil {
		t.Errorf("handler returned unexpected body: %v", err)
		return
	}
	if campaign.ID != 1 {
		t.Errorf("handler returned unexpected ID: got %v want %v", campaign.ID, 1)
	}
}
