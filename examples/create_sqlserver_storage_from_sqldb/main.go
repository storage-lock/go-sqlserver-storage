package main

import (
	"context"
	"database/sql"
	"fmt"
	sqlserver_storage "github.com/storage-lock/go-sqlserver-storage"
	storage "github.com/storage-lock/go-storage"
)

func main() {

	// 假设已经在其它地方初始化数据库连接得到了一个*sql.DB
	testDsn := "sqlserver://sa:UeGqAm8CxYGldMDLoNNt@127.0.0.1:1433?database=storage_lock_test&connection+timeout=30"
	db, err := sql.Open("mssql", testDsn)
	if err != nil {
		panic(err)
	}

	// 则可以从这个*sql.DB中创建一个SqlServer Storage
	connectionManager := storage.NewFixedSqlDBConnectionManager(db)
	options := sqlserver_storage.NewSqlServerStorageOptions().SetConnectionManager(connectionManager)
	storage, err := sqlserver_storage.NewSqlServerStorage(context.Background(), options)
	if err != nil {
		panic(err)
	}
	fmt.Println(storage.GetName())

}
