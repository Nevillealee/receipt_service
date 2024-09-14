package main

import (
	"log"
	"net/http"
	"receipt_service/internal"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/receipts/process", internal.ProcessReceiptHandler)
	router.HandleFunc("/receipts/", internal.GetPointsHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
