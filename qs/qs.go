package qs

import (
	"fmt"
	"github.com/motclub/common/json"
	"net/url"
	"strings"
)

func Parse(query string) (map[string]string, error) {
	ret := make(map[string]string)
	if !(strings.HasPrefix(query, "/") || strings.HasPrefix(query, "http") || strings.HasPrefix(query, "?")) {
		query = fmt.Sprintf("?%s", query)
	}
	u, err := url.Parse("/?a=1&b=2")
	if err != nil {
		return ret, err
	}
	q := u.Query()
	for k, v := range q {
		if len(v) == 0 {
			continue
		}
		ret[k] = v[0]
	}
	return ret, nil
}

func MustParse(v string) map[string]string {
	ret, _ := Parse(v)
	return ret
}

func Stringify(v interface{}) (string, error) {
	var data map[string]interface{}
	dat, err := json.STD().Marshal(v)
	if err != nil {
		return "", err
	}
	if err := json.STD().Unmarshal(dat, &data); err != nil {
		return "", err
	}
	query := make(url.Values)
	for k, v := range data {
		query.Set(k, fmt.Sprintf("%v", v))
	}
	return query.Encode(), nil
}

func MustStringify(v interface{}) string {
	ret, _ := Stringify(v)
	return ret
}
