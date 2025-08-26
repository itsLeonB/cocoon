package util

import (
	"regexp"
	"strings"

	"github.com/itsLeonB/ezutil"
)

func GetNameFromEmail(email string) string {
	parts := strings.SplitN(email, "@", 2)
	if len(parts) < 2 || parts[0] == "" {
		return ""
	}
	localPart := parts[0]

	re := regexp.MustCompile(`[a-zA-Z]+`)
	matches := re.FindAllString(localPart, -1)
	if len(matches) > 0 {
		name := matches[0]
		return ezutil.Capitalize(name)
	}

	return ""
}
