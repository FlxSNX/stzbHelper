package model

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var Conn *gorm.DB

func InitDB() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&TeamUser{}, &Task{})
	if err != nil {
		return
	}

	Conn = db
}
