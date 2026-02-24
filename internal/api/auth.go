package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// AuthHeaders returns the authentication headers required by the SwitchBot API v1.1.
// It generates: Authorization, sign, t, nonce
func AuthHeaders(token, secret string) map[string]string {
	t := strconv.FormatInt(time.Now().UnixMilli(), 10)
	nonce := uuid.New().String()

	payload := token + t + nonce
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return map[string]string{
		"Authorization": token,
		"sign":          sign,
		"t":             t,
		"nonce":         nonce,
		"Content-Type":  "application/json; charset=utf8",
	}
}

// ErrAPI is returned when the API responds with a non-success status code.
type ErrAPI struct {
	StatusCode int
	Message    string
}

func (e *ErrAPI) Error() string {
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}
