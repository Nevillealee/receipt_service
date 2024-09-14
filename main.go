// main.go
package main

import (
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/receipts/process", processReceiptHandler)
	router.HandleFunc("/receipts/", getPointsHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
