package url

import (
	"net/url"
)

func IsAbsolute(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}
