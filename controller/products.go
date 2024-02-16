package controller

import (
	mapreduce "MapReduce/infra/mapReduce"
	model2 "MapReduce/model"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	xql_mapreduce "github.com/charlieloom/MapReducde"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

func ExportProducts(c *gin.Context) {
	var condition model2.Condition
	condition.Id, _ = strconv.Atoi(c.Query("id"))
	condition.Name = c.Query("name")
	condition.Category = c.Query("category")
	condition.Page, _ = strconv.Atoi(c.Query("page"))
	condition.PageSize, _ = strconv.Atoi(c.Query("pageSize")) //总共要查询的数量
	condition.Sort = c.Query("sort")

	m := mapreduce.GetMapReduce()

	log.Println("map开始")
	start := time.Now()
	offest := (condition.Page - 1) * condition.PageSize //初始偏移
	for i := 0; i < condition.PageSize; i += 200 {
		value := model2.QueryMsg{
			Condition: condition,
			Offset:    i + offest, //当前偏移
			Limit:     200,
			Row:       i,
			File:      "products.csv",
		}
		paramJson, err := jsoniter.MarshalToString(value)
		if err != nil {
			log.Panic(err)
			return
		}
		//定义任务
		task := &xql_mapreduce.Task{
			Param:    paramJson,
			Taskname: "query",
			Group:    "query",
		}

		m.Map2(func() (*xql_mapreduce.Task, error) {
			return task, nil
		})
	}
	fmt.Printf("map花费时间[%vs]", time.Since(start).Seconds())

	c.JSON(http.StatusOK, gin.H{
		"data": "ok",
		"time": time.Now().String(),
		"cost": time.Since(start).Seconds(),
	})
}
