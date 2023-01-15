package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"goservice/database"
	"log"
	"strconv"
	"time"
	"context"
	"database/sql"
	"os"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// func userLogin(data database.Database) gin.HandlerFunc {
// 	fn := func(c *gin.Context) {
// 		var loginUser User
// 		if err := c.BindJSON(&loginUser); err != nil {
// 			fmt.Println(err.Error())
// 			return
// 		}
// 		ok, err := data.CreateUser(database.User(loginUser))
// 		if err != nil {
// 			fmt.Println("Error When Login", err)
// 		}
// 		c.JSON(http.StatusOK, ok)
// 	}
// 	return gin.HandlerFunc(fn)
// }

func userRegister(data database.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var newUser database.User
		response := make(map[string]string)

		// Getting params from query
		if err := c.BindJSON(&newUser); err != nil {
			fmt.Println(err)
			response["Error"] = err.Error()
			c.JSON(http.StatusBadRequest, response)
			return
		}

		// validate object
		if newUser.Email == "" || newUser.Name == "" || newUser.Password == "" || newUser.Gender == "" {
			response["Error"] = "Invalid Properties"
			c.JSON(http.StatusBadRequest, response)
			return
		}

		//Encrypt user password
		newPass, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 4)
		if err != nil {
			fmt.Println(err)
			response["Error"] = err.Error()
			c.JSON(http.StatusBadRequest, err)
			return
		}
		newUser.Password = string(newPass)

		// database query
		id, err := data.CreateUser(database.User(newUser))
		if err != nil {
			if driverErr, ok := err.(*mysql.MySQLError); ok {
				if driverErr.Number == 1062 {
					response["Error"] = "Email already Exists"
					c.JSON(http.StatusBadRequest, response)
					return
				}
			}
			response["Error"] = "Server Error"
			c.JSON(http.StatusBadRequest, response)
			fmt.Println("Error registering the user", err)
			return
		}
		response["Message"] = "User Created"
		response["UserId"] = strconv.FormatInt(id, 10)
		fmt.Println(id)
		c.JSON(http.StatusOK, response)
	}

	return gin.HandlerFunc(fn)
}

func main() {

	envErr := godotenv.Load()
	if envErr != nil {
		fmt.Printf("Error loading credentials: %v", envErr)
	}

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
		return
	}

	data := database.Database{
		SqlDb: sqlObj,
	}

	// user table creation
	query := `CREATE TABLE IF NOT EXISTS user(user_id int primary key auto_increment, name text, email varchar(244) NOT NULL UNIQUE, password text, gender text)`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err := data.SqlDb.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating users table", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
	}
	log.Printf("Rows affected when creating table: %d", rows)

	// routes declaration
	router := gin.Default()
	// router.POST("/login", userLogin(data))
	router.POST("/register", userRegister(data))

	router.Run("localhost:8080")
}
