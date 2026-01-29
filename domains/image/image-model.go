package image

import (
	"vixel/domains/user"

	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	URL             string `gorm:"not null"`
	AltText         string
	UserID          uint      `gorm:"not null"`
	User            user.User `gorm:"foreignKey:UserID"`
	Transformations []string  `gorm:"type:json"`
}
