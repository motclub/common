package mq

import (
	"github.com/motclub/common/logging"
	"time"
)

type XMessage struct {
	ID     string
	Values map[string]interface{}
}

type XGroupPendingItem struct {
	ID         string
	Consumer   string
	Idle       time.Duration
	RetryCount int64
}

type XGroupPendingResult []XGroupPendingItem

type XInfoGroupsItem struct {
	Name            string
	Consumers       int64
	Pending         int64
	LastDeliveredID string
}

type XInfoGroupsResult []XInfoGroupsItem

type IMessage interface {
	Logger() logging.ILogger
	XAdd(topic string, values map[string]interface{}) error
	XRead(topics []string, callback func(string, *XMessage) error) error

	XGroupRead(topics []string, group, consumer string, callback func(string, *XMessage) error) error
	XGroupCreate(topic, group, start string) error
	XGroupDelConsumer(stream, group, consumer string) error
	XGroupDestroy(stream, group string) error
	XGroupAck(topic, group string, ids ...string) error
	XGroupPending(topic, group string, start string, end string, count int64, consumer string) (XGroupPendingResult, error)
	XGroupClaim(topic, group, consumer string, minIdle time.Duration, ids ...string) error

	XInfoGroups(topic string) (XInfoGroupsResult, error)
	XDel(topic string, ids ...string) error
	XRange(topic, start, end string, count int64) ([]XMessage, error)
	XClose() error
}
