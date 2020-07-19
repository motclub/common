package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/motclub/common/logging"
	"github.com/motclub/common/mq"
	"strings"
	"time"
)

func NewRedisMessage(opts *redis.UniversalOptions, logger logging.ILogger) (mq.IMessage, error) {
	rdb := redis.NewUniversalClient(opts)
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	if logger == nil {
		logger = logging.DefaultLogger
	}
	return &redisMessage{
		rdb:    rdb,
		logger: logger,
	}, nil
}

type redisMessage struct {
	rdb    redis.UniversalClient
	logger logging.ILogger
}

func (r *redisMessage) Logger() logging.ILogger {
	return r.logger
}

func (r *redisMessage) XAdd(topic string, values map[string]interface{}) error {
	res := r.rdb.XAdd(context.Background(), &redis.XAddArgs{
		Stream: topic,
		ID:     "*",
		Values: values,
	})
	return res.Err()
}

func (r *redisMessage) XRead(topics []string, callback func(string, *mq.XMessage) error) error {
	if len(topics) == 0 || callback == nil {
		return nil
	}
	var streams []string
	streams = append(streams, topics...)
	streams = append(streams, "$")
	go func() {
		for {
			// block here...
			cmd := r.rdb.XRead(context.Background(), &redis.XReadArgs{
				Streams: streams,
				Block:   0,
			})
			result, err := cmd.Result()
			if err != nil {
				r.logger.ERROR(err)
				continue
			}
			r.handleXStreams(result, "", callback)
		}
	}()
	return nil
}

func (r *redisMessage) XGroupRead(topics []string, group, consumer string, callback func(string, *mq.XMessage) error) error {
	if len(topics) == 0 || group == "" || consumer == "" || callback == nil {
		return nil
	}
	// 创建消费组
	for _, topic := range topics {
		gs, err := r.XInfoGroups(topic)
		if err != nil {
			return err
		}
		var exists bool
		for _, g := range gs {
			if g.Name == group {
				exists = true
				break
			}
		}
		if !exists {
			if err := r.XGroupCreate(topic, group, "$"); err != nil && strings.Contains(err.Error(), "already exists") {
				return err
			}
		}
	}
	// 读取消息
	var streams []string
	streams = append(streams, topics...)
	streams = append(streams, ">")
	go func() {
		for {
			// block here...
			cmd := r.rdb.XReadGroup(context.Background(), &redis.XReadGroupArgs{
				Group:    group,
				Consumer: consumer,
				Streams:  streams,
				Block:    0,
			})
			result, err := cmd.Result()
			if err != nil {
				r.logger.ERROR(err)
				return
			}
			r.handleXStreams(result, group, callback)
		}
	}()
	// 监听PEL
	go func() {
		minIdle := 30 * time.Minute // 自动接盘超过30分钟未处理的消息
		for {
			for _, topic := range topics {
				list, err := r.XGroupPending(topic, group, "-", "+", 10, "")
				if err != nil {
					r.logger.ERROR(err)
					continue
				}
				var claimIds []string
				for _, item := range list {
					if item.Consumer == consumer {
						continue
					}
					if item.Idle < minIdle {
						continue
					}
					claimIds = append(claimIds, item.ID)
				}
				if len(claimIds) > 0 {
					if err := r.XGroupClaim(topic, group, consumer, minIdle, claimIds...); err != nil {
						r.logger.ERROR(err)
					}
				}
			}
			time.Sleep(5 * time.Minute) // 每5分钟监听一次
		}
	}()
	return nil
}

func (r *redisMessage) handleXStreams(streams []redis.XStream, group string, callback func(string, *mq.XMessage) error) {
	var autoAckIDs = make(map[string][]string)
	for _, stream := range streams {
		for _, msg := range stream.Messages {
			err := callback(stream.Stream, &mq.XMessage{
				ID:     msg.ID,
				Values: msg.Values,
			})
			if err != nil {
				r.logger.ERROR(err)
			} else {
				autoAckIDs[stream.Stream] = append(autoAckIDs[stream.Stream], msg.ID)
			}
		}
	}
	if len(autoAckIDs) > 0 && group != "" {
		for topic, ids := range autoAckIDs {
			_ = r.XGroupAck(topic, group, ids...)
		}
	}
}

func (r *redisMessage) XGroupCreate(topic, group, start string) error {
	if start == "" {
		start = "$"
	}
	return r.rdb.XGroupCreateMkStream(context.Background(), topic, group, start).Err()
}

func (r *redisMessage) XGroupDelConsumer(stream, group, consumer string) error {
	return r.rdb.XGroupDelConsumer(context.Background(), stream, group, consumer).Err()
}

func (r *redisMessage) XGroupDestroy(stream, group string) error {
	return r.rdb.XGroupDestroy(context.Background(), stream, group).Err()
}

func (r *redisMessage) XGroupAck(topic, group string, ids ...string) error {
	return r.rdb.XAck(context.Background(), topic, group, ids...).Err()
}

func (r *redisMessage) XGroupPending(topic, group string, start string, end string, count int64, consumer string) (mq.XGroupPendingResult, error) {
	cmd, err := r.rdb.XPendingExt(context.Background(), &redis.XPendingExtArgs{
		Stream:   topic,
		Group:    group,
		Start:    start,
		End:      end,
		Count:    count,
		Consumer: consumer,
	}).Result()
	if err != nil {
		return nil, err
	}
	var result mq.XGroupPendingResult
	for _, item := range cmd {
		result = append(result, mq.XGroupPendingItem{
			ID:         item.ID,
			Consumer:   item.Consumer,
			Idle:       item.Idle,
			RetryCount: item.RetryCount,
		})
	}
	return result, nil
}

func (r *redisMessage) XGroupClaim(topic, group, consumer string, minIdle time.Duration, ids ...string) error {
	if topic == "" || group == "" || consumer == "" || len(ids) == 0 {
		return nil
	}
	return r.rdb.XClaim(context.Background(), &redis.XClaimArgs{
		Stream:   topic,
		Group:    group,
		Consumer: consumer,
		MinIdle:  minIdle,
		Messages: ids,
	}).Err()
}

func (r *redisMessage) XInfoGroups(topic string) (mq.XInfoGroupsResult, error) {
	result, err := r.rdb.XInfoGroups(context.Background(), topic).Result()
	if err != nil {
		return nil, err
	}
	var groups mq.XInfoGroupsResult
	for _, item := range result {
		groups = append(groups, mq.XInfoGroupsItem{
			Name:            item.Name,
			Consumers:       item.Consumers,
			Pending:         item.Pending,
			LastDeliveredID: item.LastDeliveredID,
		})
	}
	return groups, nil
}

func (r *redisMessage) XDel(topic string, ids ...string) error {
	return r.rdb.XDel(context.Background(), topic, ids...).Err()
}

func (r *redisMessage) XRange(topic, start, end string, count int64) ([]mq.XMessage, error) {
	var cmd *redis.XMessageSliceCmd
	if count > 0 {
		cmd = r.rdb.XRangeN(context.Background(), topic, start, end, count)
	} else {
		cmd = r.rdb.XRange(context.Background(), topic, start, end)
	}
	result, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	var messages []mq.XMessage
	for _, item := range result {
		messages = append(messages, mq.XMessage{
			ID:     item.ID,
			Values: item.Values,
		})
	}
	return messages, nil
}

func (r *redisMessage) XClose() error {
	if r.rdb != nil {
		return r.rdb.Close()
	}
	return nil
}
