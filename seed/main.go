package main

import (
	"flag"
	"golang-trainning-database/config"
	"golang-trainning-database/infrastructure/datastore"
	"log"
	"math/rand"
	"time"
)

const batchSize = 1000

func main() {
	mode := flag.String("mode", "init", "seed mode: init | add-women")
	flag.Parse()

	config.ReadConfig()

	db := datastore.NewDB()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln(err)
	}
	defer sqlDB.Close()

	now := time.Now()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	switch *mode {
	case "init":
		companyIDs := seedCompanies(db, rng, now)
		storeIDs := seedStores(db, rng, now, companyIDs)
		womanIDs := seedWomen(db, rng, now, companyIDs)
		seedWomanStoreAssignments(db, rng, now, womanIDs, storeIDs)

	case "add-women":
		var companyIDs []uint
		if err := db.Model(&Company{}).Pluck("id", &companyIDs).Error; err != nil {
			log.Fatalln("fetch company ids:", err)
		}
		var storeIDs []uint
		if err := db.Model(&Store{}).Where("district_id IN ?", targetDistrictIDs).Pluck("id", &storeIDs).Error; err != nil {
			log.Fatalln("fetch store ids:", err)
		}
		womanIDs := seedWomen(db, rng, now, companyIDs)
		seedWomanStoreAssignments(db, rng, now, womanIDs, storeIDs)

	case "add":
		var companyIDs []uint
		if err := db.Model(&Company{}).Pluck("id", &companyIDs).Error; err != nil {
			log.Fatalln("fetch company ids:", err)
		}
		storeIDs := seedStores(db, rng, now, companyIDs)
		womanIDs := seedWomen(db, rng, now, companyIDs)
		seedWomanStoreAssignments(db, rng, now, womanIDs, storeIDs)

	default:
		log.Fatalf("unknown mode: %s\n", *mode)
	}

	log.Println("seed completed.")
}
