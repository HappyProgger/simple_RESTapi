// internal/storage/storage.go

package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"simple_RESTapi/internal/config"
)

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
)

func main() {

	dbCfg := config.MustLoad().DB
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbCfg.Adres, dbCfg.Port, dbCfg.DbUser, dbCfg.DbPassword, dbCfg.DbName)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected!")

	// Пример запроса
	rows, err := db.Query("SELECT * FROM your_table_name LIMIT 10;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
