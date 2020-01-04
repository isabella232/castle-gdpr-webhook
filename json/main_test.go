package main

import (
	"testing"
)

func Test_VerifyWebhookMAC(t *testing.T) {
	if verifyWebhookMAC("0123456789", "/L2TBNCvuNaUKLzH/6/A+NYVzPuhWv5jQpGEN219MUo=", "i'm a secret") != true {
		t.Errorf("verifyWebhookMAC failed to return true for valid hmac")
	}
	if verifyWebhookMAC("0000000000", "/L2TBNCvuNaUKLzH/6/A+NYVzPuhWv5jQpGEN219MUo=", "i'm a secret") != false {
		t.Errorf("verifyWebhookMAC failed to return false for valid hmac")
	}
}

var correctJson = `
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

var incorrectApiVersion = `
{
	"api_version": "v4",
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

var incorrectType = `
{
	"api_version": "v1",
	"app_id": "3823955555537961",
	"type": "hello",
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

var missingUrl = `
{
	"api_version": "v1",
	"app_id": "3823955555537961",
	"type": "$gdpr.subject_access_request.completed",
	"created_at": "2019-12-01T19:38:28.483Z",
	"data": {
		"id": "test",
		"download_url": "",
		"download_url_expires_at": "2020-12-12T00:00.00Z",
		"user_id": "2",
		"user_traits": {
			"id": "2",
			"email": "email@example.com"
		}
	}
}
`

func Test_HandleIncomingWebHookData(t *testing.T) {

	url, err := handleIncomingWebHookData("", "", "")
	if err.Error() != "lenght of jsonString is 0" {
		t.Errorf("handleIncomingWebHookData failed handle invalid data")
	}
	if len(url) != 0 {
		t.Errorf("handleIncomingWebHookData failed handle invalid data")
	}

	url, err = handleIncomingWebHookData(correctJson, "", "")
	if err.Error() != "castleSignature invalid" {
		t.Errorf("handleIncomingWebHookData failed handle invalid data")
	}
	if len(url) != 0 {
		t.Errorf("handleIncomingWebHookData failed handle invalid data")
	}

	url, err = handleIncomingWebHookData(correctJson, "0123", "")
	if err.Error() != "hmac key invalid" {
		t.Errorf("handleIncomingWebHookData failed handle invalid data")
	}
	if len(url) != 0 {
		t.Errorf("handleIncomingWebHookData failed handle invalid data")
	}

	url, err = handleIncomingWebHookData(correctJson, "DSgTD4h47B0IexVVMCcbyG1T80LDnFmgiYXWqEmu/gI=", "i'm a secret")
	if err != nil {
		t.Errorf("handleIncomingWebHookData failed handle return valid url: " + err.Error())
	}
	if url != "https://url/user.zip" {
		t.Errorf("handleIncomingWebHookData failed handle return valid url: " + url)
	}

	url, err = handleIncomingWebHookData(incorrectApiVersion, "7Yphb6NozJKjusGldNPcMqTr1/bfiCxaduIHTWFrcf8=", "i'm a secret")
	if err.Error() != "invalid API version: v4" {
		t.Errorf("handleIncomingWebHookData failed handle invalid data")
	}
	if len(url) != 0 {
		t.Errorf("handleIncomingWebHookData failed handle invalid data")
	}

	url, err = handleIncomingWebHookData(incorrectType, "0vef8YEJPUMgIBCHzqRx7y1fjM8hhPpI9YScECt4acM=", "i'm a secret")
	if err.Error() != "invalid type: hello" {
		t.Errorf("handleIncomingWebHookData failed handle invalid data")
	}
	if len(url) != 0 {
		t.Errorf("handleIncomingWebHookData failed handle invalid data")
	}

	url, err = handleIncomingWebHookData(missingUrl, "FHcazuyhgx0oNzQiG7L5f+G4/XeoqYnqABvf/AqfrmQ=", "i'm a secret")
	if err.Error() != "empty download url" {
		t.Errorf("handleIncomingWebHookData failed handle invalid data")
	}
	if len(url) != 0 {
		t.Errorf("handleIncomingWebHookData failed handle invalid data")
	}
}
