package handlers

import "sync"

type Photo struct {
	Filename string
	Likes    int
}

var (
	photoLikes = make(map[string]int)
	likesMutex = &sync.Mutex{}
)
