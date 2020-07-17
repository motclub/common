package rpcx

import (
	"context"
	"fmt"
	"github.com/motclub/common/caller"
	"github.com/motclub/common/json"
	"github.com/motclub/common/std"
	"github.com/smallnest/rpcx/client"
	"sync"
)

func init() {
	caller.RegisterCaller(&xCaller{})
}

var (
	clients   = make(map[string]client.XClient)
	clientsMu sync.RWMutex
)

type xCaller struct {
	Server        string `json:"server"`
	ServicePath   string `json:"service_path"`
	ServiceMethod string `json:"service_method"`
	MaxRetries    int
}

func (c *xCaller) client() client.XClient {
	server := c.Server

	clientsMu.RLock()
	xc := clients[server]
	clientsMu.RUnlock()

	if xc == nil {
		d := client.NewPeer2PeerDiscovery(server, "")
		xc = client.NewXClient(c.ServicePath, client.Failtry, client.RandomSelect, d, client.DefaultOption)

		clientsMu.Lock()
		clients[server] = xc
		clientsMu.Unlock()
	}

	return xc
}

func (c *xCaller) New(settings std.D) (caller.ICaller, error) {
	var cc xCaller
	err := json.Copy(settings, &cc)
	return &cc, err
}

func (c *xCaller) Kind() string {
	return "RPCX"
}

func (c *xCaller) Call(entry *caller.ServiceEntry, args *std.Args) *std.Reply {
	xc := c.client()

	var reply std.Reply
	err := xc.Call(context.Background(), c.ServiceMethod, args, &reply)
	if err != nil {
		return &std.Reply{
			Code:              -1,
			Data:              fmt.Errorf("failed to call service: %s", err.Error()),
			Message:           "Service call failed.",
			LocaleMessageName: "mot_service_call_failed",
		}
	}
	return &reply
}

func (c *xCaller) Close() error {
	clientsMu.RLock()
	defer clientsMu.RUnlock()

	var err error
	for _, xc := range clients {
		if e := xc.Close(); e != nil {
			err = e
		}
	}
	return err
}
