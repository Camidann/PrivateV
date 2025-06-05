package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.Static("/Videos", "./Videos")
	router.LoadHTMLFiles("src/index.html", "src/upload.html", "src/video.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.GET("/upload", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", gin.H{})
	})

	router.POST("/upload/charge", func(c *gin.Context) {
		file, err := c.FormFile("video")
		if err != nil {
			c.String(http.StatusBadRequest, "No se pudo obtener el archivo: %v", err)
			return
		}
		nombre := c.PostForm("nombre")
		if nombre == "" {
			c.String(http.StatusBadRequest, "Debe ingresar un nombre para el video")
			return
		}

		videoFileName := fmt.Sprintf("%s.mp4", nombre)
		videoPath := fmt.Sprintf("Videos/%s", videoFileName)

		if err := c.SaveUploadedFile(file, videoPath); err != nil {
			c.String(http.StatusInternalServerError, "No se pudo guardar el archivo: %v", err)
			return
		}

		newURL := fmt.Sprintf("/video/%s", videoFileName)
		c.Redirect(http.StatusFound, newURL)
	})

	// Ruta para mostrar el video subido
	router.GET("/video/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		c.HTML(http.StatusOK, "video.html", gin.H{
			"video": "/Videos/" + filename,
		})
	})

}
