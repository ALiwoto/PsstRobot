package utils

func (e *ExtractedResult) IsForEveryone() bool {
	return e.TargetID == 0 && e.Username == ""
}
