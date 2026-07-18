package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/dracory/req"
)

// GenerateCartCacheKey generates a cache key for guest cart based on IP and user agent
func GenerateCartCacheKey(r *http.Request) string {
	ip := req.GetIP(r)
	userAgent := r.UserAgent()

	// Hash the user agent for a shorter key
	hash := sha256.Sum256([]byte(userAgent))
	userAgentHash := hex.EncodeToString(hash[:16]) // Use 16 bytes for better collision resistance

	return fmt.Sprintf("cart_%s_%s", ip, userAgentHash)
}
