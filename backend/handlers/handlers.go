// handlers/handlers.go
package handlers

import (
	"database/sql"
)

var DB *sql.DB // Declare a global variable to store the db connection

func InitializeDB(db *sql.DB) {
	DB = db
}
