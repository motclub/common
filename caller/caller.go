package caller

import (
	"fmt"
	"github.com/motclub/common/cache"
	"github.com/motclub/common/std"
)

const key = "mot_services"

func New(cache cache.ICache) *Caller {
	return &Caller{cache: cache}
}

type Service struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Spec std.D  `json:"spec"`
}

type Caller struct {
	cache cache.ICache
}

func (c *Caller) CallService(name string, args *std.Args) *std.Reply {
	service := c.GetService(name)
	if service == nil {
		return &std.Reply{
			Code:              -1,
			Data:              fmt.Errorf("failed to create caller: %s", err.Error()),
			Message:           "Service call failed.",
			LocaleMessageName: "mot.service.call.failed",
		}
	}
	return nil
}

func (c *Caller) GetService(name string) *Service {
	services := c.GetServices()
	if v, has := services[name]; has {
		return &v
	}
	return nil
}

func (c *Caller) GetServices() map[string]Service {
	services := make(map[string]Service)
	c.cache.Get(key, &services)
	return services
}
