package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/motclub/common/mq"
	"time"
)

const magicString = "0BE6C14B9954431FC3AEE05C4D4CF154"

func NewRedisMQ(opts *redis.UniversalOptions) (mq.IMessageQueue, error) {
	rdb := redis.NewUniversalClient(opts)
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	return &redisMQ{rdb: rdb}, nil
}

type redisMQ struct {
	rdb redis.UniversalClient
}

func (r *redisMQ) XAddMessage(topic string, values map[string]interface{}) error {
	res := r.rdb.XAdd(context.Background(), &redis.XAddArgs{
		Stream: topic,
		ID:     "*",
		Values: values,
	})
	return res.Err()
}

func (r *redisMQ) XRegisterConsumer(topics []string, handler func(error, string, *mq.XMessage) error) {
	if len(topics) == 0 || handler == nil {
		return
	}
	if topics[len(topics)-1] != "$" {
		topics[len(topics)-1] = "$"
	}
	go func() {
		for {
			res := r.rdb.XRead(context.Background(), &redis.XReadArgs{
				Streams: topics,
			})
			streams, err := res.Result()
			if err != nil && err != redis.Nil {
				_ = handler(err, "", nil)
				return
			}
			for _, stream := range streams {
				for _, msg := range stream.Messages {
					_ = handler(nil, stream.Stream, &mq.XMessage{
						ID:     msg.ID,
						Values: msg.Values,
					})
				}
			}
		}
	}()
}

func (r *redisMQ) XRegisterConsumerGroup(topics []string, group, consumer string, handler func(error, string, *mq.XMessage) error) {
	if len(topics) == 0 || group == "" || consumer == "" || handler == nil {
		return
	}
	if topics[len(topics)-1] != ">" {
		topics[len(topics)-1] = ">"
	}
	for _, topic := range topics {
		if topic == ">" {
			continue
		}
		r.rdb.XGroupCreate(context.Background(), topic, group, "0")
	}
	go func() {
		res := r.rdb.XReadGroup(context.Background(), &redis.XReadGroupArgs{
			Group:    group,
			Consumer: consumer,
			Streams:  topics,
		})
		streams, err := res.Result()
		if err != nil && err != redis.Nil {
			_ = handler(err, "", nil)
			return
		}
		var autoAckIDs = make(map[string][]string)
		for _, stream := range streams {
			for _, msg := range stream.Messages {
				err := handler(nil, stream.Stream, &mq.XMessage{
					ID:     msg.ID,
					Values: msg.Values,
				})
				if err == nil {
					autoAckIDs[stream.Stream] = append(autoAckIDs[stream.Stream], msg.ID)
				}
			}
		}
		if len(autoAckIDs) > 0 {
			for topic, ids := range autoAckIDs {
				_ = r.XAck(topic, group, ids...)
			}
		}
	}()
}

func (r *redisMQ) XAck(topic, group string, ids ...string) error {
	if topic == "" || group == "" || len(ids) == 0 {
		return nil
	}
	return r.rdb.XAck(context.Background(), topic, group, ids...).Err()
}

func (r *redisMQ) XDel(topic string, ids ...string) error {
	return r.rdb.XDel(context.Background(), topic, ids...).Err()
}

func (r *redisMQ) XRange(topic, start, end string, count int64) ([]mq.XMessage, error) {
	res, err := r.rdb.XRangeN(context.Background(), topic, start, end, count).Result()
	if err != nil {
		return nil, err
	}
	var ret []mq.XMessage
	for _, item := range res {
		ret = append(ret, mq.XMessage{
			ID:     item.ID,
			Values: item.Values,
		})
	}
	return ret, nil
}

func (r *redisMQ) XPending(topic, group string, args *mq.XPendingArgs) ([]mq.XPendingReplyEntry, error) {
	xArgs := &redis.XPendingExtArgs{
		Stream: topic,
		Group:  group,
	}
	if args != nil {
		xArgs.Start = args.Start
		xArgs.End = args.End
		xArgs.Count = args.Count
		xArgs.Consumer = args.Consumer
	}
	res, err := r.rdb.XPendingExt(context.Background(), xArgs).Result()
	if err != nil {
		return nil, err
	}
	var ret []mq.XPendingReplyEntry
	for _, item := range res {
		ret = append(ret, mq.XPendingReplyEntry{
			ID:         item.ID,
			Consumer:   item.Consumer,
			Idle:       item.Idle,
			RetryCount: item.RetryCount,
		})
	}
	return ret, nil
}

func (r *redisMQ) XClaim(topic, group, consumer string, minIdle time.Duration, messages ...string) error {
	res := r.rdb.XClaim(context.Background(), &redis.XClaimArgs{
		Stream:   topic,
		Group:    group,
		Consumer: consumer,
		MinIdle:  minIdle,
		Messages: messages,
	})
	return res.Err()
}

func (r *redisMQ) XExclusiveRegisterConsumer(topics []string, consumer string, handler func(error, string, *mq.XMessage) error) {
	r.XRegisterConsumerGroup(topics, magicString, consumer, handler)
}

func (r *redisMQ) XExclusiveAck(topic string, ids ...string) error {
	return r.XAck(topic, magicString, ids...)
}

func (r *redisMQ) XExclusivePending(topic string, args *mq.XPendingArgs) ([]mq.XPendingReplyEntry, error) {
	return r.XPending(topic, magicString, args)
}

func (r *redisMQ) XExclusiveClaim(topic, consumer string, minIdle time.Duration, messages ...string) error {
	return r.XClaim(topic, magicString, consumer, minIdle, messages...)
}

func (r *redisMQ) XClose() error {
	if r.rdb != nil {
		return r.rdb.Close()
	}
	return nil
}
