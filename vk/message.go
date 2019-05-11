package vk

import (
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const VKURL = "https://api.vk.com/method/messages.send"

func Message(userId int64, message string) {
	randomId := rand.Int31()
	params := []string{
		"access_token=" + TOKEN,
		"peer_id=" + strconv.FormatInt(userId, 10),
		"message=" + url.QueryEscape(message),
		"v=5.95",
		"random_id=" + strconv.FormatInt(int64(randomId), 10),
	}
	paramString := strings.Join(params, "&")
	http.Get(VKURL + "?" + paramString)
}
