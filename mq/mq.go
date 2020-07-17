package mq

import "time"

type XMessage struct {
	ID     string
	Values map[string]interface{}
}

type XPendingArgs struct {
	Start    string
	End      string
	Count    int64
	Consumer string
}

type XPendingReplyEntry struct {
	ID         string
	Consumer   string
	Idle       time.Duration
	RetryCount int64
}

type IMessageQueue interface {
	XAddMessage(topic string, values map[string]interface{}) error
	XRegisterConsumer(topics []string, handler func(error, string, *XMessage) error)
	XRegisterConsumerGroup(topics []string, group, consumer string, handler func(error, string, *XMessage) error)
	XAck(topic, group string, ids ...string) error
	XDel(topic string, ids ...string) error
	XRange(topic, start, end string, count int64) ([]XMessage, error)
	XPending(topic, group string, args *XPendingArgs) ([]XPendingReplyEntry, error)
	XClaim(topic, group, consumer string, minIdle time.Duration, messages ...string) error

	XExclusiveRegisterConsumer(topics []string, consumer string, handler func(error, string, *XMessage) error)
	XExclusiveAck(topic string, ids ...string) error
	XExclusivePending(topic string, args *XPendingArgs) ([]XPendingReplyEntry, error)
	XExclusiveClaim(topic, consumer string, minIdle time.Duration, messages ...string) error

	XClose() error
}
