package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	sc "github.com/gorilla/securecookie"
)

func GenerateHMACSecret() []byte {
	return sc.GenerateRandomKey(64)
}

// Logic for this flow is from:
// https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#pseudo-code-for-implementing-hmac-csrf-tokens

func NewCSRFToken(sessionID string, HMACSecret []byte) string {
	randHexValue := hex.EncodeToString(sc.GenerateRandomKey(64))
	hmacHash := createHMACHash(sessionID, randHexValue, HMACSecret)
	return hex.EncodeToString(hmacHash) + "." + randHexValue
}

func ValidateCSRFToken(sessionID, csrfToken string, HMACSecret []byte) bool {
	parts := strings.Split(csrfToken, ".")
	if len(parts) != 2 {
		return false
	}

	hmacHex := parts[0]
	randomValueHex := parts[1]

	receivedHMAC, err := hex.DecodeString(hmacHex)
	if err != nil {
		// TODO:
	}
	expectedHMAC := createHMACHash(sessionID, randomValueHex, HMACSecret)

	return hmac.Equal(receivedHMAC, expectedHMAC)
}

func createHMACHash(sessionID, randHexValue string, HMACSecret []byte) []byte {
	message := fmt.Sprintf("%d!%s!%d!%s",
		len(sessionID),
		sessionID,
		len(randHexValue),
		randHexValue,
	)

	hmacHasher := hmac.New(sha256.New, HMACSecret)
	hmacHasher.Write([]byte(message))
	return hmacHasher.Sum(nil)
}
