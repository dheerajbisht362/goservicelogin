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
	UserID   int            `json:"id"`
	Name     sql.NullString `json:"name"`
	Password string         `json:"password"`
	Email    string         `json:"email"`
	Gender   sql.NullString `json:"gender"`
}

func (db Database) CreateUser(newUserName, newUserEmail, newUserPassword, newUserGender string) (int64, error) {

	err := db.SqlDb.PingContext(dbContext)
	if err != nil {
		fmt.Println("Error connecting to database", err)
		return -1, err
	}

	queryStatement := fmt.Sprintf("INSERT INTO user(name, email, password,gender) VALUES ('%s','%s', '%s','%s')", newUserName, newUserEmail, newUserPassword, newUserGender)

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

func (db Database) FindUser(userEmail string) (*User, error) {

	err := db.SqlDb.PingContext(dbContext)
	if err != nil {
		fmt.Println("Error connecting to database", err)
		return nil, err
	}
	fmt.Println(userEmail)

	queryStatement := fmt.Sprintf("SELECT user_id,gender,name,email,password FROM user WHERE email='%s'", userEmail)

	userMatch := &User{}
	if err := db.SqlDb.QueryRow(queryStatement).Scan(&userMatch.UserID, &userMatch.Gender, &userMatch.Name, &userMatch.Email,
		&userMatch.Password); err != nil {
		fmt.Println("Error when scaning data")
		return nil, err
	}
	return userMatch, nil
}
