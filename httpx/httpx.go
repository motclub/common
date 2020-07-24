package httpx

import (
	"errors"
	"fmt"
	"github.com/motclub/common/json"
	"github.com/motclub/common/reflectx"
	"github.com/sethgrid/pester"
	"gopkg.in/guregu/null.v4"
	"io"
	"io/ioutil"
	"net/http"
	netUrl "net/url"
	"reflect"
	"strings"
	"time"
)

var backoffStrategyMap = map[string]pester.BackoffStrategy{
	"default":            pester.DefaultBackoff,
	"linear":             pester.LinearBackoff,
	"linear_jitter":      pester.LinearJitterBackoff,
	"exponential":        pester.ExponentialBackoff,
	"exponential_jitter": pester.ExponentialJitterBackoff,
}

type Options struct {
	Url         string                 `json:"schema"`
	Method      string                 `json:"method"`
	Headers     map[string]string      `json:"headers,omitempty"`
	Query       map[string]interface{} `json:"query,omitempty"`
	ContentType string                 `json:"contentType"`
	Body        interface{}

	Timeout     string   `json:"timeout"`
	Concurrency null.Int `json:"concurrency"`
	MaxRetries  null.Int `json:"maxRetries"`
	Backoff     string   `json:"backoff"`

	client  *pester.Client
	header  http.Header
	body    io.Reader
	timeout time.Duration
	backoff pester.BackoffStrategy
}

func GET(url string, dst interface{}, options ...*Options) error {
	opts := resolveOptions(options)
	opts.Url = url
	opts.Method = http.MethodGet
	return Do(opts, dst)
}

func POST(url string, body interface{}, dst interface{}, options ...*Options) error {
	opts := resolveOptions(options)
	opts.Url = url
	opts.Method = http.MethodPost
	opts.Body = body
	return Do(opts, dst)
}

func POSTForm(url string, body interface{}, dst interface{}, options ...*Options) error {
	opts := resolveOptions(options)
	opts.Url = url
	opts.Method = http.MethodPost
	opts.Body = body
	opts.ContentType = "application/x-www-form-urlencoded"
	return Do(opts, dst)
}

func PUT(url string, body interface{}, dst interface{}, options ...*Options) error {
	opts := resolveOptions(options)
	opts.Url = url
	opts.Method = http.MethodPut
	opts.Body = body
	return Do(opts, dst)
}

func PATCH(url string, body interface{}, dst interface{}, options ...*Options) error {
	opts := resolveOptions(options)
	opts.Url = url
	opts.Method = http.MethodPatch
	opts.Body = body
	return Do(opts, dst)
}

func DELETE(url string, dst interface{}, options ...*Options) error {
	opts := resolveOptions(options)
	opts.Url = url
	opts.Method = http.MethodDelete
	return Do(opts, dst)
}

func Do(opts *Options, dst interface{}) error {
	if opts == nil {
		return errors.New("request options is nil")
	}
	if err := parseOptions(opts); err != nil {
		return err
	}
	req, err := http.NewRequest(opts.Method, opts.Url, opts.body)
	if err != nil {
		return err
	}
	if len(opts.header) > 0 {
		req.Header = opts.header
	}
	resp, err := opts.client.Do(req)
	// 无论如何记得消费响应体
	if resp != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}
	if err != nil {
		return err
	}
	if dst != nil {
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if err := json.STD().Unmarshal(buf, dst); err != nil {
			return err
		}
	} else {
		_, err = io.Copy(ioutil.Discard, resp.Body)
	}

	return nil
}

func resolveOptions(options []*Options) *Options {
	var opts *Options
	if len(options) > 0 && options[0] != nil {
		opts = options[0]
	} else {
		opts = &Options{}
	}
	return opts
}

func parseOptions(opts *Options) error {
	opts.client = pester.New()
	if opts.Timeout != "" {
		v, err := time.ParseDuration(opts.Timeout)
		if err != nil {
			return err
		}
		opts.client.Timeout = v
	}
	if opts.Concurrency.Valid {
		opts.client.Concurrency = int(opts.Concurrency.Int64)
	}
	if opts.MaxRetries.Valid {
		opts.client.MaxRetries = int(opts.MaxRetries.Int64)
	}
	if v := backoffStrategyMap[opts.Backoff]; v != nil {
		opts.client.Backoff = v
	}
	if opts.Method == "" {
		opts.Method = http.MethodGet
	}
	if !reflectx.IsNil(opts.Body) {
		data, err := json.STD().Marshal(opts.Body)
		if err != nil {
			return err
		}
		opts.body = strings.NewReader(string(data))
	}
	var query netUrl.Values
	if opts.Query != nil {
		query = make(netUrl.Values)
		for k, v := range opts.Query {
			switch reflect.Indirect(reflect.ValueOf(v)).Kind() {
			case reflect.Struct:
				query.Set(k, json.Stringify(v, false))
			case reflect.Map:
				query.Set(k, json.Stringify(v, false))
			default:
				query.Set(k, fmt.Sprintf("%v", v))
			}
		}
	}
	if len(query) > 0 {
		nu, err := netUrl.Parse(opts.Url)
		if err != nil {
			return err
		}
		a := nu.Query()
		for k := range query {
			a.Set(k, query.Get(k))
		}
		nu.RawQuery = a.Encode()
		opts.Url = nu.String()
	}
	if opts.ContentType == "" {
		opts.ContentType = "application/json"
	}
	opts.header = make(http.Header)
	opts.header.Set("Content-Type", opts.ContentType)
	if len(opts.Headers) > 0 {
		for k, v := range opts.Headers {
			values := strings.Split(v, ",")
			for _, value := range values {
				value = strings.TrimSpace(value)
				if value != "" {
					if opts.header.Get(k) != "" {
						opts.header.Add(k, value)
					} else {
						opts.header.Set(k, value)
					}
				}
			}
		}
	}
	return nil
}
