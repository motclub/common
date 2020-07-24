package std

import (
	"errors"
	"fmt"
	"github.com/motclub/common/intl"
	"github.com/motclub/common/json"
	"github.com/motclub/common/reflectx"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrUnsupportedDataType = errors.New("unsupported data type")
)

type Args struct {
	RequestID      string      `json:"request_id"`
	SessionPayload D           `json:"session_payload"`
	Data           interface{} `json:"data"`
}

type Reply struct {
	Code           int                     `json:"code"`
	Data           interface{}             `json:"data"`
	Message        string                  `json:"msg,omitempty"`
	LocaleMessage  *intl.MessageDescriptor `json:"locale_message"`
	HTTPStatusCode int                     `json:"http_status_code"`
	HTTPAction     string                  `json:"http_action"`
	HTTPCookies    []*http.Cookie          `json:"http_cookies"`
}

func (s *Reply) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{"code": s.Code}
	if !reflectx.IsNil(s.Data) {
		data["data"] = s.Data
	}
	if s.Message != "" {
		data["msg"] = s.Message
	}

	return json.STD().Marshal(data)
}

func (s *Reply) Bind(dst interface{}) error {
	return json.Copy(s.Data, dst)
}

func (s *Reply) BindInt() (int, error) {
	v := s.Data

	switch i := v.(type) {
	case int:
		return i, nil
	case int8:
		return int(i), nil
	case int16:
		return int(i), nil
	case int32:
		return int(i), nil
	case int64:
		return int(i), nil
	case uint:
		return int(i), nil
	case uint8:
		return int(i), nil
	case uint16:
		return int(i), nil
	case uint32:
		return int(i), nil
	case uint64:
		return int(i), nil
	case float32:
		return int(i), nil
	case float64:
		return int(i), nil
	case string:
		vv, err := strconv.Atoi(i)
		return vv, err
	case bool:
		var vv int
		if i {
			vv = 1
		}
		return vv, nil
	default:
		return 0, ErrUnsupportedDataType
	}
}

func (s *Reply) BindInt8() (int8, error) {
	v, err := s.BindInt()
	return int8(v), err
}

func (s *Reply) BindInt16() (int16, error) {
	v, err := s.BindInt()
	return int16(v), err
}

func (s *Reply) BindInt32() (int32, error) {
	v, err := s.BindInt()
	return int32(v), err
}

func (s *Reply) BindInt64() (int64, error) {
	v, err := s.BindInt()
	return int64(v), err
}

func (s *Reply) BindUint() (uint, error) {
	v, err := s.BindInt()
	return uint(v), err
}

func (s *Reply) BindUint8() (uint8, error) {
	v, err := s.BindInt()
	return uint8(v), err
}

func (s *Reply) BindUint16() (uint16, error) {
	v, err := s.BindInt()
	return uint16(v), err
}

func (s *Reply) BindUint32() (uint32, error) {
	v, err := s.BindInt()
	return uint32(v), err
}

func (s *Reply) BindUint64() (uint64, error) {
	v, err := s.BindInt()
	return uint64(v), err
}

func (s *Reply) BindFloat() (float64, error) {
	v := s.Data

	switch i := v.(type) {
	case int:
		return float64(i), nil
	case int8:
		return float64(i), nil
	case int16:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case uint:
		return float64(i), nil
	case uint8:
		return float64(i), nil
	case uint16:
		return float64(i), nil
	case uint32:
		return float64(i), nil
	case uint64:
		return float64(i), nil
	case float32:
		return float64(i), nil
	case float64:
		return i, nil
	case string:
		vv, err := strconv.ParseFloat(i, 64)
		return vv, err
	case bool:
		var vv float64
		if i {
			vv = 1
		}
		return vv, nil
	default:
		return 0, ErrUnsupportedDataType
	}
}

func (s *Reply) BindFloat32() (float32, error) {
	v, err := s.BindFloat()
	return float32(v), err
}

func (s *Reply) BindFloat64() (float64, error) {
	return s.BindFloat()
}

func (s *Reply) BindString() (string, error) {
	v := s.Data

	switch i := v.(type) {
	case int:
		return fmt.Sprintf("%v", i), nil
	case int8:
		return fmt.Sprintf("%v", i), nil
	case int16:
		return fmt.Sprintf("%v", i), nil
	case int32:
		return fmt.Sprintf("%v", i), nil
	case int64:
		return fmt.Sprintf("%v", i), nil
	case uint:
		return fmt.Sprintf("%v", i), nil
	case uint8:
		return fmt.Sprintf("%v", i), nil
	case uint16:
		return fmt.Sprintf("%v", i), nil
	case uint32:
		return fmt.Sprintf("%v", i), nil
	case uint64:
		return fmt.Sprintf("%v", i), nil
	case float32:
		return fmt.Sprintf("%v", i), nil
	case float64:
		return fmt.Sprintf("%v", i), nil
	case string:
		return i, nil
	case bool:
		return fmt.Sprintf("%v", i), nil
	default:
		return "", nil
	}
}

func (s *Reply) BindBool() (bool, error) {
	v := s.Data

	switch i := v.(type) {
	case bool:
		return i, nil
	default:
		return false, ErrUnsupportedDataType
	}
}

func (s *Reply) BindTime() (time.Time, error) {
	v := s.Data

	switch i := v.(type) {
	case int64:
		return time.Unix(i, 0), nil
	case time.Time:
		return i, nil
	case *time.Time:
		return *i, nil
	case string:
		return time.Parse(time.RFC3339, i)
	default:
		return time.Time{}, ErrUnsupportedDataType
	}
}
func STD(data interface{}, args ...interface{}) *STDReply {
	std := std(data, args)
	return std
}

func std(data interface{}, args []interface{}) *Reply {
	reply := Reply{
		Data: data,
		Code: 0,
	}
	if localeMessage != nil {
		reply.Message = localeMessage.Render(c.AppCache())
	}
	return &reply
}

func stdErr(data interface{}, code int, args []interface{}) *Reply {
	messageName, defaultMessage, messageValues := resolveLocaleMessage(args)
	if code == 0 {
		code = -1
	}
	body := Reply{
		Data: data,
		Code: code,
	}
	if localeMessage != nil {
		body.Message = localeMessage.Render(c.AppCache())
	}
	return &body
}

func resolveLocaleMessage(args []interface{}) (messageID, defaultMessage string, contextValues interface{}) {
	if len(args) > 0 {
		messageID = args[0].(string)
	}
	if len(args) > 1 {
		defaultMessage = args[1].(string)
	}
	if len(args) > 2 {
		contextValues = args[2]
	}
	return
}
