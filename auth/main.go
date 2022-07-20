package main

import (
	"os"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserCredentials struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type UserDatabase map[string]string

func (ud UserDatabase) Update(k string, v string) {
	ud[k] = v
}

func (ud UserDatabase) Validate(k string, v string) bool {
	if secret, ok := ud[k]; ok {
		return secret == v
	}
	return false
}

var (
	PORT     = ":8080"
	DATABASE UserDatabase = make(UserDatabase)
)

func LoadUserDatabaseFromEnv(ud *UserDatabase) {
	envvars := os.Environ()
	for _, v := range envvars {
		if !strings.HasPrefix(v, "MP_APP_USER") {
			continue
		}
		varval := strings.SplitN(v, "=", 2);
		if len(varval) != 2 {
			log.Println("Error: wrong parameter format", varval)
			continue
		}
		creds := strings.SplitN(varval[1], ":", 2)
		if len(creds) != 2 {
			log.Println("Error: wrong value format:", varval[1])
			log.Println("Error: username and password should be splitted by semicolon (':')")
			continue
		}
		ud.Update(creds[0], creds[1])
	}
}

func ValidateHandler(c *gin.Context) {
	data := UserCredentials{}
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Println("ValidateHandler: incorrect request data. Cannot parse JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if DATABASE.Validate(data.Username, data.Password) {
		c.JSON(http.StatusOK, gin.H{ "status": "OK" })
		return
	}

	c.JSON(http.StatusNotFound, gin.H{ "status": "Not Found" })
}

func main() {
	if port := os.Getenv("PORT"); port != "" {
		PORT = port
		log.Println("The service connected to port " + PORT)
	}

	LoadUserDatabaseFromEnv(&DATABASE)
	for k, v := range DATABASE {
		log.Printf("User: %s, Password: %s", k, v)
	}

	r := gin.Default()
	r.POST("/validate", ValidateHandler)

	r.Run(PORT)
}
