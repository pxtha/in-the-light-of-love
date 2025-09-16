package handlers

import (
	"sync"
	"time"

	"gorm.io/gorm"
)

type Photo struct {
	ID        uint            `gorm:"primaryKey"`
	CreatedAt time.Time       `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time       `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Filename  string
	Uploader  string
	Likes     int
	ModTime   time.Time
}

var (
	photoLikes = make(map[string]int)
	likesMutex = &sync.Mutex{}
)
