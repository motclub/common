package httpx

import (
	"github.com/motclub/common/helpers"
	"github.com/motclub/common/json"
	"github.com/motclub/common/reflectx"
	"github.com/sethgrid/pester"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var httpClient = pester.New()

func init() {
	httpClient.Concurrency = 3
	httpClient.MaxRetries = 3
	httpClient.Backoff = pester.ExponentialBackoff
	httpClient.KeepLog = true
}

type RequestOptions struct {
	Headers map[string]string
	Query   map[string]string
}

func GET(url string, dst interface{}, options ...*RequestOptions) error {
	return Request(http.MethodGet, url, nil, dst, options...)
}

func POST(url string, body interface{}, dst interface{}, options ...*RequestOptions) error {
	return Request(http.MethodPost, url, body, dst, options...)
}

func PUT(url string, body interface{}, dst interface{}, options ...*RequestOptions) error {
	return Request(http.MethodPut, url, body, dst, options...)
}

func PATCH(url string, body interface{}, dst interface{}, options ...*RequestOptions) error {
	return Request(http.MethodPatch, url, body, dst, options...)
}

func DELETE(url string, dst interface{}, options ...*RequestOptions) error {
	return Request(http.MethodDelete, url, nil, dst, options...)
}

func Request(method string, url string, args interface{}, reply interface{}, options ...*RequestOptions) error {
	var opts = &RequestOptions{}
	if len(options) > 0 && options[0] != nil {
		opts = options[0]
	}
	if opts.Headers == nil {
		opts.Headers = make(map[string]string)
	}
	if opts.Headers["Content-Type"] == "" {
		opts.Headers["Content-Type"] = "application/json"
	}
	// 构建请求
	var body io.Reader
	if !reflectx.IsNil(args) {
		data, err := json.STD().Marshal(args)
		if err != nil {
			return err
		}
		body = strings.NewReader(string(data))
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	for k, v := range opts.Headers {
		req.Header.Add(k, v)
	}
	if len(opts.Query) > 0 {
		query := make(map[string]interface{})
		for k, v := range opts.Query {
			query[k] = v
		}
		url = helpers.SetUrlQuery(url, query)
	}
	// 执行请求
	resp, err := httpClient.Do(req)
	// 处理返回
	if resp != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}
	if err != nil {
		return err
	}
	if reply != nil {
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if err := json.STD().Unmarshal(buf, reply); err != nil {
			return err
		}
	} else {
		_, err = io.Copy(ioutil.Discard, resp.Body)
	}
	return err
}
