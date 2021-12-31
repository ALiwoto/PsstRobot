package utils

import (
	"strconv"
	"strings"

	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func GetName(user *gotgbot.User) string {
	if user.FirstName != "" {
		return FixName(user.FirstName)
	}

	if user.LastName != "" {
		return FixName(user.LastName)
	}

	if user.Username != "" {
		return FixName(user.Username)
	}

	return ""
}

func ExtractUsername(text string) string {
	if text[0] == '@' {
		return text
	}
	return ""
}

func ExtractUserId(text string) int64 {
	value, _ := strconv.ParseInt(text, 10, 64)
	return value
}

func ExtractUserIdFromMessage(message *gotgbot.Message, alt ...string) int64 {
	if len(message.Entities) > 0 {
		for _, entity := range message.Entities {
			if entity.Type == "text_mention" && entity.User != nil {
				return entity.User.Id
			}
		}
	}

	if len(alt) < 1 {
		return ExtractUserId(message.Text)
	}

	var id int64
	for _, current := range alt {
		id = ExtractUserId(current)
		if id > 0 {
			return id
		}
	}

	return 0
}

func FixName(name string) string {
	if len(name) > 20 {
		return name[:20]
	}

	return name
}

func ArgsToBool(args []string) (bool, bool) {
	if len(args) < 2 {
		return false, false
	}

	myStr := strings.ToLower(args[1])
	switch myStr {
	case "on", "true", "enable":
		return true, true
	case "off", "false", "disable":
		return false, true
	}

	return false, false
}

func GetDBIndex(id int64) int {
	return int(strconv.FormatInt(id, 10)[0] - '0')
}

func ToBase32(value int64) string {
	return strconv.FormatInt(value, 32)
}

func ToBase10(value int64) string {
	return strconv.FormatInt(value, 10)
}

// UnpackInlineMessageId returns unpacked inline message id result.
// all credits goes to https:/github.com/PaulSonOfLars for writing the logic.
// the original source code is here:
// https://gist.github.com/PaulSonOfLars/49039e6422f21d5f54d0fe06e18d5143
func UnpackInlineMessageId(inlineMessageId string) (*UnpackInlineMessageResult, error) {
	bs, err := base64.RawURLEncoding.DecodeString(inlineMessageId)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode inline message ID: %w", err)
	}

	// unsigned 32
	dcId, err := readULong(bs[0:4])
	if err != nil {
		return nil, fmt.Errorf("failed to get DC ID from inline message ID: %w", err)
	}

	// unsigned 32
	messageId, err := readULong(bs[4:8])
	if err != nil {
		return nil, fmt.Errorf("failed to get message ID from inline message ID: %w", err)
	}

	// Some odd logic is required to deal with 32/64 bit ids
	var chatId int64
	var queryId int64
	if len(bs) == 20 {
		// signed 32
		chatId, err = readLong(bs[8:12])
		if err != nil {
			return nil, fmt.Errorf("failed to get 32bit chat ID from inline message ID: %w", err)
		}
		// signed 64
		queryId, err = readLongLong(bs[12:20])
		if err != nil {
			return nil, fmt.Errorf("failed to get query ID from inline message ID: %w", err)
		}
	} else {
		// signed 64
		chatId, err = readLongLong(bs[8:16])
		if err != nil {
			return nil, fmt.Errorf("failed to get 64bit chat ID from inline message ID: %w", err)
		}
		// signed 64
		queryId, err = readLongLong(bs[16:24])
		if err != nil {
			return nil, fmt.Errorf("failed to get query ID from inline message ID: %w", err)
		}
	}

	if chatId < 0 {
		fullChatId := "-100" + strconv.FormatInt(-chatId, 10)
		chatId, err = strconv.ParseInt(fullChatId, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to fix chat ID to bot API: %w", err)
		}
	}

	return &UnpackInlineMessageResult{
		InlineMessageId: inlineMessageId,
		DC:              dcId,
		MessageID:       messageId,
		ChatID:          chatId,
		QueryID:         queryId,
	}, nil
}

func readULong(b []byte) (int64, error) {
	buf := bytes.NewBuffer(b)
	var x uint32
	err := binary.Read(buf, binary.LittleEndian, &x)
	if err != nil {
		return 0, err
	}
	return int64(x), nil
}

func readLong(b []byte) (int64, error) {
	buf := bytes.NewBuffer(b)
	var x int32
	err := binary.Read(buf, binary.LittleEndian, &x)
	if err != nil {
		return 0, err
	}
	return int64(x), nil
}

func readLongLong(b []byte) (int64, error) {
	buf := bytes.NewBuffer(b)
	var x int64
	err := binary.Read(buf, binary.LittleEndian, &x)
	if err != nil {
		return 0, err
	}
	return x, nil
}

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

	result.Text = value

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
