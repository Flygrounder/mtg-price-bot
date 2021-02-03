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

const SendMessageUrl = "https://api.vk.com/method/messages.send"

func Message(userId int64, message string) {
	randomId := rand.Int63()
	params := []string{
		"access_token=" + Token,
		"peer_id=" + strconv.FormatInt(userId, 10),
		"message=" + url.QueryEscape(message),
		"v=5.95",
		"random_id=" + strconv.FormatInt(randomId, 10),
	}
	paramString := strings.Join(params, "&")
	resp, err := http.Get(SendMessageUrl + "?" + paramString)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("[error] Could not send message. User: %d", userId)
		return
	}
	responseBytes, _ := ioutil.ReadAll(resp.Body)
	var response SendMessageResponse
	_ = json.Unmarshal(responseBytes, &response)
	if response.Error.ErrorCode != 0 {
		log.Printf("[error] Message was not sent. User: %d error message: %s", userId, response.Error.ErrorMsg)
	}
}
