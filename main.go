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

	router.Run(":8080")
}
