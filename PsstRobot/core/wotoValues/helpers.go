package wotoValues

import "strings"

func IsAdvancedInline(data string) bool {
	return strings.HasPrefix(data, AdvancedInlinePrefix) &&
		strings.HasSuffix(data, AdvancedInlineSuffix)
}

func GetUniqueIdFromAdvanced(query string) string {
	result := strings.TrimSuffix(query, AdvancedInlineSuffix)
	return strings.TrimPrefix(result, AdvancedInlinePrefix)
}
