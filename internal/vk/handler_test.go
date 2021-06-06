package vk

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

