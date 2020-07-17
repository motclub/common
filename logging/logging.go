package logging

import (
	"errors"
	"fmt"
	"github.com/motclub/common/env"
	"github.com/motclub/common/json"
	"github.com/motclub/common/reflectx"
	"github.com/sirupsen/logrus"
	"reflect"
)

func init() {
	if env.Debug() {
		logrus.SetFormatter(&logrus.TextFormatter{})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}

func parseArgs(v []interface{}) (fields logrus.Fields, format string, args []interface{}, err error) {
	if len(v) == 0 || (len(v) == 1 && reflectx.IsNil(v[0])) {
		return
	}
	if vv, ok := v[0].(string); ok {
		format = vv
		if len(v) > 1 {
			args = v[1:]
		}
	} else if len(v) > 1 {
		indirectValue := reflect.Indirect(reflect.ValueOf(v[0]))
		switch indirectValue.Kind() {
		case reflect.Map, reflect.Struct:
			_ = json.Copy(v[0], &fields)
		default:
			fields = make(logrus.Fields)
			fields["data"] = fmt.Sprintf("%v", v[0])
		}
		if vv, ok := v[1].(string); ok {
			format = vv
			if len(v) > 2 {
				args = v[2:]
			}
		}
	}
	if len(fields) == 0 && format == "" {
		return fields, format, args, errors.New("unsupported parameter format")
	}
	return
}

func PRINT(v ...interface{}) {
	fields, format, args, err := parseArgs(v)
	if err != nil {
		return
	}
	logrus.WithFields(fields).Printf(format, args...)
}

func DEBUG(v ...interface{}) {
	fields, format, args, err := parseArgs(v)
	if err != nil {
		return
	}
	logrus.WithFields(fields).Debugf(format, args...)
}

func WARN(v ...interface{}) {
	fields, format, args, err := parseArgs(v)
	if err != nil {
		return
	}
	logrus.WithFields(fields).Warnf(format, args...)
}

func INFO(v ...interface{}) {
	fields, format, args, err := parseArgs(v)
	if err != nil {
		return
	}
	logrus.WithFields(fields).Infof(format, args...)
}

func ERROR(v ...interface{}) {
	fields, format, args, err := parseArgs(v)
	if err != nil {
		return
	}
	logrus.WithFields(fields).Errorf(format, args...)
}

func FATAL(v ...interface{}) {
	fields, format, args, err := parseArgs(v)
	if err != nil {
		return
	}
	logrus.WithFields(fields).Fatalf(format, args...)
}

func PANIC(v ...interface{}) {
	fields, format, args, err := parseArgs(v)
	if err != nil {
		return
	}
	logrus.WithFields(fields).Panicf(format, args...)
}
