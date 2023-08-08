package main

import (
	"context"
	"fmt"
	sqlserver_storage "github.com/storage-lock/go-sqlserver-storage"
)

func main() {

	// 数据库连接不是DSN的形式，就是一堆零散的属性，则依次设置，可以得到一个连接管理器
	host := "127.0.0.1"
	port := uint(1433)
	username := "root"
	passwd := "UeGqAm8CxYGldMDLoNNt"
	connectionManager := sqlserver_storage.NewSqlServerConnectionManager(host, port, username, passwd)

	// 然后从这个连接管理器创建SqlServer Storage
	options := sqlserver_storage.NewSqlServerStorageOptions().SetConnectionManager(connectionManager)
	storage, err := sqlserver_storage.NewSqlServerStorage(context.Background(), options)
	if err != nil {
		panic(err)
	}
	fmt.Println(storage.GetName())

}
