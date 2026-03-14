package middleware

import "strings"

func isBot(ua string, patterns []string) bool {
	ua = strings.ToLower(ua)

	for _, bot := range patterns {
		if strings.Contains(ua, strings.ToLower(bot)) {
			return true
		}
	}

	return false
}
