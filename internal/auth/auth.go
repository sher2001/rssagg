package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extracts api key from the headers of the HTTP request
// Example :
// Authorization: ApiKey {insert API key here}
func GetAPIKey(headers http.Header) (string, error) {
	value := headers.Get("Authorization")
	if value == "" {
		return "", errors.New("no Autherization token found")
	}

	vals := strings.Split(value, " ")
	if len(vals) != 2 {
		return "", errors.New("autherization token is not formed correctly")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("autherization token type is not found")
	}

	return vals[1], nil
}
