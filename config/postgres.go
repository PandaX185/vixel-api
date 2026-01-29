package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(Config.DbUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
