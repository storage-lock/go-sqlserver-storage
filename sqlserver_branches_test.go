package sqlserver_storage

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/storage-lock/go-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	"github.com/stretchr/testify/assert"
)

type takeErrorManager struct{}

func (takeErrorManager) Name() string { return "take-error" }
func (takeErrorManager) Take(context.Context) (*sql.DB, error) {
	return nil, errors.New("take failed")
}
func (takeErrorManager) Return(context.Context, *sql.DB) error { return nil }
func (takeErrorManager) Shutdown(context.Context) error        { return nil }

type returnErrorManager struct{ db *sql.DB }

func (returnErrorManager) Name() string { return "return-error" }
func (x returnErrorManager) Take(context.Context) (*sql.DB, error) {
	return x.db, nil
}
func (returnErrorManager) Return(context.Context, *sql.DB) error {
	return errors.New("return failed")
}
func (returnErrorManager) Shutdown(context.Context) error { return nil }

type defaultDriver struct{}

func (defaultDriver) Open(string) (driver.Conn, error) { return defaultConn{}, nil }

type defaultConn struct{}

func (defaultConn) Prepare(string) (driver.Stmt, error) { return nil, errors.New("unsupported") }
func (defaultConn) Close() error                        { return nil }
func (defaultConn) Begin() (driver.Tx, error)           { return nil, errors.New("unsupported") }

func init() {
	sql.Register("sqlserver-coverage", defaultDriver{})
}

func TestSqlServerConnectionManager(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mock.ExpectClose()

	fromDSN := NewSqlServerConnectionManagerFromDsn("sqlserver://u:p@h:1")
	assert.Equal(t, SqlServerConnectionManagerName, fromDSN.Name())
	assert.Equal(t, "sqlserver://u:p@h:1", fromDSN.GetDSN())
	fromDSN.open = func(string, string) (*sql.DB, error) { return nil, errors.New("open failed") }
	_, err = fromDSN.Take(context.Background())
	assert.EqualError(t, err, "open failed")
	_, err = fromDSN.Take(context.Background())
	assert.EqualError(t, err, "open failed")

	preset := &SqlServerConnectionManager{err: errors.New("preset")}
	_, err = preset.Take(context.Background())
	assert.EqualError(t, err, "preset")

	fromDB := NewSqlServerConnectionManagerFromSqlDb(db)
	got, err := fromDB.Take(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, db, got)
	assert.NoError(t, fromDB.Return(context.Background(), got))
	assert.NoError(t, fromDB.Shutdown(context.Background()))

	fromConfig := NewSqlServerConnectionManager("h", 1433, "u", "p").SetDriverName(DriverNameSqlServer)
	assert.Equal(t, DriverNameSqlServer, fromConfig.driverName)
	assert.Contains(t, fromConfig.GetDSN(), "sqlserver://u:p@h:1433")
	emptyDriver := NewSqlServerConnectionManager("h", 1433, "u", "p")
	assert.Contains(t, emptyDriver.GetDSN(), "sqlserver://")
	fromConfig.open = func(string, string) (*sql.DB, error) { return db, nil }
	got, err = fromConfig.Take(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, db, got)
	assert.NoError(t, (&SqlServerConnectionManager{}).Shutdown(context.Background()))

	defaultManager := NewSqlServerConnectionManagerFromDsn("sqlserver://u:p@127.0.0.1:1").SetDriverName("sqlserver-coverage")
	opened, err := defaultManager.Take(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, opened)
	assert.NoError(t, defaultManager.Shutdown(context.Background()))
}

func TestNewSqlServerStorage_Branches(t *testing.T) {
	_, err := NewSqlServerStorage(context.Background(), NewSqlServerStorageOptions())
	assert.ErrorIs(t, err, ErrSqlServerStorageOptionsConnectionManagerNil)

	_, err = NewSqlServerStorage(context.Background(), &SqlServerStorageOptions{ConnectionManager: takeErrorManager{}, TableName: "locks"})
	assert.EqualError(t, err, "take failed")

	dbEmpty, mockEmpty, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbEmpty.Close()
	mockEmpty.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(0, 0))
	emptyTable, err := NewSqlServerStorage(context.Background(), &SqlServerStorageOptions{
		ConnectionManager: storage.NewFixedSqlDBConnectionManager(dbEmpty),
		TableName:         "",
	})
	assert.NoError(t, err)
	assert.Equal(t, storage.DefaultStorageTableName, emptyTable.tableFullName)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectExec("CREATE TABLE").WillReturnError(errors.New("create failed"))
	mock.ExpectClose()
	_, err = NewSqlServerStorage(context.Background(), &SqlServerStorageOptions{
		ConnectionManager: returnErrorManager{db: db},
		TableName:         "locks",
	})
	assert.Error(t, err)

	db2, mock2, err := sqlmock.New()
	assert.NoError(t, err)
	defer db2.Close()
	mock2.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(0, 0))
	s, err := NewSqlServerStorage(context.Background(), &SqlServerStorageOptions{
		ConnectionManager: storage.NewFixedSqlDBConnectionManager(db2),
	})
	assert.NoError(t, err)
	assert.Equal(t, storage.DefaultStorageTableName, s.tableFullName)

	info := &storage.LockInformation{OwnerId: "owner"}
	mock2.ExpectExec("INSERT INTO").WillReturnError(errors.New("Violation of PRIMARY KEY constraint"))
	assert.ErrorIs(t, s.CreateWithVersion(context.Background(), "lock", 1, info), storage_lock.ErrVersionMiss)
	mock2.ExpectExec("INSERT INTO").WillReturnError(errors.New("other"))
	assert.EqualError(t, s.CreateWithVersion(context.Background(), "lock", 1, info), "other")
	mock2.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
	assert.NoError(t, s.CreateWithVersion(context.Background(), "lock", 1, info))

	_, err = s.GetTime(context.Background())
	assert.Error(t, err)
	mock2.ExpectQuery("SELECT CURRENT_TIMESTAMP").WillReturnError(errors.New("query failed"))
	_, err = s.GetTime(context.Background())
	assert.EqualError(t, err, "query failed")
	mock2.ExpectQuery("SELECT CURRENT_TIMESTAMP").WillReturnRows(sqlmock.NewRows([]string{"now"}))
	_, err = s.GetTime(context.Background())
	assert.EqualError(t, err, "rs server time failed")
	now := time.Now()
	mock2.ExpectQuery("SELECT CURRENT_TIMESTAMP").WillReturnRows(sqlmock.NewRows([]string{"now"}).AddRow(now))
	got, err := s.GetTime(context.Background())
	assert.NoError(t, err)
	assert.False(t, got.IsZero())
	mock2.ExpectQuery("SELECT CURRENT_TIMESTAMP").WillReturnRows(sqlmock.NewRows([]string{"now"}).AddRow("bad"))
	_, err = s.GetTime(context.Background())
	assert.Error(t, err)

	takeFailed := &SqlServerStorage{options: &SqlServerStorageOptions{ConnectionManager: takeErrorManager{}}}
	_, err = takeFailed.GetTime(context.Background())
	assert.EqualError(t, err, "take failed")

	_, err = (&SqlServerStorage{}).Watch(context.Background(), "")
	assert.Error(t, err)
}
