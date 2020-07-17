package cache

import (
	"github.com/motclub/common/getter"
	"github.com/pkg/errors"
	"time"
)

func NewCache(caches []ICache) (ICache, error) {
	if len(caches) < 2 {
		return nil, errors.New(`mot: minimum cache level is 2`)
	}
	for index, cache := range caches {
		if index == 0 {
			cache.SetParent(caches[index+1])
		} else if index > 0 && index < len(caches)-1 {
			cache.SetParent(caches[index+1])
			cache.SetChildren(caches[index-1])
		} else if index == len(caches)-1 {
			cache.SetChildren(caches[index-1])
		}
	}
	return caches[0], nil
}

// ICache
type ICache interface {
	getter.IGetter

	TTL(key string) (time.Duration, bool)
	Set(key string, value interface{}, expiration ...time.Duration) error
	HasPrefix(s string, limit ...int) (map[string]string, error)
	HasSuffix(s string, limit ...int) (map[string]string, error)
	Contains(s string, limit ...int) (map[string]string, error)
	Incr(key string) (int, error)
	IncrBy(key string, step int) (int, error)
	IncrByFloat(key string, step float64) (float64, error)
	Del(keys ...string) error
	Parent() ICache
	Children() ICache
	SetParent(parent ICache)
	SetChildren(children ICache)
	Close() error

	Publish(channel string, message interface{}) error
	Subscribe(channels []string, handler func(string, string)) error
	PSubscribe(patterns []string, handler func(string, string)) error
}
