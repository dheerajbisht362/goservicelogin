// ./database/database.go
package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
)

type Database struct {
	SqlDb *sql.DB
}

var dbContext = context.Background()

type User struct {
	UserID   int            `json:"id"`
	Name     sql.NullString `json:"name"`
	Password string         `json:"password"`
	Email    string         `json:"email"`
	Gender   sql.NullString `json:"gender"`
}

func GetDatabaseConnection() Database {
	var (
		password     = os.Getenv("MSSQL_DB_PASSWORD")
		user         = os.Getenv("MSSQL_DB_USER")
		port         = os.Getenv("MSSQL_DB_PORT")
		databaseUsed = os.Getenv("MSSQL_DB_DATABASE")
		host         = os.Getenv("MYSQL_DB_HOST")
	)

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, databaseUsed)

	sqlObj, connectionError := sql.Open("mysql", connectionString)
	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
	}

	data := Database{
		SqlDb: sqlObj,
	}
	return data
}
