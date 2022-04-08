package data

import (
	"loadbalancer/data/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DbConnection *gorm.DB

func CreateDbConnection() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("loadbalancer.db"), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Node{})

	return db, nil
}
