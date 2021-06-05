package vk

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const sendMessageUrl = "https://api.vk.com/method/messages.send"

type sender interface {
	send(userId int64, message string)
}

type ApiSender struct {
	Token  string
	Logger *log.Logger
}

type sendMessageResponse struct {
	Error errorResponse `json:"error"`
}

type errorResponse struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

func (s *ApiSender) send(userId int64, message string) {
	randomId := rand.Int63()
	params := []string{
		"access_token=" + s.Token,
		"peer_id=" + strconv.FormatInt(userId, 10),
		"message=" + url.QueryEscape(message),
		"v=5.95",
		"random_id=" + strconv.FormatInt(randomId, 10),
		"dont_parse_links=1",
	}
	joined := strings.Join(params, "&")
	reqUrl := sendMessageUrl + "?" + joined
	resp, err := http.Get(reqUrl)
	if err != nil || resp.StatusCode != http.StatusOK {
		s.Logger.Printf("[error] Could not Send message. User: %d", userId)
		return
	}
	respContent, _ := ioutil.ReadAll(resp.Body)
	var unmarshalled sendMessageResponse
	_ = json.Unmarshal(respContent, &unmarshalled)
	if unmarshalled.Error.ErrorCode != 0 {
		s.Logger.Printf("[error] Message was not sent. User: %d error message: %s", userId, unmarshalled.Error.ErrorMsg)
	}
}
