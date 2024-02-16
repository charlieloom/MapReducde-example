package mapreduce

import (
	mapreduce "github.com/charlieloom/MapReducde"
)

var m *mapreduce.KafkaMapReduce

func Init() {
	mapreduce, err := mapreduce.InitKMr([]string{"123213"})
	if err != nil {
		panic(err)
	}
	m = &mapreduce
}

func GetMapReduce() *mapreduce.KafkaMapReduce {
	return m
}
