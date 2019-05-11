package vk

type MessageRequest struct {
	Type    string      `json:"type"`
	GroupId int64       `json:"group_id"`
	Object  UserMessage `json:"object"`
	Secret  string      `json:"secret"`
}

type UserMessage struct {
	Body   string `json:"body"`
	UserId int64  `json:"user_id"`
}
