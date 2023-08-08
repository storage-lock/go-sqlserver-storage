# SqlServer Storage

# 一、这是什么
以SqlServer为存储引擎的[Storage](https://github.com/storage-lock/go-storage)实现，当前仓库为比较底层的存储层实现，你可以与[storage-lock](https://github.com/storage-lock/go-storage-lock)结合使用。

# 二、安装依赖
```bash
go get -u github.com/storage-lock/go-sqlserver-storage 
```

# 三、API Examples

## 3.1 从DSN创建SqlServerStorage

在Golang的世界中连接数据库最常见的就是DSN，下面的例子演示了如何从一个DSN创建SqlServerStorage：

```go
package main

import (
	"context"
	"fmt"
	sqlserver_storage "github.com/storage-lock/go-sqlserver-storage"
)

func main() {

	// 使用一个DSN形式的数据库连接字符串创建ConnectionManager
	testDsn := "sqlserver://sa:UeGqAm8CxYGldMDLoNNt@127.0.0.1:1433?database=storage_lock_test&connection+timeout=30"
	connectionManager := sqlserver_storage.NewSqlServerConnectionManagerFromDSN(testDsn)

	// 然后从这个ConnectionManager创建SqlServer Storage
	options := sqlserver_storage.NewSqlServerStorageOptions().SetConnectionManager(connectionManager)
	storage, err := sqlserver_storage.NewSqlServerStorage(context.Background(), options)
	if err != nil {
		panic(err)
	}
	fmt.Println(storage.GetName())

}
```

## 3.2 从连接属性（ip、端口、用户名、密码等等）中创建SqlServerStorage

或者你的配置文件中存放的并不是DSN，而是零散的几个连接属性，下面是一个创建SqlServerStorage的例子：

```go
package main

import (
	"context"
	"fmt"
	sqlserver_storage "github.com/storage-lock/go-sqlserver-storage"
)

func main() {

	// 数据库连接不是DSN的形式，就是一堆零散的属性，则依次设置，可以得到一个连接管理器
	host := "127.0.0.1"
	port := uint(1433)
	username := "root"
	passwd := "UeGqAm8CxYGldMDLoNNt"
	connectionManager := sqlserver_storage.NewSqlServerConnectionManager(host, port, username, passwd)

	// 然后从这个连接管理器创建SqlServer Storage
	options := sqlserver_storage.NewSqlServerStorageOptions().SetConnectionManager(connectionManager)
	storage, err := sqlserver_storage.NewSqlServerStorage(context.Background(), options)
	if err != nil {
		panic(err)
	}
	fmt.Println(storage.GetName())

}
```

## 3.3 从sql.DB创建SqlServerStorage

或者现在你已经有从其它渠道创建的能够连接到SqlServer的sql.DB，则也可以从这个*sql.DB创建SqlServerStorage：

```go
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
```









