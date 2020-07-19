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

var DefaultLogger = &logger{}

type ILogger interface {
	PRINT(v ...interface{})
	DEBUG(v ...interface{})
	WARN(v ...interface{})
	INFO(v ...interface{})
	ERROR(v ...interface{})
	FATAL(v ...interface{})
	PANIC(v ...interface{})
}

type logger struct{}

func (d *logger) parseArgs(v []interface{}) (fields logrus.Fields, format string, args []interface{}, err error) {
	if len(v) == 0 || (len(v) == 1 && reflectx.IsNil(v[0])) {
		return
	}
	if vv, ok := v[0].(string); ok {
		format = vv
		if len(v) > 1 {
			args = v[1:]
		}
	} else {
		indirectValue := reflect.Indirect(reflect.ValueOf(v[0]))
		switch indirectValue.Kind() {
		case reflect.Map, reflect.Struct:
			_ = json.Copy(v[0], &fields)
		default:
			fields = make(logrus.Fields)
			fields["data"] = fmt.Sprintf("%v", v[0])
		}
		if len(v) > 1 {
			if vv, ok := v[1].(string); ok {
				format = vv
				if len(v) > 2 {
					args = v[2:]
				}
			}
		}
	}
	if len(fields) == 0 && format == "" {
		return fields, format, args, errors.New("unsupported parameter format")
	}
	return
}

func (d *logger) PRINT(v ...interface{}) {
	fields, format, args, err := d.parseArgs(v)
	if err != nil {
		return
	}
	logrus.WithFields(fields).Printf(format, args...)
}

func (d *logger) DEBUG(v ...interface{}) {
	fields, format, args, err := d.parseArgs(v)
	if err != nil {
		return
	}
	logrus.WithFields(fields).Debugf(format, args...)
}

func (d *logger) WARN(v ...interface{}) {
	fields, format, args, err := d.parseArgs(v)
	if err != nil {
		return
	}
	logrus.WithFields(fields).Warnf(format, args...)
}

func (d *logger) INFO(v ...interface{}) {
	fields, format, args, err := d.parseArgs(v)
	if err != nil {
		return
	}
	logrus.WithFields(fields).Infof(format, args...)
}

func (d *logger) ERROR(v ...interface{}) {
	fields, format, args, err := d.parseArgs(v)
	if err != nil {
		return
	}
	logrus.WithFields(fields).Errorf(format, args...)
}

func (d *logger) FATAL(v ...interface{}) {
	fields, format, args, err := d.parseArgs(v)
	if err != nil {
		return
	}
	logrus.WithFields(fields).Fatalf(format, args...)
}

func (d *logger) PANIC(v ...interface{}) {
	fields, format, args, err := d.parseArgs(v)
	if err != nil {
		return
	}
	logrus.WithFields(fields).Panicf(format, args...)
}

func PRINT(v ...interface{}) {
	DefaultLogger.PRINT(v...)
}

func DEBUG(v ...interface{}) {
	DefaultLogger.DEBUG(v...)
}

func WARN(v ...interface{}) {
	DefaultLogger.WARN(v...)
}

func INFO(v ...interface{}) {
	DefaultLogger.INFO(v...)
}

func ERROR(v ...interface{}) {
	DefaultLogger.ERROR(v...)
}

func FATAL(v ...interface{}) {
	DefaultLogger.FATAL(v...)
}

func PANIC(v ...interface{}) {
	DefaultLogger.PANIC(v...)
}
