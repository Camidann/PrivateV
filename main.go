package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	//Crea el Gin router
	router := gin.Default()

	// Importa las rutas desde el archivo routes.go
	SetupRoutes(router)

	router.Run(":8080")
	// router.Run(":3000") for a hard coded port
}
