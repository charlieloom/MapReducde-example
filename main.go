package main

import (
	"MapReduce/infra/dao"
	"MapReduce/routers"
	"os"
)

func main() {
	//初始化
	dao.Init()
	port := os.Args[1]
	r := routers.SetupRouter()
	r.Run(":" + port)
}
