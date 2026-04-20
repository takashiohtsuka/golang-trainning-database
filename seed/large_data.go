package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

// --- models ---

type Company struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	BusinessTypeID uint      `gorm:"column:business_type_id"`
	Name           string    `gorm:"column:name"`
	IsActive       bool      `gorm:"column:is_active"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (Company) TableName() string { return "companies" }

type Store struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	CompanyID  uint      `gorm:"column:company_id"`
	DistrictID uint      `gorm:"column:district_id"`
	Name       string    `gorm:"column:name"`
	IsActive   bool      `gorm:"column:is_active"`
	OpenStatus string    `gorm:"column:open_status"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
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

// --- constants ---

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

// --- seed functions ---

func seedCompanies(db *gorm.DB, rng *rand.Rand, now time.Time) []uint {
	const num = 20
	companies := make([]Company, num)
	for i := range companies {
		companies[i] = Company{
			BusinessTypeID: uint(rng.Intn(5) + 1),
			Name:           fmt.Sprintf("テスト会社%03d", i+1),
			IsActive:       true,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
	}
	if err := db.CreateInBatches(&companies, batchSize).Error; err != nil {
		log.Fatalln("companies:", err)
	}
	ids := make([]uint, len(companies))
	for i, c := range companies {
		ids[i] = c.ID
	}
	log.Printf("inserted %d companies\n", len(companies))
	return ids
}

func seedStores(db *gorm.DB, rng *rand.Rand, now time.Time, companyIDs []uint) []uint {
	const num = 1000
	stores := make([]Store, num)
	for i := range stores {
		stores[i] = Store{
			CompanyID:  companyIDs[rng.Intn(len(companyIDs))],
			DistrictID: targetDistrictIDs[rng.Intn(len(targetDistrictIDs))],
			Name:       fmt.Sprintf("テスト店舗%04d", i+1),
			IsActive:   true,
			OpenStatus: "public",
			CreatedAt:  now,
			UpdatedAt:  now,
		}
	}
	if err := db.CreateInBatches(&stores, batchSize).Error; err != nil {
		log.Fatalln("stores:", err)
	}
	ids := make([]uint, len(stores))
	for i, s := range stores {
		ids[i] = s.ID
	}
	log.Printf("inserted %d stores\n", len(stores))
	return ids
}

func seedWomen(db *gorm.DB, rng *rand.Rand, now time.Time, companyIDs []uint) []uint {
	const num = 100000
	women := make([]Woman, num)
	for i := range women {
		age := rng.Intn(21) + 20 // 20-40
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
	ids := make([]uint, len(women))
	for i, w := range women {
		ids[i] = w.ID
	}
	log.Printf("inserted %d women\n", len(women))
	return ids
}

func seedWomanStoreAssignments(db *gorm.DB, rng *rand.Rand, now time.Time, womanIDs []uint, storeIDs []uint) {
	assignments := make([]WomanStoreAssignment, len(womanIDs))
	for i, wid := range womanIDs {
		assignments[i] = WomanStoreAssignment{
			WomanID:   wid,
			StoreID:   storeIDs[rng.Intn(len(storeIDs))],
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	if err := db.CreateInBatches(&assignments, batchSize).Error; err != nil {
		log.Fatalln("woman_store_assignments:", err)
	}
	log.Printf("inserted %d woman_store_assignments\n", len(assignments))
}
