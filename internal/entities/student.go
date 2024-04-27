package entities

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	TelegramID int     `gorm:"not null"`
	Group      *string `gorm:"type:varchar(255)"`
	Subgroup   *string `gorm:"type:varchar(255)"`
}
