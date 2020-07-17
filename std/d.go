package std

import (
	"fmt"
	"github.com/motclub/common/json"
	"strconv"
	"time"
)

type D map[string]interface{}

func (d D) Has(path string) bool {
	return d.HasGet(path, nil)
}

func (d D) HasGet(path string, dst interface{}) bool {
	v, has := d[path]
	if has && dst != nil {
		if err := json.Copy(v, dst); err != nil {
			return false
		}
	}
	return has
}

func (d D) HasGetInt(path string) (int, bool) {
	v, has := d[path]

	switch i := v.(type) {
	case int:
		return i, has
	case int8:
		return int(i), has
	case int16:
		return int(i), has
	case int32:
		return int(i), has
	case int64:
		return int(i), has
	case uint:
		return int(i), has
	case uint8:
		return int(i), has
	case uint16:
		return int(i), has
	case uint32:
		return int(i), has
	case uint64:
		return int(i), has
	case float32:
		return int(i), has
	case float64:
		return int(i), has
	case string:
		vv, err := strconv.Atoi(i)
		return vv, err == nil
	case bool:
		var vv int
		if i {
			vv = 1
		}
		return vv, has
	default:
		return 0, has
	}
}

func (d D) HasGetInt8(path string) (int8, bool) {
	v, has := d.HasGetInt(path)
	return int8(v), has
}

func (d D) HasGetInt16(path string) (int16, bool) {
	v, has := d.HasGetInt(path)
	return int16(v), has
}

func (d D) HasGetInt32(path string) (int32, bool) {
	v, has := d.HasGetInt(path)
	return int32(v), has
}

func (d D) HasGetInt64(path string) (int64, bool) {
	v, has := d.HasGetInt(path)
	return int64(v), has
}

func (d D) HasGetUint(path string) (uint, bool) {
	v, has := d.HasGetInt(path)
	return uint(v), has
}

func (d D) HasGetUint8(path string) (uint8, bool) {
	v, has := d.HasGetInt(path)
	return uint8(v), has
}

func (d D) HasGetUint16(path string) (uint16, bool) {
	v, has := d.HasGetInt(path)
	return uint16(v), has
}

func (d D) HasGetUint32(path string) (uint32, bool) {
	v, has := d.HasGetInt(path)
	return uint32(v), has
}

func (d D) HasGetUint64(path string) (uint64, bool) {
	v, has := d.HasGetInt(path)
	return uint64(v), has
}

func (d D) HasGetFloat(path string) (float64, bool) {
	v, has := d[path]

	switch i := v.(type) {
	case int:
		return float64(i), has
	case int8:
		return float64(i), has
	case int16:
		return float64(i), has
	case int32:
		return float64(i), has
	case int64:
		return float64(i), has
	case uint:
		return float64(i), has
	case uint8:
		return float64(i), has
	case uint16:
		return float64(i), has
	case uint32:
		return float64(i), has
	case uint64:
		return float64(i), has
	case float32:
		return float64(i), has
	case float64:
		return i, has
	case string:
		vv, err := strconv.ParseFloat(i, 64)
		return vv, err == nil
	case bool:
		var vv float64
		if i {
			vv = 1
		}
		return vv, has
	default:
		return 0, has
	}
}

func (d D) HasGetFloat32(path string) (float32, bool) {
	v, has := d.HasGetFloat(path)
	return float32(v), has
}

func (d D) HasGetFloat64(path string) (float64, bool) {
	v, has := d.HasGetFloat(path)
	return v, has
}

func (d D) HasGetString(path string) (string, bool) {
	v, has := d[path]

	switch i := v.(type) {
	case int:
		return fmt.Sprintf("%v", i), has
	case int8:
		return fmt.Sprintf("%v", i), has
	case int16:
		return fmt.Sprintf("%v", i), has
	case int32:
		return fmt.Sprintf("%v", i), has
	case int64:
		return fmt.Sprintf("%v", i), has
	case uint:
		return fmt.Sprintf("%v", i), has
	case uint8:
		return fmt.Sprintf("%v", i), has
	case uint16:
		return fmt.Sprintf("%v", i), has
	case uint32:
		return fmt.Sprintf("%v", i), has
	case uint64:
		return fmt.Sprintf("%v", i), has
	case float32:
		return fmt.Sprintf("%v", i), has
	case float64:
		return fmt.Sprintf("%v", i), has
	case string:
		return i, has
	case bool:
		return fmt.Sprintf("%v", i), has
	default:
		return "", has
	}
}

func (d D) HasGetBool(path string) (bool, bool) {
	v, has := d[path]

	switch i := v.(type) {
	case bool:
		return i, has
	default:
		return false, has
	}
}

func (d D) HasGetTime(path string) (time.Time, bool) {
	v, has := d[path]

	switch i := v.(type) {
	case int64:
		return time.Unix(i, 0), has
	case time.Time:
		return i, has
	case *time.Time:
		return *i, has
	case string:
		t, err := time.Parse(time.RFC3339, i)
		return t, has && err == nil
	default:
		return time.Time{}, has
	}
}

func (d D) Get(path string, dst interface{}) {
	_ = d.HasGet(path, dst)
}

func (d D) GetInt(path string) int {
	v, _ := d.HasGetInt(path)
	return v
}

func (d D) GetInt8(path string) int8 {
	v, _ := d.HasGetInt8(path)
	return v
}

func (d D) GetInt16(path string) int16 {
	v, _ := d.HasGetInt16(path)
	return v
}

func (d D) GetInt32(path string) int32 {
	v, _ := d.HasGetInt32(path)
	return v
}

func (d D) GetInt64(path string) int64 {
	v, _ := d.HasGetInt64(path)
	return v
}

func (d D) GetUint(path string) uint {
	v, _ := d.HasGetUint(path)
	return v
}

func (d D) GetUint8(path string) uint8 {
	v, _ := d.HasGetUint8(path)
	return v
}

func (d D) GetUint16(path string) uint16 {
	v, _ := d.HasGetUint16(path)
	return v
}

func (d D) GetUint32(path string) uint32 {
	v, _ := d.HasGetUint32(path)
	return v
}

func (d D) GetUint64(path string) uint64 {
	v, _ := d.HasGetUint64(path)
	return v
}

func (d D) GetFloat(path string) float64 {
	v, _ := d.HasGetFloat(path)
	return v
}

func (d D) GetFloat32(path string) float32 {
	v, _ := d.HasGetFloat32(path)
	return v
}

func (d D) GetFloat64(path string) float64 {
	v, _ := d.HasGetFloat64(path)
	return v
}

func (d D) GetString(path string) string {
	v, _ := d.HasGetString(path)
	return v
}

func (d D) GetBool(path string) bool {
	v, _ := d.HasGetBool(path)
	return v
}

func (d D) GetTime(path string) time.Time {
	v, _ := d.HasGetTime(path)
	return v
}

func (d D) DefaultGet(path string, dst interface{}, defaultValue interface{}) {
	if !d.HasGet(path, dst) {
		_ = json.Copy(defaultValue, dst)
	}
}

func (d D) DefaultGetInt(path string, defaultValue int) int {
	if v, has := d.HasGetInt(path); has {
		return v
	}
	return defaultValue
}

func (d D) DefaultGetInt8(path string, defaultValue int8) int8 {
	if v, has := d.HasGetInt8(path); has {
		return v
	}
	return defaultValue
}

func (d D) DefaultGetInt16(path string, defaultValue int16) int16 {
	if v, has := d.HasGetInt16(path); has {
		return v
	}
	return defaultValue
}

func (d D) DefaultGetInt32(path string, defaultValue int32) int32 {
	if v, has := d.HasGetInt32(path); has {
		return v
	}
	return defaultValue
}

func (d D) DefaultGetInt64(path string, defaultValue int64) int64 {
	if v, has := d.HasGetInt64(path); has {
		return v
	}
	return defaultValue
}

func (d D) DefaultGetUint(path string, defaultValue uint) uint {
	if v, has := d.HasGetUint(path); has {
		return v
	}
	return defaultValue
}

func (d D) DefaultGetUint8(path string, defaultValue uint8) uint8 {
	if v, has := d.HasGetUint8(path); has {
		return v
	}
	return defaultValue
}

func (d D) DefaultGetUint16(path string, defaultValue uint16) uint16 {
	if v, has := d.HasGetUint16(path); has {
		return v
	}
	return defaultValue
}

func (d D) DefaultGetUint32(path string, defaultValue uint32) uint32 {
	if v, has := d.HasGetUint32(path); has {
		return v
	}
	return defaultValue
}

func (d D) DefaultGetUint64(path string, defaultValue uint64) uint64 {
	if v, has := d.HasGetUint64(path); has {
		return v
	}
	return defaultValue
}

func (d D) DefaultGetFloat(path string, defaultValue float64) float64 {
	if v, has := d.HasGetFloat(path); has {
		return v
	}
	return defaultValue
}

func (d D) DefaultGetFloat32(path string, defaultValue float32) float32 {
	if v, has := d.HasGetFloat32(path); has {
		return v
	}
	return defaultValue
}

func (d D) DefaultGetFloat64(path string, defaultValue float64) float64 {
	if v, has := d.HasGetFloat64(path); has {
		return v
	}
	return defaultValue
}

func (d D) DefaultGetString(path string, defaultValue string) string {
	if v, has := d.HasGetString(path); has {
		return v
	}
	return defaultValue
}

func (d D) DefaultGetBool(path string, defaultValue bool) bool {
	if v, has := d.HasGetBool(path); has {
		return v
	}
	return defaultValue
}

func (d D) DefaultGetTime(path string, defaultValue time.Time) time.Time {
	if v, has := d.HasGetTime(path); has {
		return v
	}
	return defaultValue
}

func (d D) BindTo(dst interface{}) error {
	return json.Copy(d, dst)
}
