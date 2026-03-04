package db

import (
	"database/sql"
	"fmt"
	"os"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func Init(dbFile string) error {
	schema := `CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
    		date CHAR(8) NOT NULL DEFAULT "",
			title VARCHAR(256) NOT NULL DEFAULT "",
			comment TEXT DEFAULT "",
			repeat VARCHAR(128) DEFAULT ""
			);
			
			CREATE INDEX scheduler_date ON scheduler (date);`

	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
    	install = true
	}

	dbF, err := sql.Open("sqlite", dbFile)
	if err != nil {
        return fmt.Errorf("Ошибка открытия БД: %w", err)
    }

	if install {
		_, err := dbF.Exec(schema)
		dbF.Close()
		if err != nil {
			return fmt.Errorf("Ошибка создания таблицы или индекса: %w", err)
		}
	}
	db = dbF
	return nil
}

func Close() error {
	if db != nil {
		fmt.Println("close connection to bd")
		return db.Close()
	}
	return nil
}