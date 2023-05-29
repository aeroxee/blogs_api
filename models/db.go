package models

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	db = GetDB()
	db.AutoMigrate(&User{}, &Tag{}, &Article{})
}

// GetDB configuration for connection to database.
func GetDB() *gorm.DB {
	d, err := gorm.Open(sqlite.Open("blogs.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return d
}
