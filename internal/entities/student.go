package entities

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	TelegramID int
	Group      string
	Subgroup   *string
}
