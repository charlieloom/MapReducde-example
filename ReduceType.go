package main

import (
	"log"

	"github.com/IBM/sarama"
	jsoniter "github.com/json-iterator/go"
)

type ReduceType struct {
	UserHandler func(task *Task) error
}

func (r ReduceType) Setup(sarama.ConsumerGroupSession) error {
	return nil
}
func (r ReduceType) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (r ReduceType) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		task := &Task{}
		err := jsoniter.UnmarshalFromString(string(msg.Value), task)
		if err != nil {
			log.Panic(err)
			return err
		}
		//用户定义UserHandler处理逻辑
		err = r.UserHandler(task)
		if err != nil {
			log.Panic(err)
			return err
		}
		session.MarkMessage(msg, "")
	}
	return nil
}
