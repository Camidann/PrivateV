package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var SubirVideo = false
var numeroAleatorio int // Declarar como variable global

func init() {
	rand.Seed(time.Now().UnixNano())
	numeroAleatorio = rand.Intn(9999999)
}

func SetupRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {

	})

	router.GET("/upload", func(c *gin.Context) {
		// Obtener la URL de la solicitud actual
		currentURL := c.Request.URL.String()
		fmt.Println("URL actual:", currentURL)

		// Construir una nueva URL a partir de la actual y agregar numeroAleatorio
		newURL := fmt.Sprintf("%s/%d", currentURL, numeroAleatorio)
		fmt.Println("URL construida:", newURL)

		// Redirigir a la nueva URL
		c.Redirect(http.StatusFound, newURL)
	})

}
