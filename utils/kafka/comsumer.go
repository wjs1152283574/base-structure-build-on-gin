package kafka

import (
	"context"
	"fmt"
	"goweb/utils/parsecfg"
	"time"

	"github.com/Shopify/sarama"
)

var (
	group        sarama.ConsumerGroup
	KafkaBrokers = []string{}
	KafkaTopic   = ""
	GroupID      = ""
)

func init() {
	if parsecfg.GlobalConfig.Env == "dev" {
		KafkaBrokers = []string{parsecfg.GlobalConfig.Kafka.Dev.Host + parsecfg.GlobalConfig.Kafka.Dev.Port}
		KafkaTopic = parsecfg.GlobalConfig.Kafka.Dev.Topic // 本节点用户上线均以这个值注册到redis,同时也监听这个topic
		GroupID = parsecfg.GlobalConfig.Kafka.Dev.Topic
	}
	if parsecfg.GlobalConfig.Env == "prod" {
		KafkaBrokers = []string{parsecfg.GlobalConfig.Kafka.Prod.Host + parsecfg.GlobalConfig.Kafka.Prod.Port}
		KafkaTopic = parsecfg.GlobalConfig.Kafka.Prod.Topic // 本节点用户上线均以这个值注册到redis,同时也监听这个topic
		GroupID = parsecfg.GlobalConfig.Kafka.Prod.Topic
	}
	if parsecfg.GlobalConfig.Env == "stage" {
		KafkaBrokers = []string{parsecfg.GlobalConfig.Kafka.Stage.Host + parsecfg.GlobalConfig.Kafka.Stage.Port}
		KafkaTopic = parsecfg.GlobalConfig.Kafka.Stage.Topic // 本节点用户上线均以这个值注册到redis,同时也监听这个topic
		GroupID = parsecfg.GlobalConfig.Kafka.Stage.Topic
	}
}

type ConsumerGroupHandler struct{}

// 实现 Comsumer handler 接口
func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// 不要将一下代码放入另外的goroutine,因为调用这个函数本身已经是由一个新拉起goroutine来执行
	for msg := range claim.Messages() {
		fmt.Printf("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
		sess.MarkMessage(msg, "used") // 标记为已消费
	}
	return nil
}

func init() {
	var err error
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Offsets.Initial = sarama.OffsetOldest // 设为最新
	config.Consumer.Offsets.CommitInterval = 1 * time.Second
	group, err = sarama.NewConsumerGroup(KafkaBrokers, GroupID, config)
	if err != nil {
		panic(err)
	}
}

// Comsumer  持续存在的goroutine,持续监听topic，消费消息；目前主要进行消息转发
func Comsumer() {
	ctx := context.Background()
	for {
		topics := []string{KafkaTopic}
		handler := ConsumerGroupHandler{}

		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		err := group.Consume(ctx, topics, handler)
		if err != nil {
			panic(err)
		}
	}
}
func Comsumer2() {
	ctx := context.Background()
	for {
		topics := []string{KafkaTopic}
		handler := ConsumerGroupHandler{}

		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		err := group.Consume(ctx, topics, handler)
		if err != nil {
			panic(err)
		}
	}
}
