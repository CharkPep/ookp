package model

import "gorm.io/gorm"

type User struct {
	gorm.Model

	UID string `json:"uid"`
}
