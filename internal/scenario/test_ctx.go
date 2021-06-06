package scenario

import (
	"bytes"
	"log"

	"gitlab.com/flygrounder/go-mtg-vk/internal/cardsinfo"
)

type TestScenarioCtx struct {
	Scenario *Scenario
	Sender   *testSender
	LogBuf   *bytes.Buffer
}

func GetTestScenarioCtx() TestScenarioCtx {
	sender := &testSender{}
	buf := &bytes.Buffer{}
	return TestScenarioCtx{
		LogBuf: buf,
		Scenario: &Scenario{
			Sender:      sender,
			Logger:      log.New(buf, "", 0),
			InfoFetcher: &testInfoFetcher{},
			Cache: &testCache{
				table: map[string][]cardsinfo.ScgCardPrice{
					"good": {
						{
							Price:   "1",
							Edition: "alpha",
							Link:    "scg",
						},
					},
				},
			},
		},
		Sender: sender,
	}
}
