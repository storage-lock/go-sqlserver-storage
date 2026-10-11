package main

import (
	"context"
	"database/sql"
	"fmt"

	sqlserver_storage "github.com/storage-lock/go-sqlserver-storage"
	storage "github.com/storage-lock/go-storage"
)

var sqlOpen = sql.Open

func main() {
	// 假设已经在其它地方初始化数据库连接得到了一个*sql.DB
	testDsn := "sqlserver://sa:UeGqAm8CxYGldMDLoNNt@127.0.0.1:1433?database=storage_lock_test&connection+timeout=30"
	db, err := sqlOpen("mssql", testDsn)
	if err != nil {
		panic(err)
	}
	storage, err := newSqlServerStorageFromDB(context.Background(), db)
	if err != nil {
		panic(err)
	}
	fmt.Println(storage.GetName())
}

func newSqlServerStorageFromDB(ctx context.Context, db *sql.DB) (*sqlserver_storage.SqlServerStorage, error) {
	connectionManager := storage.NewFixedSqlDBConnectionManager(db)
	options := sqlserver_storage.NewSqlServerStorageOptions().SetConnectionManager(connectionManager)
	return sqlserver_storage.NewSqlServerStorage(ctx, options)
}
