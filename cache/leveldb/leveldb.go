package leveldb

import (
	"bytes"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/motclub/common/cache"
	"github.com/motclub/common/json"
	"github.com/motclub/common/std"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"os"
	"time"
)

var ErrCacheUnsupportedPubSub = errors.New(`mot: unsupported Publish/Subscribe messages`)

func NewLevelDBCache(path string, o *opt.Options) (cache.ICache, error) {
	if fi, err := os.Stat(path); err == nil {
		if !fi.IsDir() {
			return nil, fmt.Errorf("leveldb/storage: open %s: not a directory", path)
		}
		if err := os.RemoveAll(path); err != nil {
			return nil, err
		}
	}

	db, err := leveldb.OpenFile(path, o)
	if err != nil {
		return nil, err
	}
	return &levelDBCache{db: db}, nil
}

type levelDBCacheValue struct {
	ExpiredSeconds int64     `json:"expired_seconds"`
	CreatedAt      time.Time `json:"created_at"`
	Data           []byte    `json:"data"`
}

type levelDBCache struct {
	db       *leveldb.DB
	parent   cache.ICache
	children cache.ICache
}

func (l *levelDBCache) Publish(channel string, message interface{}) error {
	if l.parent == nil {
		return ErrCacheUnsupportedPubSub
	}
	return l.parent.Publish(channel, message)
}

func (l *levelDBCache) Subscribe(channels []string, handler func(string, string)) error {
	if l.parent == nil {
		return ErrCacheUnsupportedPubSub
	}
	return l.parent.Subscribe(channels, handler)
}

func (l *levelDBCache) PSubscribe(patterns []string, handler func(string, string)) error {
	if l.parent == nil {
		return ErrCacheUnsupportedPubSub
	}
	return l.parent.PSubscribe(patterns, handler)
}

func (l *levelDBCache) hasGet(path string) (*levelDBCacheValue, bool) {
	var v levelDBCacheValue
	data, err := l.db.Get([]byte(path), nil)
	if err == nil {
		err = json.STD().Unmarshal(data, &v)
		if v.ExpiredSeconds > 0 {
			expiredAt := v.CreatedAt.Add(time.Duration(v.ExpiredSeconds) * time.Second)
			if time.Now().After(expiredAt) {
				err = l.db.Delete([]byte(path), nil)
				return nil, false
			}
		} else {
			err = l.db.Delete([]byte(path), nil)
			return nil, false
		}
	}
	return &v, err != leveldb.ErrNotFound
}

func (l *levelDBCache) TTL(path string) (int64, bool) {
	v, has := l.hasGet(path)
	if !has {
		return 0, has
	}
	expiredAt := v.CreatedAt.Add(time.Duration(v.ExpiredSeconds) * time.Second)
	secs := time.Now().Sub(expiredAt).Seconds()
	return int64(secs), has
}

func (l *levelDBCache) Has(path string) bool {
	_, has := l.hasGet(path)
	if !has && l.parent != nil {
		has = l.parent.Has(path)
	}
	return has
}

func (l *levelDBCache) HasGet(path string, dst interface{}) bool {
	cv, has := l.hasGet(path)
	if has {
		_ = json.STD().Unmarshal(cv.Data, dst)
	} else if l.parent != nil {
		has = l.parent.HasGet(path, dst)
		var ttl int64
		ttl, has = l.parent.TTL(path)
		if has {
			_ = l.Set(path, dst, ttl)
		}
	}
	return has
}

func (l *levelDBCache) HasGetInt(path string) (int, bool) {
	cv, has := l.hasGet(path)
	if has {
		v, _ := jsonparser.ParseInt(cv.Data)
		return int(v), has
	} else if l.parent != nil {
		var v int
		v, has = l.parent.HasGetInt(path)
		var ttl int64
		ttl, has = l.parent.TTL(path)
		if has {
			_ = l.Set(path, v, ttl)
		}
		return v, has
	}
	return 0, has
}

func (l *levelDBCache) HasGetInt8(path string) (int8, bool) {
	v, has := l.HasGetInt(path)
	return int8(v), has
}

func (l *levelDBCache) HasGetInt16(path string) (int16, bool) {
	v, has := l.HasGetInt(path)
	return int16(v), has
}

func (l *levelDBCache) HasGetInt32(path string) (int32, bool) {
	v, has := l.HasGetInt(path)
	return int32(v), has
}

func (l *levelDBCache) HasGetInt64(path string) (int64, bool) {
	v, has := l.HasGetInt(path)
	return int64(v), has
}

func (l *levelDBCache) HasGetUint(path string) (uint, bool) {
	v, has := l.HasGetInt(path)
	return uint(v), has
}

func (l *levelDBCache) HasGetUint8(path string) (uint8, bool) {
	v, has := l.HasGetInt(path)
	return uint8(v), has
}

func (l *levelDBCache) HasGetUint16(path string) (uint16, bool) {
	v, has := l.HasGetInt(path)
	return uint16(v), has
}

func (l *levelDBCache) HasGetUint32(path string) (uint32, bool) {
	v, has := l.HasGetInt(path)
	return uint32(v), has
}

func (l *levelDBCache) HasGetUint64(path string) (uint64, bool) {
	v, has := l.HasGetInt(path)
	return uint64(v), has
}

func (l *levelDBCache) HasGetFloat(path string) (float64, bool) {
	cv, has := l.hasGet(path)
	if has {
		v, _ := jsonparser.ParseFloat(cv.Data)
		return v, has
	} else if l.parent != nil {
		var v float64
		v, has = l.parent.HasGetFloat(path)
		var ttl int64
		ttl, has = l.parent.TTL(path)
		if has {
			_ = l.Set(path, v, ttl)
		}
		return v, has
	}
	return 0, has
}

func (l *levelDBCache) HasGetFloat32(path string) (float32, bool) {
	v, has := l.HasGetFloat(path)
	return float32(v), has
}

func (l *levelDBCache) HasGetFloat64(path string) (float64, bool) {
	return l.HasGetFloat(path)
}

func (l *levelDBCache) HasGetString(path string) (string, bool) {
	cv, has := l.hasGet(path)
	if has {
		v, _ := jsonparser.ParseString(cv.Data)
		return v, has
	} else if l.parent != nil {
		var v string
		v, has = l.parent.HasGetString(path)
		var ttl int64
		ttl, has = l.parent.TTL(path)
		if has {
			_ = l.Set(path, v, ttl)
		}
		return v, has
	}
	return "", has
}

func (l *levelDBCache) HasGetBool(path string) (bool, bool) {
	cv, has := l.hasGet(path)
	if has {
		v, _ := jsonparser.ParseBoolean(cv.Data)
		return v, has
	} else if l.parent != nil {
		var v bool
		v, has = l.parent.HasGetBool(path)
		var ttl int64
		ttl, has = l.parent.TTL(path)
		if has {
			_ = l.Set(path, v, ttl)
		}
		return v, has
	}
	return false, has
}

func (l *levelDBCache) HasGetTime(path string) (time.Time, bool) {
	cv, has := l.hasGet(path)
	if has {
		var v struct{ Data time.Time }
		_ = json.Parse(fmt.Sprintf(`{"data":"%s"}`, string(cv.Data)), &v)
		return v.Data, has
	} else if l.parent != nil {
		var v time.Time
		v, has = l.parent.HasGetTime(path)
		var ttl int64
		ttl, has = l.parent.TTL(path)
		if has {
			_ = l.Set(path, v, ttl)
		}
		return v, has
	}
	return time.Time{}, has
}

func (l *levelDBCache) Get(path string, dst interface{}) {
	_ = l.HasGet(path, dst)
}

func (l *levelDBCache) GetInt(path string) int {
	v, _ := l.HasGetInt(path)
	return v
}

func (l *levelDBCache) GetInt8(path string) int8 {
	v, _ := l.HasGetInt8(path)
	return v
}

func (l *levelDBCache) GetInt16(path string) int16 {
	v, _ := l.HasGetInt16(path)
	return v
}

func (l *levelDBCache) GetInt32(path string) int32 {
	v, _ := l.HasGetInt32(path)
	return v
}

func (l *levelDBCache) GetInt64(path string) int64 {
	v, _ := l.HasGetInt64(path)
	return v
}

func (l *levelDBCache) GetUint(path string) uint {
	v, _ := l.HasGetUint(path)
	return v
}

func (l *levelDBCache) GetUint8(path string) uint8 {
	v, _ := l.HasGetUint8(path)
	return v
}

func (l *levelDBCache) GetUint16(path string) uint16 {
	v, _ := l.HasGetUint16(path)
	return v
}

func (l *levelDBCache) GetUint32(path string) uint32 {
	v, _ := l.HasGetUint32(path)
	return v
}

func (l *levelDBCache) GetUint64(path string) uint64 {
	v, _ := l.HasGetUint64(path)
	return v
}

func (l *levelDBCache) GetFloat(path string) float64 {
	v, _ := l.HasGetFloat(path)
	return v
}

func (l *levelDBCache) GetFloat32(path string) float32 {
	v, _ := l.HasGetFloat32(path)
	return v
}

func (l *levelDBCache) GetFloat64(path string) float64 {
	v, _ := l.HasGetFloat64(path)
	return v
}

func (l *levelDBCache) GetString(path string) string {
	v, _ := l.HasGetString(path)
	return v
}

func (l *levelDBCache) GetBool(path string) bool {
	v, _ := l.HasGetBool(path)
	return v
}

func (l *levelDBCache) GetTime(path string) time.Time {
	v, _ := l.HasGetTime(path)
	return v
}

func (l *levelDBCache) DefaultGet(path string, dst interface{}, defaultValue interface{}) {
	if !l.HasGet(path, dst) {
		_ = json.Copy(defaultValue, dst)
	}
}

func (l *levelDBCache) DefaultGetInt(path string, defaultValue int) int {
	if v, has := l.HasGetInt(path); has {
		return v
	}
	return defaultValue
}

func (l *levelDBCache) DefaultGetInt8(path string, defaultValue int8) int8 {
	if v, has := l.HasGetInt8(path); has {
		return v
	}
	return defaultValue
}

func (l *levelDBCache) DefaultGetInt16(path string, defaultValue int16) int16 {
	if v, has := l.HasGetInt16(path); has {
		return v
	}
	return defaultValue
}

func (l *levelDBCache) DefaultGetInt32(path string, defaultValue int32) int32 {
	if v, has := l.HasGetInt32(path); has {
		return v
	}
	return defaultValue
}

func (l *levelDBCache) DefaultGetInt64(path string, defaultValue int64) int64 {
	if v, has := l.HasGetInt64(path); has {
		return v
	}
	return defaultValue
}

func (l *levelDBCache) DefaultGetUint(path string, defaultValue uint) uint {
	if v, has := l.HasGetUint(path); has {
		return v
	}
	return defaultValue
}

func (l *levelDBCache) DefaultGetUint8(path string, defaultValue uint8) uint8 {
	if v, has := l.HasGetUint8(path); has {
		return v
	}
	return defaultValue
}

func (l *levelDBCache) DefaultGetUint16(path string, defaultValue uint16) uint16 {
	if v, has := l.HasGetUint16(path); has {
		return v
	}
	return defaultValue
}

func (l *levelDBCache) DefaultGetUint32(path string, defaultValue uint32) uint32 {
	if v, has := l.HasGetUint32(path); has {
		return v
	}
	return defaultValue
}

func (l *levelDBCache) DefaultGetUint64(path string, defaultValue uint64) uint64 {
	if v, has := l.HasGetUint64(path); has {
		return v
	}
	return defaultValue
}

func (l *levelDBCache) DefaultGetFloat(path string, defaultValue float64) float64 {
	if v, has := l.HasGetFloat(path); has {
		return v
	}
	return defaultValue
}

func (l *levelDBCache) DefaultGetFloat32(path string, defaultValue float32) float32 {
	if v, has := l.HasGetFloat32(path); has {
		return v
	}
	return defaultValue
}

func (l *levelDBCache) DefaultGetFloat64(path string, defaultValue float64) float64 {
	if v, has := l.HasGetFloat64(path); has {
		return v
	}
	return defaultValue
}

func (l *levelDBCache) DefaultGetString(path string, defaultValue string) string {
	if v, has := l.HasGetString(path); has {
		return v
	}
	return defaultValue
}

func (l *levelDBCache) DefaultGetBool(path string, defaultValue bool) bool {
	if v, has := l.HasGetBool(path); has {
		return v
	}
	return defaultValue
}

func (l *levelDBCache) DefaultGetTime(path string, defaultValue time.Time) time.Time {
	if v, has := l.HasGetTime(path); has {
		return v
	}
	return defaultValue
}

func (l *levelDBCache) Set(key string, value interface{}, expiresIn ...int64) error {
	s := json.Stringify(std.D{"data": value}, false)
	raw, _, _, err := jsonparser.Get([]byte(s), "data")
	if err != nil {
		return err
	}

	var exp int64
	if len(expiresIn) > 0 {
		exp = expiresIn[0]
	}
	cv := &levelDBCacheValue{
		ExpiredSeconds: exp,
		CreatedAt:      time.Now(),
		Data:           raw,
	}
	data, err := json.STD().Marshal(cv)
	if err != nil {
		return err
	}
	err = l.db.Put([]byte(key), data, nil)
	if err == nil && l.parent != nil {
		err = l.parent.Set(key, value, expiresIn...)
	}
	return err
}

func (l *levelDBCache) HasPrefix(s string, limit ...int) (map[string]string, error) {
	keyword := []byte(s)
	v, err := l.filter(func(key, value []byte) bool {
		return bytes.HasPrefix(key, keyword)
	}, limit...)

	if err == nil && len(v) == 0 && l.parent != nil {
		return l.parent.HasPrefix(s, limit...)
	}

	return v, err
}

func (l *levelDBCache) HasSuffix(s string, limit ...int) (map[string]string, error) {
	keyword := []byte(s)
	v, err := l.filter(func(key, value []byte) bool {
		return bytes.HasSuffix(key, keyword)
	}, limit...)

	if err == nil && len(v) == 0 && l.parent != nil {
		return l.parent.HasSuffix(s, limit...)
	}

	return v, err
}

func (l *levelDBCache) Contains(s string, limit ...int) (map[string]string, error) {
	keyword := []byte(s)
	v, err := l.filter(func(key, value []byte) bool {
		return bytes.Contains(key, keyword)
	}, limit...)

	if err == nil && len(v) == 0 && l.parent != nil {
		return l.parent.Contains(s, limit...)
	}

	return v, err
}

func (l *levelDBCache) filter(filter func(key, value []byte) bool, limit ...int) (map[string]string, error) {
	var max int
	if len(limit) > 0 {
		max = limit[0]
	}

	var v = make(map[string]string)
	iter := l.db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		if filter(key, value) {
			v[string(key)] = string(value)
		}
		if max > 0 && len(v) >= max {
			break
		}
	}
	iter.Release()
	return v, iter.Error()
}

func (l *levelDBCache) incr(key string, step int) (int, error) {
	v := l.GetInt(key)
	v = v + step
	err := l.Set(key, v)
	return v, err
}

func (l *levelDBCache) Incr(key string) (int, error) {
	if l.parent != nil {
		return l.parent.Incr(key)
	}

	return l.incr(key, 1)
}

func (l *levelDBCache) IncrBy(key string, step int) (int, error) {
	if l.parent != nil {
		return l.parent.IncrBy(key, step)
	}

	return l.incr(key, step)
}

func (l *levelDBCache) IncrByFloat(key string, step float64) (float64, error) {
	if l.parent != nil {
		return l.parent.IncrByFloat(key, step)
	}

	v := l.GetFloat(key)
	v = v + step
	err := l.Set(key, v)
	return v, err
}

func (l *levelDBCache) Del(keys ...string) error {
	var err error
	for _, key := range keys {
		e := l.db.Delete([]byte(key), nil)
		if err == nil && e != nil {
			err = e
		}
	}
	if l.parent != nil {
		err = l.parent.Del(keys...)
	}
	return err
}

func (l *levelDBCache) Parent() cache.ICache {
	return l.parent
}

func (l *levelDBCache) SetParent(parent cache.ICache) {
	if l.parent == nil {
		l.parent = parent
	}
}

func (l *levelDBCache) Children() cache.ICache {
	return l.children
}

func (l *levelDBCache) SetChildren(children cache.ICache) {
	if l.children == nil {
		l.children = children
	}
}

func (l *levelDBCache) Close() error {
	return l.db.Close()
}
