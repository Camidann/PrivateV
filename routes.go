//Cuando empece este codigo, solo dios y yo sabian que tenia dentro.
//Ahora ninguno de los dos sabemos que mierda dice ni 1/4 del choclo este//

package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// seleciona las rutas http

//NO TOCAR NADA, FUNCIONA POR ARTE DE MAGIA, SI TOCAS ALGO SE ROMPE TODO, NO ME HAGAS SUFRIR MAS POR FAVOR

func SetupRoutes(router *gin.Engine) {
	// Quitar Static ya que ahora servimos desde DB
	// router.Static("/Videos", "./Videos")
	router.LoadHTMLFiles("src/logout.html", "src/iniciodesession.html", "src/register.html", "src/index.html", "src/upload.html", "src/video.html")

	router.GET("/galeriaVideos", func(c *gin.Context) {
		videos := listVideos()

		usuario, err := c.Cookie("usuario")
		logueado := err == nil && usuario != ""

		c.HTML(http.StatusOK, "index.html", gin.H{
			"videos":   videos,
			"logueado": logueado,
		})
	})

	router.GET("/upload", AuthRequired(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", gin.H{})
	})

	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{})
	})

	router.POST("/register", func(c *gin.Context) {
		usuario := strings.TrimSpace(c.PostForm("usuario"))
		password := strings.TrimSpace(c.PostForm("password"))
		if usuario == "" || password == "" {
			c.String(http.StatusBadRequest, "Usuario y contraseña requeridos")
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //encriptar la contraseña
		if err != nil {
			c.String(http.StatusInternalServerError, "No se pudo procesar la contraseña")
			return
		}

		_, err = db.Exec("INSERT INTO usuarios(usuario, password) VALUES (?, ?)", usuario, string(hash))
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE") { //si el error es por usuario duplicado
				c.String(http.StatusConflict, "El usuario ya existe")
				return
			}
			c.String(http.StatusInternalServerError, "No se pudo registrar el usuario")
			return
		}
		c.Redirect(http.StatusFound, "/sign-in")
	})

	router.GET("/sign-in", func(c *gin.Context) {
		c.HTML(http.StatusOK, "iniciodesession.html", gin.H{})
	})

	router.POST("/sign-in", func(c *gin.Context) {
		usuario := strings.TrimSpace(c.PostForm("usuario"))
		password := strings.TrimSpace(c.PostForm("password"))
		var dbUser, dbPass string
		err := db.QueryRow("SELECT usuario, password FROM usuarios WHERE usuario = ?", usuario).Scan(&dbUser, &dbPass)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(password)) != nil {
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

		// Leer el contenido del archivo
		src, err := file.Open()
		if err != nil {
			c.String(http.StatusInternalServerError, "No se pudo leer el archivo")
			return
		}
		defer src.Close()
		contenido, err := io.ReadAll(src)
		if err != nil {
			c.String(http.StatusInternalServerError, "No se pudo leer el contenido del archivo")
			return
		}

		// Obtener usuario_id
		usuario, err := c.Cookie("usuario")
		if err != nil {
			c.String(http.StatusUnauthorized, "Usuario no autenticado")
			return
		}
		usuarioID, err := getUsuarioID(usuario)
		if err != nil {
			c.String(http.StatusInternalServerError, "No se pudo obtener el usuario")
			return
		}

		_, err = db.Exec("INSERT OR REPLACE INTO videos(filename, titulo, contenido, usuario_id) VALUES (?, ?, ?, ?)", videoFileName, nombre, contenido, usuarioID)
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
			"video":  "/video-content/" + filename,
			"titulo": titulo,
		})
	})

	router.GET("/video-content/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		var contenido []byte
		err := db.QueryRow("SELECT contenido FROM videos WHERE filename = ?", filename).Scan(&contenido)
		if err != nil {
			c.String(http.StatusNotFound, "Video no encontrado")
			return
		}
		c.Header("Content-Type", "video/mp4")
		c.Data(http.StatusOK, "video/mp4", contenido)
	})

}
