package vk

import (
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const SendMessageUrl = "https://api.vk.com/method/messages.send"

func Message(userId int64, message string) {
	randomId := rand.Int31()
	params := []string{
		"access_token=" + Token,
		"peer_id=" + strconv.FormatInt(userId, 10),
		"message=" + url.QueryEscape(message),
		"v=5.95",
		"random_id=" + strconv.FormatInt(int64(randomId), 10),
	}
	paramString := strings.Join(params, "&")
	resp, err := http.Get(SendMessageUrl + "?" + paramString)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Print("Could not send message\n user: %lld", userId)
	}
}
