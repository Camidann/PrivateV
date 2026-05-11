package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// crea el path de cada video nuevo
// Amo esta funcion.

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
