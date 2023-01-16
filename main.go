package main

import (
	"fmt"
	"goservice/database"
	"goservice/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {

	envErr := godotenv.Load()
	if envErr != nil {
		fmt.Printf("Error loading credentials: %v", envErr)
	}

	data := database.GetDatabaseConnection()

	_, err := database.CreateUsersTable(data)

	if err != nil {
		fmt.Printf("Error creating users table")
	}
	// routes declaration
	router := gin.Default()
	router.POST("/register", routes.UserRegister(data))
	router.POST("/login", routes.UserLogin(data))

	router.Run("localhost:8080")
}
