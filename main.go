package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	//Crea el Gin router
	initDB()
	router := gin.Default()
	fmt.Println("EMPEZO LA CARGA DE LA PAGINA")
	SetupRoutes(router)
	//Inicia el servidor en el puerto 8080
	router.Run(":8080")
}
