package main

import (
	"os"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserCredentials struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var (
	PORT     = ":8080"
	DATABASE = []UserCredentials{
		{ "admin",     "admin"    },
		{ "developer", "go0dc0de" },
	}
)

func ValidateHandler(c *gin.Context) {
	data := UserCredentials{}
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println("ValidateHandler: incorrect request data. Cannot parse JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("username: %s, password: %s\n", data.Username, data.Password)
	for _, uc := range DATABASE {
		if data.Username == uc.Username && data.Password == uc.Password {
			c.JSON(http.StatusOK, gin.H{ "status": "OK" })
			return 
		}
	}
	c.JSON(http.StatusNotFound, gin.H{ "status": "Not Found" })
}

func main() {
	if port := os.Getenv("PORT"); port != "" {
		PORT = port
		log.Println("The service connected to port " + PORT)
	}

	r := gin.Default()
	r.POST("/validate", ValidateHandler)

	r.Run(PORT)
}
