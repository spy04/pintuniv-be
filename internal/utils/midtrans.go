package utils

import (
	"crypto/sha512"
	"encoding/hex"
)

func VerifyMidtransSignature(
	orderID string,
	statusCode string,
	grossAmount string,
	serverKey string,
	signatureKey string,
) bool {

	raw := orderID + statusCode + grossAmount + serverKey

	hash := sha512.Sum512([]byte(raw))
	expected := hex.EncodeToString(hash[:])

	return expected == signatureKey
}
