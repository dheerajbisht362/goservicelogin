// ./database/database.go
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Database struct {
	SqlDb *sql.DB
}

var dbContext = context.Background()

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Gender   string `json:"gender"`
}

func (db Database) CreateUser(newUser User) (int64, error) {
	var err error

	err = db.SqlDb.PingContext(dbContext)
	if err != nil {
		fmt.Println("Error connecting to database", err)
		return -1, err
	}

	queryStatement := fmt.Sprintf("INSERT INTO user(name, email, password ) VALUES ('%s','%s', '%s')", newUser.Name, newUser.Email, newUser.Password)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.SqlDb.ExecContext(ctx, queryStatement)

	if err != nil {
		log.Println("Error when inserting query", err)
		return -1, err
	}

	log.Println(res)
	id, _ := res.LastInsertId()
	return id, nil
}
