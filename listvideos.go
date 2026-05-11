package main

import "log"

type VideoItem struct {
	FileName string
	IsVideo  bool
}

func listVideos() []VideoItem {
	rows, err := db.Query("SELECT filename FROM videos")
	if err != nil {
		log.Println("Error querying videos:", err)
		return []VideoItem{}
	}
	defer rows.Close()

	var videos []VideoItem
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			log.Println("Error scanning video:", err)
			continue
		}
		videos = append(videos, VideoItem{
			FileName: filename,
			IsVideo:  contentType(filename) != "image/gif",
		})
	}
	return videos
}
