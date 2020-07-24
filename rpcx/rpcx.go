package rpcx

import (
	"context"
	"github.com/smallnest/rpcx/client"
	"gopkg.in/guregu/null.v4"
	"sync"
	"time"
)

type Options struct {
	Server         string   `json:"server"`
	Path           string   `json:"path"`
	Method         string   `json:"method"`
	ConnectTimeout string   `json:"connectTimeout"`
	Retries        null.Int `json:"retries"`
}

var (
	clients   = make(map[string]client.XClient)
	clientsMu sync.RWMutex
)

func getClient(opts *Options) client.XClient {
	clientsMu.RLock()
	xc := clients[opts.Server]
	clientsMu.RUnlock()

	if xc == nil {
		d := client.NewPeer2PeerDiscovery(opts.Server, "")
		xc = client.NewXClient(opts.Path, g.FailMode, g.SelectMode, d.Clone(opts.Path), g.Option)
		clientOpts := client.DefaultOption
		if opts.ConnectTimeout != "" {
			dur, err := time.ParseDuration(opts.ConnectTimeout)
			if err == nil {
				clientOpts.ConnectTimeout = dur
			}
		}
		if opts.Retries.Valid {
			clientOpts.Retries = int(opts.Retries.Int64)
		}
		xc = client.NewXClient(opts.Path, client.Failtry, client.RandomSelect, d.Clone(opts.Path), clientOpts)

		clientsMu.Lock()
		clients[opts.Server] = xc
		clientsMu.Unlock()
	}

	return xc
}

func Call(opts *Options, args interface{}, dst interface{}) error {
	xc := getClient(opts)
	err := xc.Call(context.Background(), opts.Method, args, dst)
	if err != nil {
		return err
	}
	return nil
}
