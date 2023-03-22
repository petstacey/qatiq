package service

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Service struct {
	Name string
	DB   *sql.DB
}
