package vk

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestHandler_HandleMessage_Confirm(t *testing.T) {
	testCtx := getTestHandlerCtx()
	ctx := getTestRequestCtx(&messageRequest{
		Type:    "confirmation",
		GroupId: testCtx.handler.GroupId,
		Secret:  testCtx.handler.SecretKey,
	}, testCtx.recorder)
	testCtx.handler.HandleMessage(ctx)
	assert.Equal(t, testCtx.handler.ConfirmationString, testCtx.recorder.Body.String())
}

func TestHandler_HandleMessage_Message(t *testing.T) {
	testCtx := getTestHandlerCtx()
	ctx := getTestRequestCtx(&messageRequest{
		Type:   "message_new",
		Secret: testCtx.handler.SecretKey,
	}, testCtx.recorder)
	testCtx.handler.HandleMessage(ctx)
	assert.Equal(t, "ok", testCtx.recorder.Body.String())
}

func TestHandler_HandleMessage_NoSecretKey(t *testing.T) {
	testCtx := getTestHandlerCtx()
	ctx := getTestRequestCtx(&messageRequest{
		Type: "message_new",
	}, testCtx.recorder)
	testCtx.handler.HandleMessage(ctx)
	assert.Equal(t, "", testCtx.recorder.Body.String())
}

func TestHandler_handleSearch_BadCommand(t *testing.T) {
	testCtx := getTestHandlerCtx()
	testCtx.handler.handleSearch(&messageRequest{
		Object: userMessage{
			Body:   "!s",
			UserId: 1,
		},
	})
	assert.Equal(t, []testMessage{
		{
			userId:  1,
			message: incorrectMessage,
		},
	}, testCtx.sender.sent)
	assert.True(t, strings.Contains(testCtx.logBuf.String(), "[info]"))
}

func TestHandler_handleSearch_GoodCommand(t *testing.T) {
	testCtx := getTestHandlerCtx()
	testCtx.handler.handleSearch(&messageRequest{
		Object: userMessage{
			Body:   "!s grn 228",
			UserId: 1,
		},
	})
	assert.Equal(t, []testMessage{
		{
			userId:  1,
			message: "good",
		},
	}, testCtx.sender.sent)
}

func TestHandler_handleSearch_NotFoundCard(t *testing.T) {
	testCtx := getTestHandlerCtx()
	testCtx.handler.handleSearch(&messageRequest{
		Object: userMessage{
			Body:   "absolutely_random_card",
			UserId: 1,
		},
	})
	assert.Equal(t, []testMessage{
		{
			userId:  1,
			message: cardNotFoundMessage,
		},
	}, testCtx.sender.sent)
	assert.True(t, strings.Contains(testCtx.logBuf.String(), "[info]"))
}

func TestHandler_handleSearch_BadCard(t *testing.T) {
	testCtx := getTestHandlerCtx()
	testCtx.handler.handleSearch(&messageRequest{
		Object: userMessage{
			Body:   "bad",
			UserId: 1,
		},
	})
	assert.Equal(t, []testMessage{
		{
			userId:  1,
			message: pricesUnavailableMessage,
		},
	}, testCtx.sender.sent)
	assert.True(t, strings.Contains(testCtx.logBuf.String(), "[error]"))
}
func TestHandler_handleSearch_Uncached(t *testing.T) {
	testCtx := getTestHandlerCtx()
	testCtx.handler.handleSearch(&messageRequest{
		Object: userMessage{
			Body:   "uncached",
			UserId: 1,
		},
	})
	assert.Equal(t, []testMessage{
		{
			userId:  1,
			message: "uncached",
		},
	}, testCtx.sender.sent)
	msg, _ := testCtx.handler.Cache.Get("uncached")
	assert.Equal(t, "uncached", msg)
}
