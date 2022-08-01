package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New(path string) (db *gorm.DB, err error) {
	db, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return
	}

	db.AutoMigrate(&Alarm{})
	db.AutoMigrate(&Registration{})
	db.AutoMigrate(&FoodTag{})
	db.AutoMigrate(&FoodTagRelation{})
	db.AutoMigrate(&Action{})
	return
}
