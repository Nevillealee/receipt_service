package internal

import "sync"

var (
	receipts = make(map[string]*Receipt)
	points   = make(map[string]int)
	mu       sync.Mutex
)
