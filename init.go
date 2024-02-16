package main

import (
	"log"
	"time"

	"github.com/IBM/sarama"
)

var producer sarama.SyncProducer

func Init(address []string) error {
	err := initProducer(address)
	return err
}
func initProducer(address []string) error {
	// 配置
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second
	p, err := sarama.NewSyncProducer(address, config)
	producer = p

	if err != nil {
		log.Printf("new sync producer error : %s \n", err.Error())
		return err
	}
	return nil
}
