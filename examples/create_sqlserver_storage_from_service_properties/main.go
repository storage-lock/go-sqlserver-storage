package main

import (
	"context"
	"fmt"

	sqlserver_storage "github.com/storage-lock/go-sqlserver-storage"
)

func main() {
	// 数据库连接不是DSN的形式，就是一堆零散的属性，则依次设置，可以得到一个连接管理器
	storage, err := newSqlServerStorageFromProperties(context.Background(), "127.0.0.1", 1433, "root", "UeGqAm8CxYGldMDLoNNt")
	if err != nil {
		panic(err)
	}
	fmt.Println(storage.GetName())
}

func newSqlServerStorageFromProperties(ctx context.Context, host string, port uint, username, passwd string) (*sqlserver_storage.SqlServerStorage, error) {
	connectionManager := sqlserver_storage.NewSqlServerConnectionManager(host, port, username, passwd)
	options := sqlserver_storage.NewSqlServerStorageOptions().SetConnectionManager(connectionManager)
	return sqlserver_storage.NewSqlServerStorage(ctx, options)
}
