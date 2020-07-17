package caller

import (
	"fmt"
	"github.com/motclub/common/std"
	"strings"
	"sync"
)

var (
	callers   = make(map[string]ICaller)
	callersMu sync.RWMutex
)

func RegisterCaller(c ICaller) {
	if c == nil || c.Kind() == "" {
		return
	}

	callersMu.Lock()
	defer callersMu.Unlock()

	callers[c.Kind()] = c
}

type ICaller interface {
	New(settings interface{}) (ICaller, error)
	Kind() string
	Call(entry *ServiceEntry, args *std.Args) *std.Reply
	Close() error
}

type ServiceEntry struct {
	Name     string
	Kind     string
	Retry    bool
	Settings std.D
}

func Call(entry *ServiceEntry, args *std.Args) *std.Reply {
	if entry == nil || args == nil {
		return nil
	}
	kind := strings.ToUpper(entry.Kind)

	callersMu.RLock()
	caller := callers[kind]
	callersMu.RUnlock()
	if caller == nil {
		return &std.Reply{
			Code:              -1,
			Data:              fmt.Errorf("unsupported caller type: %s", kind),
			Message:           "Service call failed.",
			LocaleMessageName: "mot_service_call_failed",
		}
	}

	c, err := caller.New(entry.Settings)
	if err != nil {
		return &std.Reply{
			Code:              -1,
			Data:              fmt.Errorf("failed to create caller: %s", err.Error()),
			Message:           "Service call failed.",
			LocaleMessageName: "mot_service_call_failed",
		}
	}

	return c.Call(entry, args)
}
