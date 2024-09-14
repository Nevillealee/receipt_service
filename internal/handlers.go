package internal

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func generateID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var receipt Receipt
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &receipt)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	id := generateID()

	mu.Lock()
	receipts[id] = &receipt
	mu.Unlock()

	pts := calculatePoints(&receipt)

	mu.Lock()
	points[id] = pts
	mu.Unlock()

	response := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) != 3 || pathParts[0] != "receipts" || pathParts[2] != "points" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	id := pathParts[1]

	mu.Lock()
	pts, exists := points[id]
	mu.Unlock()

	if !exists {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	response := map[string]int{"points": pts}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
