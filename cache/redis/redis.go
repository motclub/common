package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/motclub/common/cache"
	"github.com/motclub/common/json"
	"strings"
	"time"
)

const redisDelKeysChannel = "__MOT_DEL_KEYS_CHANNEL__"

func NewRedisCache(opts *redis.UniversalOptions) (cache.ICache, error) {
	cmd := redis.NewUniversalClient(opts)
	if _, err := cmd.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	cache := &redisCache{rdb: cmd}
	err := cache.Subscribe([]string{redisDelKeysChannel}, func(channel string, data string) {
		if data == "" {
			return
		}

		if cache.children != nil {
			keys := strings.Split(data, ",")
			_ = cache.children.Del(keys...)
		}
	})
	return cache, err
}

type redisCache struct {
	rdb      redis.UniversalClient
	parent   cache.ICache
	children cache.ICache
}

func (r *redisCache) SetParent(parent cache.ICache) {
	if r.parent == nil {
		r.parent = parent
	}
}

func (r *redisCache) SetChildren(children cache.ICache) {
	if r.children == nil {
		r.children = children
	}
}

func (r *redisCache) Publish(channel string, message interface{}) error {
	return r.rdb.Publish(context.Background(), channel, message).Err()
}

func (r *redisCache) Subscribe(channels []string, handler func(string, string)) error {
	ps := r.rdb.Subscribe(context.Background(), channels...)
	if _, err := ps.Receive(context.Background()); err != nil {
		return err
	}
	ch := ps.Channel()
	go func(ch <-chan *redis.Message) {
		for msg := range ch {
			handler(msg.Channel, msg.Payload)
		}
	}(ch)
	return nil
}

func (r *redisCache) PSubscribe(patterns []string, handler func(string, string)) error {
	pubsub := r.rdb.PSubscribe(context.Background(), patterns...)
	if _, err := pubsub.Receive(context.Background()); err != nil {
		return err
	}
	ch := pubsub.Channel()
	go func(ch <-chan *redis.Message) {
		for msg := range ch {
			handler(msg.Channel, msg.Payload)
		}
	}(ch)
	return nil
}

func (r *redisCache) TTL(path string) (int64, bool) {
	dur, err := r.rdb.TTL(context.Background(), path).Result()
	var v int64
	if err == nil {
		if s := int64(dur.Seconds()); s > 0 {
			v = s
		}
	}
	return v, err == nil && v > 0
}

func (r *redisCache) Has(key string) bool {
	has := r.rdb.Exists(context.Background(), key).Val() > 1
	if !has && r.parent != nil {
		has = r.parent.Has(key)
	}
	return has
}

func (r *redisCache) HasGet(key string, dst interface{}) bool {
	s, err := r.rdb.Get(context.Background(), key).Result()
	has := err == nil
	if has && s != "" {
		var v struct {
			ExpiredSeconds int64       `json:"expired_seconds"`
			CreatedAt      time.Time   `json:"created_at"`
			Data           interface{} `json:"data"`
		}
		if err = json.Parse(s, &v); err == nil {
			err = json.Copy(v.Data, dst)
		}
	} else if r.parent != nil {
		has = r.parent.HasGet(key, dst)
		var ttl int64
		ttl, has = r.parent.TTL(key)
		if has {
			_ = r.Set(key, dst, ttl)
		}
	}
	return has
}

func (r *redisCache) HasGetInt(key string) (int, bool) {
	s, err := r.rdb.Get(context.Background(), key).Result()
	has := err == nil
	if !has && r.parent != nil {
		var v int
		v, has = r.parent.HasGetInt(key)
		var ttl int64
		ttl, has = r.parent.TTL(key)
		if has {
			_ = r.Set(key, v, ttl)
		}
		return v, has
	}
	var v struct {
		ExpiredSeconds int64     `json:"expired_seconds"`
		CreatedAt      time.Time `json:"created_at"`
		Data           int       `json:"data"`
	}
	err = json.Parse(s, &v)
	return v.Data, has
}

func (r *redisCache) HasGetInt8(key string) (int8, bool) {
	v, has := r.HasGetInt(key)
	return int8(v), has
}

func (r *redisCache) HasGetInt16(key string) (int16, bool) {
	v, has := r.HasGetInt(key)
	return int16(v), has
}

func (r *redisCache) HasGetInt32(key string) (int32, bool) {
	v, has := r.HasGetInt(key)
	return int32(v), has
}

func (r *redisCache) HasGetInt64(key string) (int64, bool) {
	v, has := r.HasGetInt(key)
	return int64(v), has
}

func (r *redisCache) HasGetUint(key string) (uint, bool) {
	v, has := r.HasGetUint64(key)
	return uint(v), has
}

func (r *redisCache) HasGetUint8(key string) (uint8, bool) {
	v, has := r.HasGetUint64(key)
	return uint8(v), has
}

func (r *redisCache) HasGetUint16(key string) (uint16, bool) {
	v, has := r.HasGetUint64(key)
	return uint16(v), has
}

func (r *redisCache) HasGetUint32(key string) (uint32, bool) {
	v, has := r.HasGetUint64(key)
	return uint32(v), has
}

func (r *redisCache) HasGetUint64(key string) (uint64, bool) {
	v, has := r.HasGetInt(key)
	return uint64(v), has
}

func (r *redisCache) HasGetFloat(key string) (float64, bool) {
	s, err := r.rdb.Get(context.Background(), key).Result()
	has := err == nil
	if !has && r.parent != nil {
		var v float64
		v, has = r.parent.HasGetFloat(key)
		var ttl int64
		ttl, has = r.parent.TTL(key)
		if has {
			_ = r.Set(key, v, ttl)
		}
		return v, has
	}
	var v struct {
		ExpiredSeconds int64     `json:"expired_seconds"`
		CreatedAt      time.Time `json:"created_at"`
		Data           float64   `json:"data"`
	}
	err = json.Parse(s, &v)
	return v.Data, has
}

func (r *redisCache) HasGetFloat32(key string) (float32, bool) {
	v, has := r.HasGetFloat(key)
	return float32(v), has
}

func (r *redisCache) HasGetFloat64(key string) (float64, bool) {
	return r.HasGetFloat(key)
}

func (r *redisCache) HasGetString(key string) (string, bool) {
	s, err := r.rdb.Get(context.Background(), key).Result()
	has := err == nil
	if !has && r.parent != nil {
		var v string
		v, has = r.parent.HasGetString(key)
		var ttl int64
		ttl, has = r.parent.TTL(key)
		if has {
			_ = r.Set(key, v, ttl)
		}
		return v, has
	}
	var v struct {
		ExpiredSeconds int64     `json:"expired_seconds"`
		CreatedAt      time.Time `json:"created_at"`
		Data           string    `json:"data"`
	}
	err = json.Parse(s, &v)
	return v.Data, has
}

func (r *redisCache) HasGetBool(key string) (bool, bool) {
	s, err := r.rdb.Get(context.Background(), key).Result()
	has := err == nil
	if !has && r.parent != nil {
		var v bool
		v, has = r.parent.HasGetBool(key)
		var ttl int64
		ttl, has = r.parent.TTL(key)
		if has {
			_ = r.Set(key, v, ttl)
		}
		return v, has
	}
	var v struct {
		ExpiredSeconds int64     `json:"expired_seconds"`
		CreatedAt      time.Time `json:"created_at"`
		Data           bool      `json:"data"`
	}
	err = json.Parse(s, &v)
	return v.Data, has
}

func (r *redisCache) HasGetTime(key string) (time.Time, bool) {
	var v time.Time
	has := r.HasGet(key, &v)
	return v, has
}

func (r *redisCache) Get(key string, dst interface{}) {
	_ = r.HasGet(key, dst)
}

func (r *redisCache) GetInt(key string) int {
	v, _ := r.HasGetInt(key)
	return v
}

func (r *redisCache) GetInt8(key string) int8 {
	v, _ := r.HasGetInt8(key)
	return v
}

func (r *redisCache) GetInt16(key string) int16 {
	v, _ := r.HasGetInt16(key)
	return v
}

func (r *redisCache) GetInt32(key string) int32 {
	v, _ := r.HasGetInt32(key)
	return v
}

func (r *redisCache) GetInt64(key string) int64 {
	v, _ := r.HasGetInt64(key)
	return v
}

func (r *redisCache) GetUint(key string) uint {
	v, _ := r.HasGetUint(key)
	return v
}

func (r *redisCache) GetUint8(key string) uint8 {
	v, _ := r.HasGetUint8(key)
	return v
}

func (r *redisCache) GetUint16(key string) uint16 {
	v, _ := r.HasGetUint16(key)
	return v
}

func (r *redisCache) GetUint32(key string) uint32 {
	v, _ := r.HasGetUint32(key)
	return v
}

func (r *redisCache) GetUint64(key string) uint64 {
	v, _ := r.HasGetUint64(key)
	return v
}

func (r *redisCache) GetFloat(key string) float64 {
	v, _ := r.HasGetFloat(key)
	return v
}

func (r *redisCache) GetFloat32(key string) float32 {
	v, _ := r.HasGetFloat32(key)
	return v
}

func (r *redisCache) GetFloat64(key string) float64 {
	v, _ := r.HasGetFloat64(key)
	return v
}

func (r *redisCache) GetString(key string) string {
	v, _ := r.HasGetString(key)
	return v
}

func (r *redisCache) GetBool(key string) bool {
	v, _ := r.HasGetBool(key)
	return v
}

func (r *redisCache) GetTime(key string) time.Time {
	v, _ := r.HasGetTime(key)
	return v
}

func (r *redisCache) DefaultGet(key string, dst interface{}, defaultValue interface{}) {
	if !r.HasGet(key, dst) {
		_ = json.Copy(defaultValue, dst)
	}
}

func (r *redisCache) DefaultGetInt(key string, defaultValue int) int {
	if v, has := r.HasGetInt(key); has {
		return v
	}
	return defaultValue
}

func (r *redisCache) DefaultGetInt8(key string, defaultValue int8) int8 {
	if v, has := r.HasGetInt8(key); has {
		return v
	}
	return defaultValue
}

func (r *redisCache) DefaultGetInt16(key string, defaultValue int16) int16 {
	if v, has := r.HasGetInt16(key); has {
		return v
	}
	return defaultValue
}

func (r *redisCache) DefaultGetInt32(key string, defaultValue int32) int32 {
	if v, has := r.HasGetInt32(key); has {
		return v
	}
	return defaultValue
}

func (r *redisCache) DefaultGetInt64(key string, defaultValue int64) int64 {
	if v, has := r.HasGetInt64(key); has {
		return v
	}
	return defaultValue
}

func (r *redisCache) DefaultGetUint(key string, defaultValue uint) uint {
	if v, has := r.HasGetUint(key); has {
		return v
	}
	return defaultValue
}

func (r *redisCache) DefaultGetUint8(key string, defaultValue uint8) uint8 {
	if v, has := r.HasGetUint8(key); has {
		return v
	}
	return defaultValue
}

func (r *redisCache) DefaultGetUint16(key string, defaultValue uint16) uint16 {
	if v, has := r.HasGetUint16(key); has {
		return v
	}
	return defaultValue
}

func (r *redisCache) DefaultGetUint32(key string, defaultValue uint32) uint32 {
	if v, has := r.HasGetUint32(key); has {
		return v
	}
	return defaultValue
}

func (r *redisCache) DefaultGetUint64(key string, defaultValue uint64) uint64 {
	if v, has := r.HasGetUint64(key); has {
		return v
	}
	return defaultValue
}

func (r *redisCache) DefaultGetFloat(key string, defaultValue float64) float64 {
	if v, has := r.HasGetFloat(key); has {
		return v
	}
	return defaultValue
}

func (r *redisCache) DefaultGetFloat32(key string, defaultValue float32) float32 {
	if v, has := r.HasGetFloat32(key); has {
		return v
	}
	return defaultValue
}

func (r *redisCache) DefaultGetFloat64(key string, defaultValue float64) float64 {
	if v, has := r.HasGetFloat64(key); has {
		return v
	}
	return defaultValue
}

func (r *redisCache) DefaultGetString(key string, defaultValue string) string {
	if v, has := r.HasGetString(key); has {
		return v
	}
	return defaultValue
}

func (r *redisCache) DefaultGetBool(key string, defaultValue bool) bool {
	if v, has := r.HasGetBool(key); has {
		return v
	}
	return defaultValue
}

func (r *redisCache) DefaultGetTime(key string, defaultValue time.Time) time.Time {
	if v, has := r.HasGetTime(key); has {
		return v
	}
	return defaultValue
}

func (r *redisCache) Parent() cache.ICache {
	return r.parent
}

func (r *redisCache) Children() cache.ICache {
	return r.children
}

func (r *redisCache) HasPrefix(s string, limit ...int) (map[string]string, error) {
	v, err := r.contains(fmt.Sprintf("%s*", s), limit...)

	if err == nil && len(v) == 0 && r.parent != nil {
		return r.parent.HasPrefix(s, limit...)
	}

	return v, err
}

func (r *redisCache) HasSuffix(s string, limit ...int) (map[string]string, error) {
	v, err := r.contains(fmt.Sprintf("*%s", s), limit...)

	if err == nil && len(v) == 0 && r.parent != nil {
		return r.parent.HasSuffix(s, limit...)
	}

	return v, err
}

func (r *redisCache) Contains(s string, limit ...int) (map[string]string, error) {
	v, err := r.contains(fmt.Sprintf("*%s*", s), limit...)

	if err == nil && len(v) == 0 && r.parent != nil {
		return r.parent.Contains(s, limit...)
	}

	return v, err
}

func (r *redisCache) Set(key string, value interface{}, expiresIn ...int64) error {
	var exp int64
	if len(expiresIn) > 0 {
		exp = expiresIn[0]
	}
	dur := time.Duration(exp) * time.Second
	cv := struct {
		ExpiredSeconds int64       `json:"expired_seconds"`
		CreatedAt      time.Time   `json:"created_at"`
		Data           interface{} `json:"data"`
	}{
		ExpiredSeconds: exp,
		CreatedAt:      time.Now(),
		Data:           value,
	}
	v := json.Stringify(&cv, false)
	err := r.rdb.Set(context.Background(), key, v, dur).Err()
	if err == nil && r.parent != nil {
		err = r.parent.Set(key, value, expiresIn...)
	}
	return err
}

func (r *redisCache) Incr(key string) (int, error) {
	if r.parent != nil {
		return r.parent.Incr(key)
	}

	v, err := r.rdb.Incr(context.Background(), key).Result()
	return int(v), err
}

func (r *redisCache) IncrBy(key string, step int) (int, error) {
	if r.parent != nil {
		return r.parent.IncrBy(key, step)
	}

	v, err := r.rdb.IncrBy(context.Background(), key, int64(step)).Result()
	return int(v), err
}

func (r *redisCache) IncrByFloat(key string, step float64) (float64, error) {
	if r.parent != nil {
		return r.parent.IncrByFloat(key, step)
	}

	return r.rdb.IncrByFloat(context.Background(), key, step).Result()
}

func (r *redisCache) Del(keys ...string) error {
	err := r.rdb.Del(context.Background(), keys...).Err()

	if r.parent != nil {
		err = r.parent.Del(keys...)
	} else {
		err = r.Publish(redisDelKeysChannel, strings.Join(keys, ","))
	}

	return err
}

func (r *redisCache) Close() error {
	if r.rdb != nil {
		return r.rdb.Close()
	}
	return nil
}

func (r *redisCache) contains(pattern string, limit ...int) (map[string]string, error) {

	var count int64
	if len(limit) > 0 {
		count = int64(limit[0])
	}

	var (
		cursor uint64
		values = make(map[string]string)
	)

	for {
		var keys []string
		var err error
		keys, cursor, err = r.rdb.Scan(context.Background(), cursor, pattern, count).Result()
		if err != nil {
			return nil, err
		}
		for _, key := range keys {
			s, err := r.rdb.Get(context.Background(), key).Result()
			if err != nil {
				return nil, err
			}
			var v struct {
				ExpiredSeconds int64     `json:"expired_seconds"`
				CreatedAt      time.Time `json:"created_at"`
				Data           string    `json:"data"`
			}
			if err = json.Parse(s, &v); err == nil {
				values[key] = v.Data
			}
		}
		if cursor == 0 {
			break
		}
	}

	return values, nil
}
