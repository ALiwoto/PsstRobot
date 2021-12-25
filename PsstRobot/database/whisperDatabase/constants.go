package whisperDatabase

const (
	MaxTextLength = 199
)

const (
	MediaGroupSep = "\u200D\u200F"
	CaptionSep    = "\u200D\u200F"
)

// whisper message types
const (
	WhisperTypePlainText WhisperType = iota
	WhisperTypePhoto
	WhisperTypeVideo
	WhisperTypeAudio
	WhisperTypeSticker
	WhisperTypeDocument
	WhisperTypeAnimation
)
