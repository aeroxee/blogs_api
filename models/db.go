package models

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func init() {
	db = GetDB()
	db.AutoMigrate(&User{}, &Tag{}, &Article{})
}

// GetDB configuration for connection to database.
func GetDB() *gorm.DB {
	d, err := gorm.Open(sqlite.Open("blogs.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}

	return d
}
