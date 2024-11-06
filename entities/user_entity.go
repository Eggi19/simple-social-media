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
	Name           string
	Email          string
	FollowerCount  int64
	FollowingCount int64
}

type FollowerFirestore struct {
	FollowerId string
	FollowedAt time.Time
}

type FollowingFirestore struct {
	FollowingId string
	FollowedAt  time.Time
}
