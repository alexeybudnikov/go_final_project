package database

import (
	_ "github.com/mattn/go-sqlite3"

	"database/sql"
	"log"
	"os"
	"path/filepath"
)

func InitDatabase() (*sql.DB, error) {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	envFile := os.Getenv("TODO_DBFILE")
	if len(envFile) > 0 {
		dbFile = envFile
	}

	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	if install {
		query := ` 
				CREATE TABLE IF NOT EXISTS scheduler (
    				id INTEGER PRIMARY KEY AUTOINCREMENT,
    				date TEXT(8) NOT NULL,
    				title TEXT NOT NULL,
    				comment TEXT,
   				 repeat TEXT(128)
				);

				CREATE INDEX IF NOT EXISTS idx_scheduler_date on scheduler(date);`
		_, err = db.Exec(query)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
