package main

import (
	"os"
	"strings"
)

// Esta es la carga de videos para subir a index y a videos
// Se suben en upload y lo que hace es buscar en la carpeta Videos, el video y su nombre
// entiendo el codigo? si
// lo puedo leer? si (porfavor no me hagas sufrir mas)
func listVideos() []string {
	files, err := os.ReadDir("./Videos")
	if err != nil {
		return []string{}
	}
	var videos []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".mp4") {
			videos = append(videos, file.Name())
		}
	}
	return videos
}
