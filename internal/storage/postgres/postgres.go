package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/lib/pq"

	"simple_RESTapi/internal/config"
	"simple_RESTapi/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(log *slog.Logger) (*Storage, error) {

	const op = "storage.sqlite.NewStorage" // Имя текущей функции для логов и ошибок

	// Подключение к базе данных
	dbConfig := config.MustLoad().DB

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Adres, strconv.Itoa(dbConfig.Port), dbConfig.DbUser, dbConfig.DbPassword, dbConfig.DbName)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
		panic("Can't connect to postgres")
	}

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully connected to Postgres!")

	// Создаем таблицу, если ее еще нет
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
			id SERIAL  PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
		`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//добавляем индексы в таблицу по alias
	stmt, err = db.Prepare(`
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
		`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.postgresql.SaveURL"

	// Подготавливаем запрос с использованием RETURNING для получения последнего вставленного id
	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	defer stmt.Close()

	// Выполняем запрос
	var lastInsertedId int64
	err = stmt.QueryRow(urlToSave, alias).Scan(&lastInsertedId)
	if err != nil {
		// Проверяем на наличие уже существующей записи
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	// Возвращаем ID
	return lastInsertedId, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgresql.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = $1")
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

func (s *Storage) URLDelete(alias string) error {
	const op = "storage.postgresql.DeleteURL"

	var exists bool
	stmt, err := s.db.Prepare("SELECT EXISTS(SELECT 1 FROM url WHERE alias = $1)")

	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	err = stmt.QueryRow(alias).Scan(&exists)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("%s: query row: %w", op, err)
	}

	stmt, err = s.db.Prepare("DELETE FROM url WHERE alias = $1")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	defer stmt.Close()

	stmt.Exec()

	if err != nil {
		// Проверяем на наличие уже существующей записи
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}

		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	if !exists {
		return storage.ErrURLNotFound
	}

	defer stmt.Close()

	_, err = stmt.Exec(alias)

	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrURLNotFound
	}

	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return nil
}
