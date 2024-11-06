package entities

import (
	"database/sql"
	"time"
)

type User struct {
	Id           int64
	Name         string
	Email        string
	Password     string
	FcmToken     sql.NullString
	PhotoProfile sql.NullString
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    sql.NullTime
}

type UserFirestore struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	FollowerCount  int64  `json:"follower_count"`
	FollowingCount int64  `json:"followeing_count"`
}
