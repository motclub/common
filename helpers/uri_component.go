package helpers

import (
	"net/url"
	"strings"
)

func EncodeURIComponent(str string) string {
	r := url.QueryEscape(str)
	r = strings.Replace(r, "+", "%20", -1)
	return r
}

func DecodeURIComponent(str string) (string, error) {
	r, err := url.QueryUnescape(str)
	if err != nil {
		return "", err
	}
	return r, nil
}
