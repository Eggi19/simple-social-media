package entities

import (
	"database/sql"
	"time"
)

type Tweet struct {
	Id        int64
	Tweet     string
	UserId    int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}
