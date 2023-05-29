package models

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	db = GetDB()
	db.AutoMigrate(&User{}, &Tag{}, &Article{})
}

// GetDB configuration for connection to database.
func GetDB() *gorm.DB {
	dsn := "host=localhost port=5432 user=fajhri password=root dbname=blogs sslmode=disable TimeZone=Asia/Jakarta"
	d, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return d
}
