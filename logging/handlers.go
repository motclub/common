package logging

import (
	"github.com/sirupsen/logrus"
	"sync"
)

func init() {
	logrus.AddHook(&handlersHook{})
}

var (
	handlers   []func(*logrus.Entry) error
	handlersMu sync.RWMutex
)

func AddHandler(hs ...func(*logrus.Entry) error) {
	handlersMu.Lock()
	defer handlersMu.Unlock()

	handlers = append(handlers, hs...)
}

type handlersHook struct{}

func (h *handlersHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *handlersHook) Fire(entry *logrus.Entry) error {
	handlersMu.RLock()
	defer handlersMu.RUnlock()

	var err error
	for _, handler := range handlers {
		if e := handler(entry); e != nil {
			err = e
		}
	}
	return err
}
