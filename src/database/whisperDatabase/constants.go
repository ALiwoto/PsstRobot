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
	WhisperTypeUnknown WhisperType = iota
	WhisperTypePlainText
	WhisperTypePhoto
	WhisperTypeVideo
	WhisperTypeAudio
	WhisperTypeVoice
	WhisperTypeSticker
	WhisperTypeDocument
	WhisperTypeVideoNote
	WhisperTypeAnimation
	WhisperTypeDice
)
