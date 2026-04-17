package main

import (
	"golang-trainning-database/config"
	"golang-trainning-database/infrastructure/datastore"
	"log"
)

func main() {
	config.ReadConfig()

	db := datastore.NewDB()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln(err)
	}
	defer sqlDB.Close()
}
