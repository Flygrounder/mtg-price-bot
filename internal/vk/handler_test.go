package vk

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gitlab.com/flygrounder/go-mtg-vk/internal/cardsinfo"
	"io"
	"log"
	"net/http/httptest"
	"strings"
	"testing"
)

func getTestRequestCtx(msgReq *messageRequest, recorder *httptest.ResponseRecorder) *gin.Context {
	ctx, _ := gin.CreateTestContext(recorder)
	body, _ := json.Marshal(msgReq)
	ctx.Request = httptest.NewRequest("POST", "/", bytes.NewReader(body))
	return ctx
}

type testCtx struct {
	handler  *Handler
	recorder *httptest.ResponseRecorder
	sender   *testSender
	logBuf   *bytes.Buffer
}

type testMessage struct {
	userId  int64
	message string
}

type testSender struct {
	sent []testMessage
}

func (s *testSender) Send(userId int64, message string) {
	s.sent = append(s.sent, testMessage{
		userId:  userId,
		message: message,
	})
}

type testCache struct {
	table map[string]string
}

func (t *testCache) Get(cardName string) (string, error) {
	msg, ok := t.table[cardName]
	if !ok {
		return "", errors.New("test")
	}
	return msg, nil
}

func (t *testCache) Set(cardName string, message string) {
	t.table[cardName] = message
}

func getTestHandlerCtx() testCtx {
	sender := &testSender{}
	buf := &bytes.Buffer{}
	return testCtx{
		logBuf: buf,
		handler: &Handler{
			SecretKey:          "sec",
			GroupId:            10,
			ConfirmationString: "con",
			Sender:             sender,
			Logger:             log.New(buf, "", 0),
			InfoFetcher:        &testInfoFetcher{},
			Cache: &testCache{
				table: map[string]string{
					"good": "good",
				},
			},
		},
		sender:   sender,
		recorder: httptest.NewRecorder(),
	}
}

type testInfoFetcher struct{}

func (t *testInfoFetcher) GetPrices(name string) ([]cardsinfo.CardPrice, error) {
	if name == "good" || name == "uncached" {
		return nil, nil
	}
	return nil, errors.New("test")
}

func (t *testInfoFetcher) FormatCardPrices(name string, _ []cardsinfo.CardPrice) string {
	return name
}

func (t *testInfoFetcher) GetNameByCardId(_ string, _ string) string {
	return "good"
}

func (t *testInfoFetcher) GetOriginalName(name string, _ io.Reader) string {
	if name == "good" || name == "bad" || name == "uncached" {
		return name
	}
	return ""
}

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
