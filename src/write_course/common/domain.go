package common

import (
	"github.com/google/uuid"
	"time"
)

type (
	Course struct {
		PK          string `json:"user_id"`
		SK          string `json:"course_id"`
		Thumb       string `json:"thumb"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Owner       string `json:"owner"`
		CartOpen    bool   `json:"cart_open"`
		CreatedAt   int64  `json:"created_at"`
		UpdatedAt   int64  `json:"updated_at"`
	}
)

func (c *Course) BeforeCreate() {
	c.SK = uuid.NewString()
	c.CreatedAt = time.Now().UnixMilli()
	c.UpdatedAt = time.Now().UnixMilli()
}

func (c *Course) BeforeUpdate() {
	c.UpdatedAt = time.Now().UnixMilli()
}
