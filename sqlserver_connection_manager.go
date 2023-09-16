package sqlserver_storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/storage-lock/go-storage"
	"sync"
)

const (
	DriverNameMsSql     = "mssql"
	DriverNameSqlServer = "sqlserver"
)

// SqlServerConnectionManager 负责管理与Sql Server服务器的连接
type SqlServerConnectionManager struct {

	// 主机的名字
	Host string

	// 主机的端口
	Port uint

	// 用户名
	User string

	// 密码
	Passwd string

	// DSN
	// Example: "sqlserver://sa:UeGqAm8CxYGldMDLoNNt@192.168.128.206:1433"
	DSN string

	driverName string

	// 初始化好的数据库实例
	db   *sql.DB
	err  error
	once sync.Once
}

var _ storage.ConnectionManager[*sql.DB] = &SqlServerConnectionManager{}

// NewSqlServerConnectionManagerFromDsn 从DSN创建SqlServer连接
func NewSqlServerConnectionManagerFromDsn(dsn string) *SqlServerConnectionManager {
	return &SqlServerConnectionManager{
		DSN: dsn,
	}
}

func NewSqlServerConnectionManagerFromSqlDb(db *sql.DB) *SqlServerConnectionManager {
	return &SqlServerConnectionManager{
		db: db,
	}
}

// NewSqlServerConnectionManager 从服务器属性创建数据库连接
func NewSqlServerConnectionManager(host string, port uint, user, passwd string) *SqlServerConnectionManager {
	return &SqlServerConnectionManager{
		Host:   host,
		Port:   port,
		User:   user,
		Passwd: passwd,
	}
}

func (x *SqlServerConnectionManager) SetDriverName(driverName string) *SqlServerConnectionManager {
	x.driverName = driverName
	return x
}

const SqlServerConnectionManagerName = "sql-server-connection-manager"

func (x *SqlServerConnectionManager) Name() string {
	return SqlServerConnectionManagerName
}

func (x *SqlServerConnectionManager) GetDSN() string {
	if x.DSN != "" {
		return x.DSN
	}
	return fmt.Sprintf("sqlserver://%s:%s@%s:%d", x.User, x.Passwd, x.Host, x.Port)
}

// Take 获取到数据库的连接
func (x *SqlServerConnectionManager) Take(ctx context.Context) (*sql.DB, error) {
	x.once.Do(func() {
		if x.db != nil {
			return
		}
		driverName := x.driverName
		if driverName == "" {
			driverName = DriverNameMsSql
		}
		db, err := sql.Open(driverName, x.GetDSN())
		if err != nil {
			x.err = err
			return
		}
		x.db = db
	})
	return x.db, x.err
}

func (x *SqlServerConnectionManager) Return(ctx context.Context, db *sql.DB) error {
	return nil
}

func (x *SqlServerConnectionManager) Shutdown(ctx context.Context) error {
	if x.db != nil {
		return x.db.Close()
	}
	return nil
}
