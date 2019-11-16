package model

import (
	"time"
)

type Todo struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	IsDone    bool      `db:"is_done" json:"is_done"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
