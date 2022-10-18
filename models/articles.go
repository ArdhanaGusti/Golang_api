package models

import (
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Title  string
	Tag    string
	Slug   string `gorm:"unique_index"`
	Desc   string `sql:"type:text;"`
	UserID uint
	User   User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
