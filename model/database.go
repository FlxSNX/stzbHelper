package model

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var Conn *gorm.DB

func InitDB(databaseName string) {
	db, err := gorm.Open(sqlite.Open(databaseName+".db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&TeamUser{}, &Task{}, &Report{})
	if err != nil {
		return
	}

	Conn = db
}
