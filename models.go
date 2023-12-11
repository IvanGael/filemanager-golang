// models.go

package main

import (
	"time"
	// "gorm.io/gorm"
)

type File struct {
	// gorm.Model
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	FileName  string `gorm:"unique;not null"`
	FilePath  string `gorm:"not null"`
}
