package main

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	sqlserver_storage "github.com/storage-lock/go-sqlserver-storage"
	"github.com/stretchr/testify/assert"
)

func TestNewSqlServerStorageFromDSN_InitFailsWithoutDatabase(t *testing.T) {
	_, err := newSqlServerStorageFromDSN(context.Background(), "sqlserver://sa:pass@127.0.0.1:1?database=db")
	assert.Error(t, err)
}

func TestMain_PrintsStorageName(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(0, 0))
	orig := sqlserver_storage.SqlServerOpen
	sqlserver_storage.SqlServerOpen = func(string, string) (*sql.DB, error) { return db, nil }
	defer func() { sqlserver_storage.SqlServerOpen = orig }()
	assert.NotPanics(t, main)
}

func TestMain_PanicsWhenInitFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectExec(".").WillReturnError(assert.AnError)
	orig := sqlserver_storage.SqlServerOpen
	sqlserver_storage.SqlServerOpen = func(string, string) (*sql.DB, error) { return db, nil }
	defer func() { sqlserver_storage.SqlServerOpen = orig }()
	assert.Panics(t, main)
}
