package html

import "strings"

func MakeSureHasDoctype(html []byte) []byte {
	trimStrHTML := strings.TrimSpace(string(html))
	html = []byte(trimStrHTML)

	if len(html) > 15 && strings.ToLower(string(html[:15])) == "<!doctype html>" {
		return html
	}
	return append([]byte("<!DOCTYPE html>"), html...)
}
