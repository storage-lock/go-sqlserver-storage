package main

import (
	"context"
	"fmt"

	sqlserver_storage "github.com/storage-lock/go-sqlserver-storage"
)

func main() {
	// 使用一个DSN形式的数据库连接字符串创建ConnectionManager
	testDsn := "sqlserver://sa:UeGqAm8CxYGldMDLoNNt@127.0.0.1:1433?database=storage_lock_test&connection+timeout=30"
	storage, err := newSqlServerStorageFromDSN(context.Background(), testDsn)
	if err != nil {
		panic(err)
	}
	fmt.Println(storage.GetName())
}

func newSqlServerStorageFromDSN(ctx context.Context, dsn string) (*sqlserver_storage.SqlServerStorage, error) {
	connectionManager := sqlserver_storage.NewSqlServerConnectionManagerFromDsn(dsn)
	options := sqlserver_storage.NewSqlServerStorageOptions().SetConnectionManager(connectionManager)
	return sqlserver_storage.NewSqlServerStorage(ctx, options)
}
