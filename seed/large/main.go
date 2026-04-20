package main

import (
	"fmt"
	"golang-trainning-database/config"
	"golang-trainning-database/infrastructure/datastore"
	"log"
	"math/rand"
	"time"
)

const batchSize = 1000

var targetDistrictIDs = []uint{1, 2, 4, 5}
var bloodTypes = []string{"A", "B", "O", "AB"}
var birthplaces = []string{
	"北海道", "青森県", "岩手県", "宮城県", "秋田県",
	"山形県", "福島県", "茨城県", "栃木県", "群馬県",
	"埼玉県", "千葉県", "東京都", "神奈川県", "新潟県",
	"富山県", "石川県", "福井県", "山梨県", "長野県",
	"岐阜県", "静岡県", "愛知県", "三重県", "滋賀県",
	"京都府", "大阪府", "兵庫県", "奈良県", "和歌山県",
	"鳥取県", "島根県", "岡山県", "広島県", "山口県",
	"徳島県", "香川県", "愛媛県", "高知県", "福岡県",
	"佐賀県", "長崎県", "熊本県", "大分県", "宮崎県",
	"鹿児島県", "沖縄県",
}

type Company struct {
	ID uint `gorm:"primaryKey"`
}

func (Company) TableName() string { return "companies" }

type Store struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	CompanyID      uint      `gorm:"column:company_id"`
	DistrictID     uint      `gorm:"column:district_id"`
	ContractPlanID uint      `gorm:"column:contract_plan_id"`
	BusinessTypeID uint      `gorm:"column:business_type_id"`
	Name           string    `gorm:"column:name"`
	IsActive       bool      `gorm:"column:is_active"`
	OpenStatus     string    `gorm:"column:open_status"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (Store) TableName() string { return "stores" }

type Woman struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	CompanyID  uint      `gorm:"column:company_id"`
	Name       string    `gorm:"column:name"`
	Age        *int      `gorm:"column:age"`
	Birthplace string    `gorm:"column:birthplace"`
	BloodType  string    `gorm:"column:blood_type"`
	IsActive   bool      `gorm:"column:is_active"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (Woman) TableName() string { return "women" }

type WomanStoreAssignment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	WomanID   uint      `gorm:"column:woman_id"`
	StoreID   uint      `gorm:"column:store_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (WomanStoreAssignment) TableName() string { return "woman_store_assignments" }

func main() {
	config.ReadConfig()

	db := datastore.NewDB()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln(err)
	}
	defer sqlDB.Close()

	now := time.Now()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var companyIDs []uint
	if err := db.Model(&Company{}).Pluck("id", &companyIDs).Error; err != nil {
		log.Fatalln("fetch company ids:", err)
	}
	log.Printf("found %d companies\n", len(companyIDs))

	// stores
	log.Println("inserting stores...")
	stores := make([]Store, 1000)
	for i := range stores {
		stores[i] = Store{
			CompanyID:      companyIDs[rng.Intn(len(companyIDs))],
			DistrictID:     targetDistrictIDs[rng.Intn(len(targetDistrictIDs))],
			ContractPlanID: uint(rng.Intn(2) + 1), // 1 or 2
			BusinessTypeID: uint(rng.Intn(5) + 1), // 1-5
			Name:           fmt.Sprintf("テスト店舗%04d", i+1),
			IsActive:       true,
			OpenStatus:     "public",
			CreatedAt:      now,
			UpdatedAt:      now,
		}
	}
	if err := db.CreateInBatches(&stores, batchSize).Error; err != nil {
		log.Fatalln("stores:", err)
	}
	log.Printf("inserted %d stores\n", len(stores))

	storeIDs := make([]uint, len(stores))
	for i, s := range stores {
		storeIDs[i] = s.ID
	}

	// women
	log.Println("inserting women...")
	women := make([]Woman, 100000)
	for i := range women {
		age := rng.Intn(21) + 20
		women[i] = Woman{
			CompanyID:  companyIDs[rng.Intn(len(companyIDs))],
			Name:       fmt.Sprintf("テスト女性%06d", i+1),
			Age:        &age,
			Birthplace: birthplaces[rng.Intn(len(birthplaces))],
			BloodType:  bloodTypes[rng.Intn(len(bloodTypes))],
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
	}
	if err := db.CreateInBatches(&women, batchSize).Error; err != nil {
		log.Fatalln("women:", err)
	}
	log.Printf("inserted %d women\n", len(women))

	// assignments
	log.Println("inserting woman_store_assignments...")
	assignments := make([]WomanStoreAssignment, len(women))
	for i, w := range women {
		assignments[i] = WomanStoreAssignment{
			WomanID:   w.ID,
			StoreID:   storeIDs[rng.Intn(len(storeIDs))],
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	if err := db.CreateInBatches(&assignments, batchSize).Error; err != nil {
		log.Fatalln("assignments:", err)
	}
	log.Printf("inserted %d assignments\n", len(assignments))

	log.Println("seed completed.")
}
