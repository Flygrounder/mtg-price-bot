package scenario

import "gitlab.com/flygrounder/go-mtg-vk/internal/cardsinfo"

type testSender struct {
	sent []testMessage
}

func (s *testSender) SendPrices(userId int64, cardName string, prices []cardsinfo.ScgCardPrice) {
	s.sent = append(s.sent, testMessage{
		userId:  userId,
		message: cardName,
	})
}

type testMessage struct {
	userId  int64
	message string
}

func (s *testSender) Send(userId int64, message string) {
	s.sent = append(s.sent, testMessage{
		userId:  userId,
		message: message,
	})
}
