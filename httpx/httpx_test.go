package httpx

import (
	"fmt"
	"github.com/motclub/common/json"
	"github.com/motclub/common/parser"
	"github.com/motclub/common/std"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func writeData(data std.D) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		defer func() { _ = r.Body.Close() }()
		buf, _ := ioutil.ReadAll(r.Body)
		body := make(std.D)
		_ = json.STD().Unmarshal(buf, &body)
		for k, v := range data {
			body[k] = v
		}
		_, _ = io.WriteString(w, json.Stringify(map[string]interface{}{
			"header": r.Header,
			"method": r.Method,
			"query":  r.URL.Query(),
			"body":   body,
			"form":   r.Form,
		}, false))
	})
}

type reply struct {
	Header http.Header `json:"header"`
	Method string      `json:"method"`
	Query  url.Values  `json:"query"`
	Body   std.D       `json:"body"`
	Form   url.Values  `json:"form"`
}

func TestGET(t *testing.T) {
	ts := httptest.NewServer(writeData(std.D{
		"login": "John Doe",
		"roles": []string{"moderator"},
		"id":    "b06cd03f-75d0-413a-b94b-35e155444d70",
	}))
	var dst reply
	if err := GET(fmt.Sprintf("%s?a=1&b=2", ts.URL), &dst); err != nil {
		t.Fatal(err)
	}
	body, _ := parser.NewParser(dst.Body)
	assert.Equal(t, http.MethodGet, dst.Method)
	assert.Equal(t, "application/json", dst.Header.Get("Content-Type"))
	assert.Equal(t, "1", dst.Query.Get("a"))
	assert.Equal(t, "2", dst.Query.Get("b"))
	assert.Equal(t, "b06cd03f-75d0-413a-b94b-35e155444d70", body.GetString("id"))
	assert.Equal(t, "John Doe", body.GetString("login"))
	assert.Equal(t, "moderator", body.GetString("roles[0]"))
}

func TestPOST(t *testing.T) {
	ts := httptest.NewServer(writeData(std.D{
		"login": "John Doe",
		"roles": []string{"moderator"},
		"id":    "b06cd03f-75d0-413a-b94b-35e155444d70",
	}))
	var dst reply
	if err := POST(fmt.Sprintf("%s?a=1&b=2", ts.URL), std.D{"username": "chan", "gender": "female", "age": 30}, &dst); err != nil {
		t.Fatal(err)
	}
	body, _ := parser.NewParser(dst.Body)
	assert.Equal(t, http.MethodPost, dst.Method)
	assert.Equal(t, "application/json", dst.Header.Get("Content-Type"))
	assert.Equal(t, "1", dst.Query.Get("a"))
	assert.Equal(t, "2", dst.Query.Get("b"))
	assert.Equal(t, "b06cd03f-75d0-413a-b94b-35e155444d70", body.GetString("id"))
	assert.Equal(t, "John Doe", body.GetString("login"))
	assert.Equal(t, "moderator", body.GetString("roles[0]"))
	assert.Equal(t, "chan", body.GetString("username"))
	assert.Equal(t, 30, body.GetInt("age"))
}

func TestPOSTForm(t *testing.T) {
	ts := httptest.NewServer(writeData(std.D{
		"login": "John Doe",
		"roles": []string{"moderator"},
		"id":    "b06cd03f-75d0-413a-b94b-35e155444d70",
	}))
	var dst reply
	if err := POSTForm(fmt.Sprintf("%s?a=1&b=2", ts.URL), std.D{"username": "chan", "gender": "female", "age": 30}, &dst); err != nil {
		t.Fatal(err)
	}
	body, _ := parser.NewParser(dst.Body)
	assert.Equal(t, http.MethodPost, dst.Method)
	assert.Equal(t, "application/x-www-form-urlencoded", dst.Header.Get("Content-Type"))
	assert.Equal(t, "1", dst.Query.Get("a"))
	assert.Equal(t, "2", dst.Query.Get("b"))
	assert.Equal(t, "b06cd03f-75d0-413a-b94b-35e155444d70", body.GetString("id"))
	assert.Equal(t, "John Doe", body.GetString("login"))
	assert.Equal(t, "moderator", body.GetString("roles[0]"))
	assert.Equal(t, "chan", body.GetString("username"))
	assert.Equal(t, 30, body.GetInt("age"))
}

func TestPUT(t *testing.T) {
	ts := httptest.NewServer(writeData(std.D{
		"login": "John Doe",
		"roles": []string{"moderator"},
		"id":    "b06cd03f-75d0-413a-b94b-35e155444d70",
	}))
	var dst reply
	if err := PUT(fmt.Sprintf("%s?a=1&b=2", ts.URL), std.D{"username": "chan", "gender": "female", "age": 30}, &dst); err != nil {
		t.Fatal(err)
	}
	body, _ := parser.NewParser(dst.Body)
	assert.Equal(t, http.MethodPut, dst.Method)
	assert.Equal(t, "application/json", dst.Header.Get("Content-Type"))
	assert.Equal(t, "1", dst.Query.Get("a"))
	assert.Equal(t, "2", dst.Query.Get("b"))
	assert.Equal(t, "b06cd03f-75d0-413a-b94b-35e155444d70", body.GetString("id"))
	assert.Equal(t, "John Doe", body.GetString("login"))
	assert.Equal(t, "moderator", body.GetString("roles[0]"))
	assert.Equal(t, "chan", body.GetString("username"))
	assert.Equal(t, 30, body.GetInt("age"))
}

func TestPATCH(t *testing.T) {
	ts := httptest.NewServer(writeData(std.D{
		"login": "John Doe",
		"roles": []string{"moderator"},
		"id":    "b06cd03f-75d0-413a-b94b-35e155444d70",
	}))
	var dst reply
	if err := PATCH(fmt.Sprintf("%s?a=1&b=2", ts.URL), std.D{"username": "chan", "gender": "female", "age": 30}, &dst); err != nil {
		t.Fatal(err)
	}
	body, _ := parser.NewParser(dst.Body)
	assert.Equal(t, http.MethodPatch, dst.Method)
	assert.Equal(t, "application/json", dst.Header.Get("Content-Type"))
	assert.Equal(t, "1", dst.Query.Get("a"))
	assert.Equal(t, "2", dst.Query.Get("b"))
	assert.Equal(t, "b06cd03f-75d0-413a-b94b-35e155444d70", body.GetString("id"))
	assert.Equal(t, "John Doe", body.GetString("login"))
	assert.Equal(t, "moderator", body.GetString("roles[0]"))
	assert.Equal(t, "chan", body.GetString("username"))
	assert.Equal(t, 30, body.GetInt("age"))
}

func TestDELETE(t *testing.T) {
	ts := httptest.NewServer(writeData(std.D{
		"login": "John Doe",
		"roles": []string{"moderator"},
		"id":    "b06cd03f-75d0-413a-b94b-35e155444d70",
	}))
	var dst reply
	if err := DELETE(fmt.Sprintf("%s?a=1&b=2", ts.URL), &dst); err != nil {
		t.Fatal(err)
	}
	body, _ := parser.NewParser(dst.Body)
	assert.Equal(t, http.MethodDelete, dst.Method)
	assert.Equal(t, "application/json", dst.Header.Get("Content-Type"))
	assert.Equal(t, "1", dst.Query.Get("a"))
	assert.Equal(t, "2", dst.Query.Get("b"))
	assert.Equal(t, "b06cd03f-75d0-413a-b94b-35e155444d70", body.GetString("id"))
	assert.Equal(t, "John Doe", body.GetString("login"))
	assert.Equal(t, "moderator", body.GetString("roles[0]"))
}
