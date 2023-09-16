package main

import (
	"context"
	"fmt"
	sqlserver_storage "github.com/storage-lock/go-sqlserver-storage"
)

func main() {

	// 使用一个DSN形式的数据库连接字符串创建ConnectionManager
	testDsn := "sqlserver://sa:UeGqAm8CxYGldMDLoNNt@127.0.0.1:1433?database=storage_lock_test&connection+timeout=30"
	connectionManager := sqlserver_storage.NewSqlServerConnectionManagerFromDsn(testDsn)

	// 然后从这个ConnectionManager创建SqlServer Storage
	options := sqlserver_storage.NewSqlServerStorageOptions().SetConnectionManager(connectionManager)
	storage, err := sqlserver_storage.NewSqlServerStorage(context.Background(), options)
	if err != nil {
		panic(err)
	}
	fmt.Println(storage.GetName())

}
