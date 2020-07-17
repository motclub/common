package getter

import "time"

type IGetter interface {
	Has(path string) bool
	HasGet(path string, dst interface{}) bool
	HasGetInt(path string) (int, bool)
	HasGetInt8(path string) (int8, bool)
	HasGetInt16(path string) (int16, bool)
	HasGetInt32(path string) (int32, bool)
	HasGetInt64(path string) (int64, bool)
	HasGetUint(path string) (uint, bool)
	HasGetUint8(path string) (uint8, bool)
	HasGetUint16(path string) (uint16, bool)
	HasGetUint32(path string) (uint32, bool)
	HasGetUint64(path string) (uint64, bool)
	HasGetFloat(path string) (float64, bool)
	HasGetFloat32(path string) (float32, bool)
	HasGetFloat64(path string) (float64, bool)
	HasGetString(path string) (string, bool)
	HasGetBool(path string) (bool, bool)
	HasGetTime(path string) (time.Time, bool)

	Get(path string, dst interface{})
	GetInt(path string) int
	GetInt8(path string) int8
	GetInt16(path string) int16
	GetInt32(path string) int32
	GetInt64(path string) int64
	GetUint(path string) uint
	GetUint8(path string) uint8
	GetUint16(path string) uint16
	GetUint32(path string) uint32
	GetUint64(path string) uint64
	GetFloat(path string) float64
	GetFloat32(path string) float32
	GetFloat64(path string) float64
	GetString(path string) string
	GetBool(path string) bool
	GetTime(path string) time.Time

	DefaultGet(path string, dst interface{}, defaultValue interface{})
	DefaultGetInt(path string, defaultValue int) int
	DefaultGetInt8(path string, defaultValue int8) int8
	DefaultGetInt16(path string, defaultValue int16) int16
	DefaultGetInt32(path string, defaultValue int32) int32
	DefaultGetInt64(path string, defaultValue int64) int64
	DefaultGetUint(path string, defaultValue uint) uint
	DefaultGetUint8(path string, defaultValue uint8) uint8
	DefaultGetUint16(path string, defaultValue uint16) uint16
	DefaultGetUint32(path string, defaultValue uint32) uint32
	DefaultGetUint64(path string, defaultValue uint64) uint64
	DefaultGetFloat(path string, defaultValue float64) float64
	DefaultGetFloat32(path string, defaultValue float32) float32
	DefaultGetFloat64(path string, defaultValue float64) float64
	DefaultGetString(path string, defaultValue string) string
	DefaultGetBool(path string, defaultValue bool) bool
	DefaultGetTime(path string, defaultValue time.Time) time.Time
}
