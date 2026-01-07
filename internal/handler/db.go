package handler

import "database/sql"

var db *sql.DB

func SetDB(d *sql.DB) {
	db = d
}
