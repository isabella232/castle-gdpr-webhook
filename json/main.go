package main

import (
	"encoding/json"
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
