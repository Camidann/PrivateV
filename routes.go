//Cuando empece este codigo, solo dios y yo sabian que tenia dentro. Ahora ninguno de los dos
//sabemos que mierda dice ni 1/4 del choclo este

package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// seleciona las rutas http
func SetupRoutes(router *gin.Engine) {
	router.Static("/Videos", "./Videos")
	router.LoadHTMLFiles("src/logout.html", "src/iniciodesession.html", "src/register.html", "src/index.html", "src/upload.html", "src/video.html")

	router.GET("/galeriaVideos", func(c *gin.Context) {
		videos := listVideos()
		c.HTML(http.StatusOK, "index.html", gin.H{
			"videos": videos,
		})
	})

	router.GET("/upload", AuthRequired(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", gin.H{})
	})

	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{})
	})

	router.POST("/register", func(c *gin.Context) {
		usuario := c.PostForm("usuario")
		password := c.PostForm("password")
		if usuario == "" || password == "" {
			c.String(http.StatusBadRequest, "Usuario y contraseña requeridos")
			return
		}
		_, err := db.Exec("INSERT INTO usuarios(usuario, password) VALUES (?, ?)", usuario, password)
		if err != nil {
			c.String(http.StatusInternalServerError, "No se pudo registrar el usuario")
			return
		}
		c.Redirect(http.StatusFound, "/sign-in")
	})

	router.GET("/sign-in", func(c *gin.Context) {
		c.HTML(http.StatusOK, "iniciodesession.html", gin.H{})
	})

	router.POST("/sign-in", func(c *gin.Context) {
		usuario := c.PostForm("usuario")
		password := c.PostForm("password")
		var dbUser, dbPass string
		err := db.QueryRow("SELECT usuario, password FROM usuarios WHERE usuario = ?", usuario).Scan(&dbUser, &dbPass)
		if err != nil || dbPass != password {
			c.String(http.StatusUnauthorized, "Usuario o contraseña incorrectos")
			return
		}
		// Guardar sesión en cookie
		c.SetCookie("usuario", usuario, 3600, "/", "", false, true)
		c.Redirect(http.StatusFound, "/galeriaVideos")
	})

	router.GET("/logout", AuthRequired(), func(c *gin.Context) {
		c.SetCookie("usuario", "", -1, "/", "", false, true)
		c.HTML(http.StatusOK, "logout.html", gin.H{})
	})

	//se encarga de la carga del video
	router.POST("/upload/charge", AuthRequired(), func(c *gin.Context) {
		file, err := c.FormFile("video")
		if err != nil {
			c.String(http.StatusBadRequest, "No se pudo obtener el archivo")
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
			c.String(http.StatusInternalServerError, "No se pudo guardar el archivo")
			return
		}

		_, err = db.Exec("INSERT OR REPLACE INTO videos(filename, titulo) VALUES (?, ?)", videoFileName, nombre)
		if err != nil {
			c.String(http.StatusInternalServerError, "No se pudo guardar en la base de datos")
			return
		}

		c.Redirect(http.StatusFound, fmt.Sprintf("/video/%s", videoFileName))
	})

	router.GET("/video/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		var titulo, descripcion string
		err := db.QueryRow("SELECT titulo, descripcion FROM videos WHERE filename = ?",
			filename).Scan(&titulo, &descripcion)
		if err != nil {
			titulo = filename
		}
		c.HTML(http.StatusOK, "video.html", gin.H{
			"video":  "/Videos/" + filename,
			"titulo": titulo,
		})
	})

}
