package http

import (
	"fmt"
	"github.com/motclub/common/caller"
	"github.com/motclub/common/httpx"
	"github.com/motclub/common/json"
	"github.com/motclub/common/std"
	"net/http"
	"time"
)

func init() {
	caller.RegisterCaller(&httpCaller{})
}

type httpCaller struct {
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Query       map[string]string `json:"query"`
	Headers     map[string]string `json:"headers"`
	IdleTimeout time.Duration     `json:"idle_timeout"`
	MaxRetries  uint              `json:"max_retries"`
}

func (c *httpCaller) New(settings interface{}) (caller.ICaller, error) {
	var cc httpCaller
	err := json.Copy(settings, &cc)
	return &cc, err
}

func (c *httpCaller) Kind() string {
	return "HTTP"
}

func (c *httpCaller) Call(entry *caller.ServiceEntry, args *std.Args) *std.Reply {
	if c.Method == "" {
		c.Method = http.MethodGet
	}
	var (
		err   error
		reply std.Reply
		opts  = httpx.RequestOptions{
			Headers: c.Headers,
			Query:   c.Query,
		}
	)
	if c.Method == http.MethodGet {
		c.Query["args"] = json.Stringify(args, false)
		err = httpx.GET(c.URL, &reply, &opts)
	} else {
		err = httpx.Request(c.Method, c.URL, args, &reply, &opts)
	}
	if err != nil {
		return &std.Reply{
			Code:              -1,
			Data:              fmt.Errorf("failed to request: %s", err.Error()),
			Message:           "Service call failed.",
			LocaleMessageName: "mot_service_call_failed",
		}
	}
	return &reply
}

func (c *httpCaller) Close() error {
	return nil
}
