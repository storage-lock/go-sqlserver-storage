package sqlserver_storage

import (
	"database/sql"
	"github.com/storage-lock/go-storage"
)

type SqlServerStorageOptions struct {

	// 存放锁的表的名字
	TableName string

	// 用于获取数据库连接
	ConnectionManager storage.ConnectionManager[*sql.DB]
}

func NewSqlServerStorageOptions() *SqlServerStorageOptions {
	return &SqlServerStorageOptions{
		TableName: storage.DefaultStorageTableName,
	}
}

func (x *SqlServerStorageOptions) SetTableName(tableName string) *SqlServerStorageOptions {
	x.TableName = tableName
	return x
}

func (x *SqlServerStorageOptions) SetConnectionProvider(connectionManager storage.ConnectionManager[*sql.DB]) *SqlServerStorageOptions {
	x.ConnectionManager = connectionManager
	return x
}
