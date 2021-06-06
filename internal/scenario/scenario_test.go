package scenario

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScenario_HandleSearch_BadCommand(t *testing.T) {
	testCtx := GetTestScenarioCtx()
	testCtx.Scenario.HandleSearch(&UserMessage{
		Body:   "!s",
		UserId: 1,
	})
	assert.Equal(t, []testMessage{
		{
			userId:  1,
			message: incorrectMessage,
		},
	}, testCtx.Sender.sent)
	assert.True(t, strings.Contains(testCtx.LogBuf.String(), "[info]"))
}

func TestScenario_HandleSearch_GoodCommand(t *testing.T) {
	testCtx := GetTestScenarioCtx()
	testCtx.Scenario.HandleSearch(&UserMessage{
		Body:   "!s grn 228",
		UserId: 1,
	})
	assert.Equal(t, []testMessage{
		{
			userId:  1,
			message: "good",
		},
	}, testCtx.Sender.sent)
}

func TestScenario_HandleSearch_NotFoundCard(t *testing.T) {
	testCtx := GetTestScenarioCtx()
	testCtx.Scenario.HandleSearch(&UserMessage{
		Body:   "absolutely_random_card",
		UserId: 1,
	})
	assert.Equal(t, []testMessage{
		{
			userId:  1,
			message: cardNotFoundMessage,
		},
	}, testCtx.Sender.sent)
	assert.True(t, strings.Contains(testCtx.LogBuf.String(), "[info]"))
}

func TestScenario_HandleSearch_BadCard(t *testing.T) {
	testCtx := GetTestScenarioCtx()
	testCtx.Scenario.HandleSearch(&UserMessage{
		Body:   "bad",
		UserId: 1,
	})
	assert.Equal(t, []testMessage{
		{
			userId:  1,
			message: pricesUnavailableMessage,
		},
	}, testCtx.Sender.sent)
	assert.True(t, strings.Contains(testCtx.LogBuf.String(), "[error]"))
}
func TestScenario_HandleSearch_Uncached(t *testing.T) {
	testCtx := GetTestScenarioCtx()
	testCtx.Scenario.HandleSearch(&UserMessage{
		Body:   "uncached",
		UserId: 1,
	})
	assert.Equal(t, []testMessage{
		{
			userId:  1,
			message: "uncached",
		},
	}, testCtx.Sender.sent)
	_, err := testCtx.Scenario.Cache.Get("uncached")
	assert.Nil(t, err)
}
