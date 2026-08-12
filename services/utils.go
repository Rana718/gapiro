package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// generateID creates a unique ID with the given prefix (e.g., "rq", "rs", "fl", "ws")
func generateID(prefix string) string {
	b := make([]byte, 8)
	rand.Read(b)
	ts := time.Now().UnixMilli() % 1000000
	return fmt.Sprintf("%s_%06d%s", prefix, ts, hex.EncodeToString(b))
}
