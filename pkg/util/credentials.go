package util

import (
	b64 "encoding/base64"
	"strings"
)

// Base64Credentials returns the base64 encoded "username:password"
func Base64Credentials(username, password string) string {
	var sb strings.Builder

	sb.WriteString(username)
	sb.WriteRune(':')
	sb.WriteString(password)

	return b64.StdEncoding.EncodeToString([]byte(sb.String()))
}
