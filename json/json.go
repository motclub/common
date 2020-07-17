package json

import (
	"github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

func init() {
	extra.SupportPrivateFields()
	//extra.SetNamingStrategy(extra.LowerCaseWithUnderscores)
}

func STD() jsoniter.API {
	return jsoniter.ConfigCompatibleWithStandardLibrary
}

func Copy(src, dst interface{}) (err error) {
	var data []byte
	if data, err = STD().Marshal(src); err == nil {
		err = STD().Unmarshal(data, dst)
	}
	return
}

// Stringify
func Stringify(value interface{}, format bool) string {
	var data []byte
	if format {
		data, _ = STD().MarshalIndent(value, "", "  ")
	} else {
		data, _ = STD().Marshal(value)
	}
	return string(data)
}

// Parse
func Parse(v string, r interface{}) error {
	return STD().Unmarshal([]byte(v), r)
}
