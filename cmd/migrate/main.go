package main

import (
	"database/sql"
	"log"
	"path/filepath"
	"runtime"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"

	"golang-trainning-database/config"
)

func main() {
	config.ReadConfig()

	migrationsDir := filepath.Join(rootDir(), "db", "migrations")

	if err := goose.SetDialect("mysql"); err != nil {
		log.Fatalln(err)
	}

	migrate(dsn(config.C.Database.DBName), migrationsDir)
	migrate(dsn(config.C.Database.DBName+"_test"), migrationsDir)
}

func migrate(dsn, migrationsDir string) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	if err := goose.Up(db, migrationsDir); err != nil {
		log.Fatalln(err)
	}
}

func dsn(dbName string) string {
	c := config.C.Database
	return c.User + ":" + c.Password + "@tcp(" + c.Addr + ")/" + dbName + "?parseTime=true"
}

func rootDir() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "..", "..")
}
