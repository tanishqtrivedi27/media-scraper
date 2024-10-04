package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     int
}

func NewPostgreSQLStorage(cfg PostgresConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

type DBStorage struct {
	db *sql.DB
}

func NewDBStorage(db *sql.DB) (*DBStorage, error) {
	return &DBStorage{db: db}, nil
}

func (dbStore *DBStorage) InsertTuple(realUrl, storedUrl string) error {
	_, err := dbStore.db.Exec("INSERT INTO url_metadata(realUrl, storedUrl) VALUES ($1, $2)", realUrl, storedUrl)
	return err
}
