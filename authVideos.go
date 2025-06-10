package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		usuario, err := c.Cookie("usuario")
		if err != nil || usuario == "" {
			c.Redirect(http.StatusFound, "/sign-in")
			c.Abort()
			return
		}
		c.Next()
	}
}
