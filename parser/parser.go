package parser

import (
	"github.com/buger/jsonparser"
	"github.com/motclub/common/getter"
	"github.com/motclub/common/json"
	"github.com/motclub/common/std"
	"github.com/pkg/errors"
	"regexp"
	"strings"
	"sync"
	"time"
)

func NewParser(v std.D) (getter.IGetter, error) {
	data, err := json.STD().Marshal(v)
	if err != nil {
		return nil, errors.Wrap(err, "parser")
	}
	return &parser{data: data}, nil
}

var (
	parsedPaths   = make(map[string][]string)
	parsedPathsMu sync.RWMutex
)

type parser struct {
	data []byte
}

func (c *parser) parsePath(path string) []string {
	parsedPathsMu.Lock()
	defer parsedPathsMu.Unlock()

	if v, has := parsedPaths[path]; has {
		return v
	}

	var reg *regexp.Regexp
	reg = regexp.MustCompile(`(\[\d+\])`)
	path = reg.ReplaceAllString(path, ".$1.")
	reg = regexp.MustCompile(`\.+`)
	path = reg.ReplaceAllString(path, ".")
	path = strings.Trim(path, ".")
	keys := strings.Split(path, ".")

	parsedPaths[path] = keys
	return keys
}

func (c *parser) Has(path string) bool {
	return c.HasGet(path, nil)
}

func (c *parser) HasGet(path string, dst interface{}) bool {
	keys := c.parsePath(path)
	data, _, _, err := jsonparser.Get(c.data, keys...)

	if err == nil && dst != nil {
		err = json.STD().Unmarshal(data, dst)
	}

	return err == nil
}

func (c *parser) HasGetInt(path string) (int, bool) {
	keys := c.parsePath(path)
	val, err := jsonparser.GetInt(c.data, keys...)

	return int(val), err == nil
}

func (c *parser) HasGetInt8(path string) (int8, bool) {
	v, has := c.HasGetInt(path)
	return int8(v), has
}

func (c *parser) HasGetInt16(path string) (int16, bool) {
	v, has := c.HasGetInt(path)
	return int16(v), has
}

func (c *parser) HasGetInt32(path string) (int32, bool) {
	v, has := c.HasGetInt(path)
	return int32(v), has
}

func (c *parser) HasGetInt64(path string) (int64, bool) {
	v, has := c.HasGetInt(path)
	return int64(v), has
}

func (c *parser) HasGetUint(path string) (uint, bool) {
	v, has := c.HasGetInt(path)
	return uint(v), has
}

func (c *parser) HasGetUint8(path string) (uint8, bool) {
	v, has := c.HasGetInt(path)
	return uint8(v), has
}

func (c *parser) HasGetUint16(path string) (uint16, bool) {
	v, has := c.HasGetInt(path)
	return uint16(v), has
}

func (c *parser) HasGetUint32(path string) (uint32, bool) {
	v, has := c.HasGetInt(path)
	return uint32(v), has
}

func (c *parser) HasGetUint64(path string) (uint64, bool) {
	v, has := c.HasGetInt(path)
	return uint64(v), has
}

func (c *parser) HasGetFloat(path string) (float64, bool) {
	keys := c.parsePath(path)
	val, err := jsonparser.GetFloat(c.data, keys...)

	return val, err == nil
}

func (c *parser) HasGetFloat32(path string) (float32, bool) {
	v, has := c.HasGetFloat(path)
	return float32(v), has
}

func (c *parser) HasGetFloat64(path string) (float64, bool) {
	return c.HasGetFloat(path)
}

func (c *parser) HasGetString(path string) (string, bool) {
	keys := c.parsePath(path)
	val, err := jsonparser.GetString(c.data, keys...)

	return val, err == nil
}

func (c *parser) HasGetBool(path string) (bool, bool) {
	keys := c.parsePath(path)
	val, err := jsonparser.GetBoolean(c.data, keys...)

	return val, err == nil
}

func (c *parser) HasGetTime(path string) (time.Time, bool) {
	var t time.Time
	has := c.HasGet(path, &t)
	return t, has
}

func (c *parser) Get(path string, dst interface{}) {
	c.HasGet(path, dst)
}

func (c *parser) GetInt(path string) int {
	v, _ := c.HasGetInt(path)
	return v
}

func (c *parser) GetInt8(path string) int8 {
	v, _ := c.HasGetInt8(path)
	return v
}

func (c *parser) GetInt16(path string) int16 {
	v, _ := c.HasGetInt16(path)
	return v
}

func (c *parser) GetInt32(path string) int32 {
	v, _ := c.HasGetInt32(path)
	return v
}

func (c *parser) GetInt64(path string) int64 {
	v, _ := c.HasGetInt64(path)
	return v
}

func (c *parser) GetUint(path string) uint {
	v, _ := c.HasGetUint(path)
	return v
}

func (c *parser) GetUint8(path string) uint8 {
	v, _ := c.HasGetUint8(path)
	return v
}

func (c *parser) GetUint16(path string) uint16 {
	v, _ := c.HasGetUint16(path)
	return v
}

func (c *parser) GetUint32(path string) uint32 {
	v, _ := c.HasGetUint32(path)
	return v
}

func (c *parser) GetUint64(path string) uint64 {
	v, _ := c.HasGetUint64(path)
	return v
}

func (c *parser) GetFloat(path string) float64 {
	v, _ := c.HasGetFloat(path)
	return v
}

func (c *parser) GetFloat32(path string) float32 {
	v, _ := c.HasGetFloat32(path)
	return v
}

func (c *parser) GetFloat64(path string) float64 {
	v, _ := c.HasGetFloat64(path)
	return v
}

func (c *parser) GetString(path string) string {
	v, _ := c.HasGetString(path)
	return v
}

func (c *parser) GetBool(path string) bool {
	v, _ := c.HasGetBool(path)
	return v
}

func (c *parser) GetTime(path string) time.Time {
	v, _ := c.HasGetTime(path)
	return v
}

func (c *parser) DefaultGet(path string, dst interface{}, defaultValue interface{}) {
	if has := c.HasGet(path, dst); !has {
		_ = json.Copy(defaultValue, dst)
	}
}

func (c *parser) DefaultGetInt(path string, defaultValue int) int {
	v, has := c.HasGetInt(path)
	if has {
		return v
	}
	return defaultValue
}

func (c *parser) DefaultGetInt8(path string, defaultValue int8) int8 {
	v, has := c.HasGetInt8(path)
	if has {
		return v
	}
	return defaultValue
}

func (c *parser) DefaultGetInt16(path string, defaultValue int16) int16 {
	v, has := c.HasGetInt16(path)
	if has {
		return v
	}
	return defaultValue
}

func (c *parser) DefaultGetInt32(path string, defaultValue int32) int32 {
	v, has := c.HasGetInt32(path)
	if has {
		return v
	}
	return defaultValue
}

func (c *parser) DefaultGetInt64(path string, defaultValue int64) int64 {
	v, has := c.HasGetInt64(path)
	if has {
		return v
	}
	return defaultValue
}

func (c *parser) DefaultGetUint(path string, defaultValue uint) uint {
	v, has := c.HasGetUint(path)
	if has {
		return v
	}
	return defaultValue
}

func (c *parser) DefaultGetUint8(path string, defaultValue uint8) uint8 {
	v, has := c.HasGetUint8(path)
	if has {
		return v
	}
	return defaultValue
}

func (c *parser) DefaultGetUint16(path string, defaultValue uint16) uint16 {
	v, has := c.HasGetUint16(path)
	if has {
		return v
	}
	return defaultValue
}

func (c *parser) DefaultGetUint32(path string, defaultValue uint32) uint32 {
	v, has := c.HasGetUint32(path)
	if has {
		return v
	}
	return defaultValue
}

func (c *parser) DefaultGetUint64(path string, defaultValue uint64) uint64 {
	v, has := c.HasGetUint64(path)
	if has {
		return v
	}
	return defaultValue
}

func (c *parser) DefaultGetFloat(path string, defaultValue float64) float64 {
	v, has := c.HasGetFloat(path)
	if has {
		return v
	}
	return defaultValue
}

func (c *parser) DefaultGetFloat32(path string, defaultValue float32) float32 {
	v, has := c.HasGetFloat32(path)
	if has {
		return v
	}
	return defaultValue
}

func (c *parser) DefaultGetFloat64(path string, defaultValue float64) float64 {
	v, has := c.HasGetFloat64(path)
	if has {
		return v
	}
	return defaultValue
}

func (c *parser) DefaultGetString(path string, defaultValue string) string {
	v, has := c.HasGetString(path)
	if has {
		return v
	}
	return defaultValue
}

func (c *parser) DefaultGetBool(path string, defaultValue bool) bool {
	v, has := c.HasGetBool(path)
	if has {
		return v
	}
	return defaultValue
}

func (c *parser) DefaultGetTime(path string, defaultValue time.Time) time.Time {
	v, has := c.HasGetTime(path)
	if has {
		return v
	}
	return defaultValue
}
