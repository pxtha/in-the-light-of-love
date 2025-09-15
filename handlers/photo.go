package handlers

import (
	"sync"
	"time"
)

type Photo struct {
	Filename string
	Likes    int
	ModTime  time.Time
}

var (
	photoLikes = make(map[string]int)
	likesMutex = &sync.Mutex{}
)
