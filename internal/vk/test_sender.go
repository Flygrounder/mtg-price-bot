package vk

type testSender struct {
	sent []testMessage
}

type testMessage struct {
	userId  int64
	message string
}

func (s *testSender) send(userId int64, message string) {
	s.sent = append(s.sent, testMessage{
		userId:  userId,
		message: message,
	})
}
