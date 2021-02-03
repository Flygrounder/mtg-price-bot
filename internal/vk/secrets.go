package vk

import (
	"os"
	"strconv"
)

var Token = os.Getenv("VK_TOKEN")
var SecretKey = os.Getenv("VK_SECRET_KEY")
var GroupId, _ = strconv.ParseInt(os.Getenv("VK_GROUP_ID"), 10, 64)
var ConfirmationString = os.Getenv("VK_CONFIRMATION_STRING")
