package models

import "gorm.io/gorm"

type Recommendation struct {
	gorm.Model
	From_user string `gorm:"not null"`
	To_user   string `gorm:"not null"`
	Url       string `gorm:"not null"`
	Name      string `gorm:"not null"`
	ImgUrl    string `gorm:"not null"`
}
