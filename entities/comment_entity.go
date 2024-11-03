package entities

import (
	"database/sql"
	"time"
)

type Comment struct {
	Id        int64
	Comment   string
	UserId    int64
	TweetId   int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}
