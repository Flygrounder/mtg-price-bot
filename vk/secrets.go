package vk

import (
	"os"
	"strconv"
)

var TOKEN = os.Getenv("VK_TOKEN")
var SECRET_KEY = os.Getenv("VK_SECRET_KEY")
var GROUPID, _ = strconv.ParseInt(os.Getenv("VK_GROUP_ID"), 10, 64)
var CONFIRMATION_STRING = os.Getenv("VK_CONFIRMATION_STRING")
