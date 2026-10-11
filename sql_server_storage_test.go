package sqlserver_storage

import (
	"context"
	"database/sql"
	"testing"

	"github.com/storage-lock/go-storage"
	"github.com/stretchr/testify/assert"
)

func TestNewSqlServerStorageOptions(t *testing.T) {
	opts := NewSqlServerStorageOptions()
	assert.NotNil(t, opts)
	assert.Equal(t, storage.DefaultStorageTableName, opts.TableName)
	assert.Nil(t, opts.ConnectionManager)
}

func TestSqlServerStorageOptions_SetTableName(t *testing.T) {
	opts := NewSqlServerStorageOptions()
	assert.NotNil(t, opts)
	opts.SetTableName("custom_table")
	assert.Equal(t, "custom_table", opts.TableName)
}

func TestSqlServerStorageOptions_SetConnectionManager(t *testing.T) {
	opts := NewSqlServerStorageOptions()
	assert.NotNil(t, opts)
	db := &sql.DB{}
	opts.SetConnectionManager(NewSqlServerConnectionManagerFromSqlDb(db))
	assert.NotNil(t, opts.ConnectionManager)
}

func TestSqlServerStorageOptions_Check_NilConnectionManager(t *testing.T) {
	opts := NewSqlServerStorageOptions()
	assert.NotNil(t, opts)
	err := opts.Check()
	assert.Error(t, err)
	assert.Equal(t, ErrSqlServerStorageOptionsConnectionManagerNil, err)
}

func TestSqlServerStorageOptions_Check_ValidOptions(t *testing.T) {
	opts := NewSqlServerStorageOptions()
	assert.NotNil(t, opts)
	db := &sql.DB{}
	opts.SetConnectionManager(NewSqlServerConnectionManagerFromSqlDb(db))
	err := opts.Check()
	assert.NoError(t, err)
}

func TestNewSqlServerStorage_NilOptions(t *testing.T) {
	assert.Panics(t, func() {
		NewSqlServerStorage(context.Background(), nil)
	})
}

func TestSqlServerStorage_GetName(t *testing.T) {
	s := &SqlServerStorage{}
	name := s.GetName()
	assert.Equal(t, StorageName, name)
}

func TestSqlServerStorage_Capabilities(t *testing.T) {
	s := &SqlServerStorage{}
	caps := s.Capabilities()
	assert.NotNil(t, caps)
	assert.Len(t, caps, 3)
	assert.Equal(t, storage.CapabilityCAS, caps[0])
	assert.Equal(t, storage.CapabilityReliableTime, caps[1])
	assert.Equal(t, storage.CapabilityWatch, caps[2])
}
