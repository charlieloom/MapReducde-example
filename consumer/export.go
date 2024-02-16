package consumer

import (
	"MapReduce/infra/dao"
	mapreduce "MapReduce/infra/mapReduce"
	"fmt"
	"log"
	"os"
	"strings"

	model2 "MapReduce/model"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"

	xql_mapreduce "github.com/charlieloom/MapReducde"
)

func Query(task *xql_mapreduce.Task) error {
	m := mapreduce.GetMapReduce()

	paramquery := &model2.QueryMsg{}
	err := jsoniter.UnmarshalFromString(task.Param, paramquery)
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
	task = &xql_mapreduce.Task{
		Param:    paramJson,
		Taskname: "export",
		Group:    "export",
	}

	//将查询的结果发送出去
	m.Map2(func() (*xql_mapreduce.Task, error) {
		return task, nil
	})

	return nil
}

func Export(task *xql_mapreduce.Task) error {
	paramexport := &model2.ExportMsg{}
	err := jsoniter.UnmarshalFromString(task.Param, paramexport)
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
}
