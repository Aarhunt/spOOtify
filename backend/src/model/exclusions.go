package model

import "gorm.io/gorm"

type AutoExclusion struct {
	gorm.Model
	Substring string `json:"substring"`
}
