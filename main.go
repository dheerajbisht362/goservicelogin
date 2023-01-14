package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

var users = []User{{ID: 1, Name: "Dhrj", Password: "1234", Email: "dhrj@gmail.com"}}

func userLogin(c *gin.Context) {
	var loginUser User
	if err := c.BindJSON(&loginUser); err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, v := range users {
		if v.Email == loginUser.Email {
			err := bcrypt.CompareHashAndPassword([]byte(v.Password), []byte(loginUser.Password))
			if err != nil {
				println(err)
				c.String(http.StatusUnauthorized, "Email and Password does not match")
				return
			}
			c.JSON(http.StatusOK, "Logged In Success")
			return
		}
	}
	c.String(http.StatusUnauthorized, "Email not found")
}

func userRegister(c *gin.Context) {
	var newUser User

	if err := c.BindJSON(&newUser); err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, v := range users {
		if v.Email == newUser.Email {
			c.String(http.StatusMethodNotAllowed, "Email already present")
			return
		}
	}

	newPass, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 4)
	if err != nil {
		fmt.Println(err)
		return
	}
	newUser.ID = len(users)
	newUser.Password = string(newPass)

	users = append(users, newUser)

	c.JSON(http.StatusCreated, newUser)

}

func main() {
	router := gin.Default()
	router.POST("/login", userLogin)
	router.POST("/register", userRegister)

	router.Run("localhost:8080")
}
