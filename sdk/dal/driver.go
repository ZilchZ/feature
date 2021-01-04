package dal

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	MaxConn  int
}

func connDB(cfg DatabaseConfig) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("dbname=%s host=%s user=%s password=%s port=%d sslmode=disable",
		cfg.DBName, cfg.Host, cfg.User, cfg.Password, cfg.Port)
	db, err := sqlx.Open(cfg.Driver, connStr)
	if err != nil || db == nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return db, err
	}
	return db, nil
}
