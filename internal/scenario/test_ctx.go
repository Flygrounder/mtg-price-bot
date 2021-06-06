package scenario

import (
	"bytes"
	"log"
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
				table: map[string]string{
					"good": "good",
				},
			},
		},
		Sender: sender,
	}
}
