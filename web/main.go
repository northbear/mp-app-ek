package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const UserDataName = "AuthorizedUser"

type UserCredentials struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"username" binding:"required"`
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
	log.Printf("received credentials: { username: %s, password: %s }", uc.Username, uc.Password)
	session.Set("AuthorizedUser", uc.Username)
	session.Save()

	if uc.Username == "admin" {
		c.Redirect(http.StatusSeeOther, "/restricted")
	} else {
		c.String(http.StatusUnauthorized,
			fmt.Sprintf("Hi %s! You cannot get into this part!\n", uc.Username))
	}
}

func targetHandler(c *gin.Context) {
	session := sessions.Default(c)

	user := session.Get("AuthorizedUser")
	if user == "admin" {
		c.String(http.StatusOK, "Bingo Admin!!!\n")
	} else {
		c.String(http.StatusUnauthorized, "Who are you???\n")
	}
}

func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/",                rootHandler)
	r.GET("/login_form", loginFormHandler)
	r.GET("/restricted",    targetHandler)

	r.POST("/login", loginHandler)

	r.Run()
}
