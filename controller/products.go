package controller

import (
	"MapReduce/infra/dao"
	model2 "MapReduce/model"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	mapreduce "github.com/charlieloom/MapReducde"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
)

func ExportProducts(c *gin.Context) {
	var condition model2.Condition
	condition.Id, _ = strconv.Atoi(c.Query("id"))
	condition.Name = c.Query("name")
	condition.Category = c.Query("category")
	condition.Page, _ = strconv.Atoi(c.Query("page"))
	condition.PageSize, _ = strconv.Atoi(c.Query("pageSize")) //总共要查询的数量
	condition.Sort = c.Query("sort")

	address := []string{"127.0.0.1:9092"}

	//初始化mapreduce
	mr, err := mapreduce.InitKMr(address)
	if err != nil {
		log.Panic(err)
		return
	}

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
		task := &mapreduce.Task{
			Param:    paramJson,
			Taskname: "query",
			Group:    "query",
		}

		mr.Map2(func() (*mapreduce.Task, error) {
			return task, nil
		})
	}
	fmt.Printf("map花费时间[%vs]", time.Since(start).Seconds())

	//reduce 查询
	go mr.Reduce([]string{"query"}, func(task *mapreduce.Task) error {
		paramquery := &model2.QueryMsg{}
		err = jsoniter.UnmarshalFromString(task.Param, paramquery)
		if err != nil {
			log.Panic(err)
			return err
		}
		//查询
		productlist, err := dao.GetAllproducts(&paramquery.Condition, paramquery.Offset, paramquery.Limit)
		if err != nil {
			fmt.Println(err)
			return err
		}
		value := model2.ExportMsg{
			Productlist: productlist,
			File:        paramquery.File,
			Row:         paramquery.Row,
		}

		paramJson, _ := jsoniter.MarshalToString(value)
		//定义任务
		task = &mapreduce.Task{
			Param:    paramJson,
			Taskname: "export",
			Group:    "export",
		}

		//将查询的结果发送出去
		mr.Map2(func() (*mapreduce.Task, error) {
			return task, nil
		})

		return nil
	})

	//导出
	go mr.Reduce([]string{"export"}, func(task *mapreduce.Task) error {
		paramexport := &model2.ExportMsg{}
		err = jsoniter.UnmarshalFromString(task.Param, paramexport)
		if err != nil {
			log.Panic(err)
			return err
		}

		//处理
		file, err := os.OpenFile(paramexport.File, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
		if err != nil {
			log.Println("open file error :", err, paramexport.File)
			return err
		}
		defer file.Close()
		for _, product := range paramexport.Productlist {
			c := []string{cast.ToString(product.ID), product.Name, product.Description, product.Category, cast.ToString(product.Price), cast.ToString(product.StockQuantity), product.CountryOfManufacture, cast.ToString(product.DateAdded), cast.ToString(product.LastUpdated), cast.ToString(product.UnitsSold), cast.ToString(product.NumberOfReviews), cast.ToString(product.AverageRating)}
			row := strings.Join(c, ",")
			row += "\n"
			_, err := file.WriteString(row)
			if err != nil {
				fmt.Println("write Error:", err)
				return err
			}
		}
		return nil
	})

	c.JSON(http.StatusOK, gin.H{
		"data": "ok",
		"time": time.Now().String(),
		"cost": time.Since(start).Seconds(),
	})
}
