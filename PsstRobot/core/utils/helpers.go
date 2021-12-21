package utils

import (
	"strconv"
	"strings"
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
		result.Text = strings.TrimPrefix(result.Text, myStrs[0])
		return result
	}

	id, err = strconv.ParseInt(myStrs[last], 10, 64)
	if err == nil {
		result.TargetID = id
		result.Text = strings.TrimSuffix(result.Text, myStrs[last])
		return result
	}

	return result
}
