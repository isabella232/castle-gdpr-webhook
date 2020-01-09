package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	/*
		"log"
		"os"
		"testing"
	*/)

type GdprSar struct {
	ApiVersion string `json:"api_version"`
	AppId      string `json:"app_id"`
	Type       string `json:"type"`
	CreatedAt  string `json:"created_at"`
	Data       Data   `json:"data"`
}

type Data struct {
	Id                   string     `json:"id"`
	DownloadUrl          string     `json:"download_url"`
	DownloadUrlExpiresAt string     `json:"download_url_expires_at"`
	UserId               string     `json:"user_id"`
	UserTraits           UserTraits `json:"user_traits"`
}

type UserTraits struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

// verifies the webhook by computing the HMAC SHA256 signature of the payload
// message and key are both strings however messageMAC is a base64 encoded string
// see castle docs for explanation
func verifyWebhookMAC(message, messageMACBase64, key string) bool {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	computedMAC := mac.Sum(nil)

	// the messageMACBase64 is base64 encoded to we have to decode it
	messageMAC, err := base64.StdEncoding.DecodeString(messageMACBase64)
	if err != nil {
		fmt.Println("error:", err)
		return false
	}
	messageMACB64 := base64.StdEncoding.EncodeToString(computedMAC)
	log.Printf("HMac of %s expected: %s computed: %s\n", message, messageMACBase64, messageMACB64)
	return hmac.Equal(messageMAC, computedMAC)
}

// processes the incoming webhook data, checks the signature and the api conforms to the expected format
// is successful returns the url of the GDPR SAR and the user_id
func HandleIncomingWebHookData(jsonString, castleSignature, key string) (string, string, error) {
	verifySignature := true // TODO: figure out how to pass custom header via API Gateway
	if len(jsonString) == 0 {
		return "", "", errors.New("lenght of jsonString is 0")
	}
	if verifySignature {
		if len(castleSignature) == 0 {
			return "", "", errors.New("castleSignature invalid")
		}
	}
	if len(key) == 0 {
		return "", "", errors.New("hmac key invalid")
	}

	if verifySignature {

		// first check the signature
		if verifyWebhookMAC(jsonString, castleSignature, key) == false {
			return "", "", errors.New("hmac invalid")
		}
	}

	b := []byte(jsonString)
	var sar GdprSar
	err := json.Unmarshal(b, &sar)
	if err != nil {
		return "", "", err
	}

	fmt.Printf("%+v\n", sar)

	if sar.ApiVersion != "v1" {
		return "", "", errors.New("invalid API version: " + sar.ApiVersion)
	}

	if sar.Type != "$gdpr.subject_access_request.completed" {
		return "", "", errors.New("invalid type: " + sar.Type)
	}

	if len(sar.Data.DownloadUrl) == 0 {
		return "", "", errors.New("empty download url")
	}
	return sar.Data.DownloadUrl, sar.Data.UserId, nil
}
