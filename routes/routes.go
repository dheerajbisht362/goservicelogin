package routes

import (
	"database/sql"
	"fmt"
	"goservice/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type userRegisterFields struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Gender   string `json:"gender"`
}

func UserRegister(data database.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var newUser userRegisterFields
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
		id, err := data.CreateUser(newUser.Name, newUser.Email, string(newPass), newUser.Gender)
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

func UserLogin(data database.Database) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var loginUser database.User
		response := make(map[string]string)

		if err := c.BindJSON(&loginUser); err != nil {
			fmt.Println(err)
			response["Error"] = err.Error()
			c.JSON(http.StatusBadRequest, response)
			return
		}

		ok, err := data.FindUser(loginUser.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("Error When Login", err)
				response["Error"] = "Email not found"
				c.JSON(http.StatusBadRequest, response)
				return
			}
			fmt.Println("Error When Login", err)
			response["Error"] = "Server Error"
			c.JSON(http.StatusBadRequest, response)
			return
		}
		fmt.Println("User Found", ok.Password)

		err = bcrypt.CompareHashAndPassword([]byte(ok.Password), []byte(loginUser.Password))
		if err != nil {
			fmt.Println("Error decrypting password", err)
			response["Error"] = "Email or password Invalid"
			c.JSON(http.StatusBadRequest, response)
			return
		}
		response["Message"] = "Successfully logged in"
		response["Name"] = ok.Name.String
		response["Gender"] = ok.Gender.String
		c.JSON(http.StatusOK, response)
	}
	return gin.HandlerFunc(fn)
}
