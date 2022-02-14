package utils

type ExtractedResult struct {
	Text     string
	TargetID int64
	Username string
}

type UnpackInlineMessageResult struct {
	InlineMessageId string
	DC              int64
	MessageID       int64
	ChatID          int64
	QueryID         int64
}
