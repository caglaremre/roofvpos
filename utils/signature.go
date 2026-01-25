package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"
	"roof/vpos/repository"
	"time"

	"github.com/google/uuid"
)

func CalculateSignature(body string, bolt *repository.Bolt) http.Header {
	nonce, err := uuid.NewV7()
	if err != nil {
		log.Panicf("Couldn't create uuid: %s", err)
	}

	clientToken, secretKey := bolt.ConfigRepo.GetClientAndSecretKey()
	if clientToken == "" || secretKey == "" {
		return http.Header{}
	}

	clientTokenDigest := sha256.Sum256([]byte(clientToken))
	clientTokenDigestEncoded := base64.StdEncoding.EncodeToString(clientTokenDigest[:])

	timestamp := time.Now().UTC().Format("20060102150405")

	signatureText := clientTokenDigestEncoded + secretKey + nonce.String() + timestamp + body
	signatureTextDigest := sha256.Sum256([]byte(signatureText))
	signatureTextDigestEncoded := base64.StdEncoding.EncodeToString(signatureTextDigest[:])

	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	headers.Add("x_client_token", clientToken)
	headers.Add("x_nonce", nonce.String())
	headers.Add("x_timestamp", timestamp)
	headers.Add("x_signature", signatureTextDigestEncoded)
	headers.Add("x_scenario", "mock")
	return headers

}
