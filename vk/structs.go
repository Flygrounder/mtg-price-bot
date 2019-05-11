package vk

type MessageRequest struct {
	Object UserMessage `json:"object"`
	Secret string      `json:"secret"`
}

type UserMessage struct {
	Body   string `json:"body"`
	UserId int64  `json:"user_id"`
}
