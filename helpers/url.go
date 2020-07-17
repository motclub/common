package helpers

import (
	"fmt"
	"net/url"
)

func SetUrlQuery(rawUrl string, values map[string]interface{}, replace ...bool) string {
	u, _ := url.Parse(rawUrl)
	if len(replace) > 0 && replace[0] {
		u.RawQuery = ""
	}
	query := u.Query()
	for k, v := range values {
		query.Set(k, fmt.Sprintf("%v", v))
	}
	u.RawQuery = query.Encode()
	return u.String()
}
