package sqlserver_storage

import "fmt"

var (
	ErrSqlServerStorageOptionsConnectionManagerNil = fmt.Errorf("SqlServerStorageOptions.ConnectionManager can not nil")
)
