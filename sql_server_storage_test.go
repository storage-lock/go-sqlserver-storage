package sqlserver_storage

import (
	"context"
	storage_test_helper "github.com/storage-lock/go-storage-test-helper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewSqlServerStorage(t *testing.T) {
	envName := "STORAGE_LOCK_SQLSERVER_DSN"
	dsn := os.Getenv(envName)
	assert.NotEmpty(t, dsn)
	connectionGetter := NewSqlServerStorageConnectionGetterFromDSN(dsn)
	storage, err := NewSqlServerStorage(context.Background(), &SqlServerStorageOptions{
		ConnectionManager: connectionGetter,
		TableName:         storage_test_helper.TestTableName,
	})
	assert.Nil(t, err)
	storage_test_helper.TestStorage(t, storage)
}
