package utils

func (e *ExtractedResult) IsForEveryone() bool {
	return e.TargetID == 0 && e.Username == ""
}

func (e *ExtractedResult) IsUsernameValid() bool {
	return e.Username == "" || len(e.Username) > MinUsernameLength
}
