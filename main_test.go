package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"receipt_service/internal"
	"testing"
)

var (
	receipts = make(map[string]*internal.Receipt)
	points   = make(map[string]int)
)

func TestProcessReceiptHandler_Success(t *testing.T) {
	receipt := internal.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []internal.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
			{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
			{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
		},
		Total: "35.35",
	}
	body, _ := json.Marshal(receipt)
	req, err := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(body))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(internal.ProcessReceiptHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &resp)

	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	if _, ok := resp["id"]; !ok {
		t.Errorf("Response does not contain id")
	}
}

func TestProcessReceiptHandler_BadRequest(t *testing.T) {
	req, err := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer([]byte("invalid json")))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(internal.ProcessReceiptHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestGetPointsHandler_Success(t *testing.T) {
	receipt := internal.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []internal.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
			{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
			{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
		},
		Total: "35.35",
	}
	id := internal.GenerateID()
	receipts[id] = &receipt

	pts := internal.CalculatePoints(&receipt)

	points[id] = pts

	req, err := http.NewRequest("GET", "/receipts/"+id+"/points", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(internal.GetPointsHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp map[string]int
	err = json.Unmarshal(rr.Body.Bytes(), &resp)

	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	if resp["points"] != pts {
		t.Errorf("Expected points %v, got %v", pts, resp["points"])
	}
}

func TestGetPointsHandler_NotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/receipts/invalid_id/points", nil)

	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(internal.GetPointsHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}
