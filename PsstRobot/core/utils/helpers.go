package utils

import (
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func ExtractRecipient(value string) *ExtractedResult {
	if len(value) == 0 {
		return nil
	}

	result := new(ExtractedResult)

	// supported formats:
	// @username message
	// message @username
	// ID message
	// message ID
	myStrs := strings.Fields(value)
	if myStrs[0][0] == '@' {
		result.Username = myStrs[0]
		result.Text = strings.TrimPrefix(value, myStrs[0])
		return result
	}

	last := len(myStrs) - 1
	if myStrs[last][0] == '@' {
		result.Username = myStrs[last]
		result.Text = strings.TrimSuffix(value, result.Username)
		return result
	}

	id, err := strconv.ParseInt(myStrs[0], 10, 64)
	if err == nil && id > MinUserID {
		result.TargetID = id
		result.Text = strings.TrimPrefix(value, myStrs[0])
		return result
	}

	id, err = strconv.ParseInt(myStrs[last], 10, 64)
	if err == nil {
		result.TargetID = id
		result.Text = strings.TrimSuffix(value, myStrs[last])
		return result
	}

	return result
}

// IsStartedForWhisper util function returns true if the text is
// for getting whisper text. basically you can use this function
// to determine if user has started bot to get their whisper text.
func IsStartedForWhisper(text string) bool {
	// we are sure that a start command with whisper data is
	// something like this:
	// "/start a=b"
	// so it's basically impossible to consider the text a whisper
	// with callback data if the length of the text is less than 10.
	return len(text) > 10 && strings.Contains(text, "=")
}

func GetWhisperUniqueId(message *gotgbot.Message) string {
	myStrs := strings.Split(message.Text, " ")
	if len(myStrs) < 2 {
		return ""
	}

	return myStrs[1]
}
