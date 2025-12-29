package checktgauth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
)

func New(botToken string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			r.Body.Close()
			var data map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &data); err != nil {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			if !checkTgHash(data, botToken) {
				http.Error(w, "wrong hash", http.StatusUnauthorized)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func checkTgHash(req map[string]interface{}, botToken string) bool {
	secretKey := sha256.Sum256([]byte(botToken))
	hash := req["hash"].(string)
	delete(req, "hash")
	var parts []string
	for k, v := range req {
		switch v.(type) {
		case float64:
			parts = append(parts, fmt.Sprintf("%s=%.0f", k, v))
		default:
			parts = append(parts, fmt.Sprintf("%s=%v", k, v))
		}

	}
	sort.Strings(parts)
	log.Println(parts)
	dataCheckString := strings.Join(parts, "\n")
	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(dataCheckString))
	calculatedHash := hex.EncodeToString(h.Sum(nil))
	return calculatedHash == hash
}
