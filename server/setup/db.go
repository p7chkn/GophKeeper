package setup

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
	"path/filepath"
	"runtime"
)

func MustSetupDatabase(db *sql.DB) error {
	log.Println("Enter a migrations start")
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	migrationsPath := basePath + "/migrations"
	err := goose.Up(db, migrationsPath)
	if err != nil {
		return err
	}
	return nil
}
