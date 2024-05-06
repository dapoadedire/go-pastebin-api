package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"pastebin/config"
)

var DB *sql.DB

func InitDB() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", config.DB_HOST, config.DB_PORT, config.DB_USER, config.DB_NAME, config.DB_PASSWORD)
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	if err = DB.Ping(); err != nil {
		panic(err)
	}
	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS pastes (id TEXT PRIMARY KEY, content TEXT, created_at TIMESTAMP)")
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to the database!")
}
