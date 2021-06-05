package vk

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestApiSender_Send_OK(t *testing.T) {
	defer gock.Off()

	gock.New(sendMessageUrl).MatchParams(
		map[string]string{
			"access_token":     "token",
			"peer_id":          "1",
			"message":          "msg",
			"v":                "5.95",
			"dont_parse_links": "1",
		},
	).ParamPresent("random_id").Reply(http.StatusOK)

	sender := ApiSender{Token: "token"}
	sender.send(1, "msg")
	assert.False(t, gock.HasUnmatchedRequest())
}

func TestApiSender_Send_NotOK(t *testing.T) {
	defer gock.Off()

	gock.New(sendMessageUrl).Reply(http.StatusInternalServerError)

	b := &bytes.Buffer{}
	sender := ApiSender{
		Token:  "token",
		Logger: log.New(b, "", 0),
	}
	sender.send(1, "msg")
	assert.True(t, strings.Contains(b.String(), "[error]"))
}

func TestApiSender_Send_ErrorCode(t *testing.T) {
	defer gock.Off()

	gock.New(sendMessageUrl).Reply(http.StatusOK).JSON(
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
	sender.send(1, "msg")
	assert.True(t, strings.Contains(b.String(), "[error]"))
}
