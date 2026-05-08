package main

import (
	"log"
)

// Esta es la carga de videos para subir a index y a videos
// Se suben en upload y lo que hace es buscar en la DB, el video y su nombre
func listVideos() []string {
	rows, err := db.Query("SELECT filename FROM videos")
	if err != nil {
		log.Println("Error querying videos:", err)
		return []string{}
	}
	defer rows.Close()

	var videos []string
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			log.Println("Error scanning video:", err)
			continue
		}
		videos = append(videos, filename)
	}
	return videos
}
