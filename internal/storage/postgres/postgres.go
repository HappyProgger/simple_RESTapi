package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/lib/pq"

	"simple_RESTapi/internal/config"
	"simple_RESTapi/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New() (*Storage, error) {

	const op = "storage.sqlite.NewStorage" // Имя текущей функции для логов и ошибок

	// Подключение к базе данных
	dbConfig := config.MustLoad().DB
	println(dbConfig.DbUser)
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Adres, strconv.Itoa(dbConfig.Port), dbConfig.DbUser, dbConfig.DbPassword, dbConfig.DbName)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected!")

	// Создаем таблицу, если ее еще нет
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
		`)

	defer stmt.Close()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	// Подготавливаем запрос
	stmt, err := s.db.Prepare("INSERT INTO url(url,alias) values(?,?)")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	defer stmt.Close()

	// Выполняем запрос
	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	// Получаем ID созданной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	// Возвращаем ID
	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	defer stmt.Close()

	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resURL, nil
}
