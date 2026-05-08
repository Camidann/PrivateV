package main

import (
	// ...otros imports...
	"database/sql"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "./videos.db")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS videos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    filename TEXT UNIQUE,
    titulo TEXT,
    descripcion TEXT,
    contenido BLOB,
    usuario_id INTEGER,
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id)
)`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS usuarios (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    usuario TEXT UNIQUE,
    password TEXT
)`)
	if err != nil {
		panic(err)
	}

}
func getUsuarioID(usuario string) (int, error) {
	var id int
	err := db.QueryRow("SELECT id FROM usuarios WHERE usuario = ?", usuario).Scan(&id)
	return id, err
}
