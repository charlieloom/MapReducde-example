package dao

import (
	"MapReduce/dal"
	"MapReduce/dal/query"
)

const MySQLDSN = "root:123456@tcp(127.0.0.1:3307)/Shop?charset=utf8mb4&parseTime=True"

func Init() {
	dal.DB = dal.ConnectDB(MySQLDSN)
	if dal.DB == nil {
		panic("connect db fail")
	}
	query.SetDefault(dal.DB)
}
