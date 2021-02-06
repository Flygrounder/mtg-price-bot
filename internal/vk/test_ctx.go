package vk

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http/httptest"
)

type testCtx struct {
	handler  *Handler
	recorder *httptest.ResponseRecorder
	sender   *testSender
	logBuf   *bytes.Buffer
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

func getTestRequestCtx(msgReq *messageRequest, recorder *httptest.ResponseRecorder) *gin.Context {
	ctx, _ := gin.CreateTestContext(recorder)
	body, _ := json.Marshal(msgReq)
	ctx.Request = httptest.NewRequest("POST", "/", bytes.NewReader(body))
	return ctx
}
