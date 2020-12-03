package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func dbClient(dns string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
