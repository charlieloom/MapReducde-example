package main

import (
	"MapReduce/consumer"
	"MapReduce/infra/dao"
	mapreduce "MapReduce/infra/mapReduce"
	"MapReduce/routers"
	"os"
)

func main() {
	//初始化
	dao.Init()
	mapreduce.Init()
	consumer.Init()

	port := os.Args[1]
	r := routers.SetupRouter()
	r.Run(":" + port)
}
