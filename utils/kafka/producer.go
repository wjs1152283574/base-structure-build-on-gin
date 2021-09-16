package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

//
var (
	producer sarama.SyncProducer
)

func init() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	brokers := KafkaBrokers
	var err error
	producer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		fmt.Printf("init producer failed -> %v \n", err)
		panic(err)
	} else {
		fmt.Println("producer init success")
	}
}

// Q:如何决定发到哪个topic 及设置多少个分区？
// A:kafka配置文件配置每个topic的partition数量，使用sarama提供的负载均衡机制分发到每个partition.

// KafkaSend 发送消息
func KafkaSend(msg, mobile, topic string) {

	msgX := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(mobile), // 相同的key 会分配到同一个partition,这里发送消息时同时传入用户电话作为分区依据
		Value: sarama.StringEncoder(msg),
	}
	partition, offset, err := producer.SendMessage(msgX)
	if err != nil {
		fmt.Printf("send msg error:%s \n", err)
	} else {
		fmt.Printf("msg send success, message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	}
}
