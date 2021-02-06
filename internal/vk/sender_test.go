package vk

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestApiSender_Send_OK(t *testing.T) {
	defer gock.Off()

	gock.New(SendMessageUrl).MatchParams(
		map[string]string{
			"access_token": "token",
			"peer_id":      "1",
			"message":      "msg",
			"v":            "5.95",
		},
	).ParamPresent("random_id").Reply(http.StatusOK)

	sender := ApiSender{Token: "token"}
	sender.Send(1, "msg")
	assert.False(t, gock.HasUnmatchedRequest())
}

func TestApiSender_Send_NotOK(t *testing.T) {
	defer gock.Off()

	gock.New(SendMessageUrl).Reply(http.StatusInternalServerError)

	b := &bytes.Buffer{}
	sender := ApiSender{
		Token:  "token",
		Logger: log.New(b, "", 0),
	}
	sender.Send(1, "msg")
	assert.True(t, strings.Contains(b.String(), "[error]"))
}

func TestApiSender_Send_ErrorCode(t *testing.T) {
	defer gock.Off()

	gock.New(SendMessageUrl).Reply(http.StatusOK).JSON(
		map[string]interface{}{
			"error": map[string]interface{}{
				"error_code": 100,
				"error_msg":  "bad user",
			},
		},
	)

	b := &bytes.Buffer{}
	sender := ApiSender{
		Token:  "token",
		Logger: log.New(b, "", 0),
	}
	sender.Send(1, "msg")
	assert.True(t, strings.Contains(b.String(), "[error]"))
}
