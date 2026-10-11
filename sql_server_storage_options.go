package sqlserver_storage

import (
	"database/sql"

	"github.com/storage-lock/go-storage"
)

// SqlServerStorageOptions 创建基于SqlServer作为存储的Storage的选项
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

func (x *SqlServerStorageOptions) SetConnectionManager(connManager storage.ConnectionManager[*sql.DB]) *SqlServerStorageOptions {
	x.ConnectionManager = connManager
	return x
}

func (x *SqlServerStorageOptions) Check() error {
	if x.ConnectionManager == nil {
		return ErrSqlServerStorageOptionsConnectionManagerNil
	}
	return nil
}
