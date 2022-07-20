package main

import (
	"bytes"
	_ "encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const UserDataName = "AuthorizedUser"

var (
	PORT string = ":8080"
	AUTH_SERVICE string = "localhost:10080"
)

type UserCredentials struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"username" form:"password" binding:"required"`
}

func rootHandler(c *gin.Context) {
	msg := "Login page: <a href=\"login_form\">Login Page</a>\n"
	c.Data(http.StatusUnauthorized, "text/html; charset=utf-8", []byte(msg))
}

func loginFormHandler(c *gin.Context) {
	c.File("pages/login.html")
}

func loginHandler(c *gin.Context) {
	session := sessions.Default(c)
	uc := UserCredentials{}

	if c.ShouldBind(&uc) != nil {
		log.Println("Error: don't succeed to parse form data")
	}
	log.Printf("received credentials: { username: %s, password: ***** }", uc.Username, uc.Password)
	if QueryAuthService(uc) {
		log.Printf("The credentials get an approval", uc.Username, uc.Password)
		session.Set("AuthorizedUser", uc.Username)
		session.Save()
		c.Redirect(http.StatusSeeOther, "/restricted")
	} else {
		log.Printf("The credentials doesn't get an approval", uc.Username, uc.Password)
		c.String(http.StatusUnauthorized, "Hi! But... You are not authorized!")
	}
}

func targetHandler(c *gin.Context) {
	session := sessions.Default(c)

	user := session.Get("AuthorizedUser").(string)
	if user == "admin" {
		c.String(http.StatusOK, "Bingo Admin!!!\n")
	} else {
		c.String(http.StatusOK, "Hi " + user + ", Welcome to our service!\n")
	}
}

func QueryAuthService(uc UserCredentials) bool {
	query := "http://" + AUTH_SERVICE + "/validate"

	data := []byte(fmt.Sprintf("{ \"username\": \"%s\", \"password\": \"%s\" }", uc.Username, uc.Password))
	log.Printf("QueryAuthService: requested data: %s", uc)

	if r, err := http.Post(query, "application/json", bytes.NewReader(data)); err == nil {
		log.Printf("QueryAuthService: received response on %s: %s", uc.Username, r.Status)
		return r.StatusCode == http.StatusOK
	} else {
		log.Println("QueryAuthService: The request failed:", err)
		return false
	}
}

func main() {
	if port := os.Getenv("PORT"); port != "" {
		PORT = port
		log.Println("The service connected to port " + PORT)
	}
	if auth := os.Getenv("AUTH_SERVICE"); auth != "" {
		AUTH_SERVICE = auth
		log.Println("Auth service is awaited on " + AUTH_SERVICE)
	}

	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mp_session", store))

	r.GET("/", loginFormHandler)
	r.GET("/restricted",    targetHandler)

	r.POST("/login", loginHandler)

	r.Run(PORT)
}
