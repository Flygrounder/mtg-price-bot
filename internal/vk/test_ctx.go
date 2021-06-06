package vk

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"gitlab.com/flygrounder/go-mtg-vk/internal/scenario"
)

type testCtx struct {
	handler  *Handler
	recorder *httptest.ResponseRecorder
}

func getTestHandlerCtx() testCtx {
    s := scenario.GetTestScenarioCtx()
	return testCtx{
		handler: &Handler{
			SecretKey:          "sec",
			GroupId:            10,
			ConfirmationString: "con",
			Scenario: s.Scenario,
		},
		recorder: httptest.NewRecorder(),
	}
}

func getTestRequestCtx(msgReq *messageRequest, recorder *httptest.ResponseRecorder) *gin.Context {
	ctx, _ := gin.CreateTestContext(recorder)
	body, _ := json.Marshal(msgReq)
	ctx.Request = httptest.NewRequest("POST", "/", bytes.NewReader(body))
	return ctx
}
