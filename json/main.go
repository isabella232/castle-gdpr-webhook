package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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

var validJson = `
{
	"api_version": "v1",
	"app_id": "3823955555537961",
	"type": "$gdpr.subject_access_request.completed",
	"created_at": "2019-12-01T19:38:28.483Z",
	"data": {
		"id": "test",
		"download_url": "https://url/user.zip",
		"download_url_expires_at": "2020-12-12T00:00.00Z",
		"user_id": "2",
		"user_traits": {
			"id": "2",
			"email": "email@example.com"
		}
	}
}
`

// verifies the webhook by computing the HMAC SHA256 signature of the payload
// message and key are both string however messageMAC is a base64 encoded string
// see castle docs for explanation
func verifyWebhookMAC(message, messageMAC, key string) bool {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	expectedMAC := mac.Sum(nil)

	// the messageMAC is base64 encoded to we have to decode itours as well
	messageMACbytes, err := base64.StdEncoding.DecodeString(messageMAC)
	if err != nil {
		fmt.Println("error:", err)
		return false
	}
	fmt.Printf("HMac of %sexpected: %s computed: %s\n", message, string(expectedMAC), string(messageMACbytes))
	return hmac.Equal(messageMACbytes, expectedMAC)
}

// processes the incoming webhook data, checks the signature and the api conforms to the expected format
// is successful returns the url of the GDPR SAR
func handleIncomingWebHookData(jsonString, castleSignature, key string) (string, error) {
	if len(jsonString) == 0 {
		return "", errors.New("lenght of jsonString is 0")
	}
	if len(castleSignature) == 0 {
		return "", errors.New("castleSignature invalid")
	}
	if len(key) == 0 {
		return "", errors.New("hmac key invalid")
	}

	// first check the signature
	if verifyWebhookMAC(jsonString, castleSignature, key) == false {
		return "", errors.New("hmac invalid")
	}

	b := []byte(jsonString)
	var sar GdprSar
	err := json.Unmarshal(b, &sar)
	if err != nil {
		return "", err
	}

	fmt.Printf("%+v\n", sar)

	if sar.ApiVersion != "v1" {
		return "", errors.New("invalid API version: " + sar.ApiVersion)
	}

	if sar.Type != "$gdpr.subject_access_request.completed" {
		return "", errors.New("invalid type: " + sar.Type)
	}

	if len(sar.Data.DownloadUrl) == 0 {
		return "", errors.New("empty download url")
	}
	return sar.Data.DownloadUrl, nil
}

func main() {
	b := []byte(validJson)
	var s GdprSar
	err := json.Unmarshal(b, &s)
	if err != nil {
		fmt.Printf("internal error, validJson not valid: %s", err.Error())
	}

	fmt.Printf("%+v\n", s)

	/*
		n := NewConsole()
		logger := log.New(os.Stdout, "[CSPReport]", log.LstdFlags)
		e := n.Process(r, nil, logger)
		if e != nil {
			t.Fatalf("Console failed print valid Json: %s", err.Error())
		}
	*/
}
