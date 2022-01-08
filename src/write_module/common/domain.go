package common

import (
	"github.com/google/uuid"
	"time"
)

type (
	SectionsOrder map[string]string
	Module        struct {
		PK            string        `json:"course_id"`
		SK            string        `json:"module_id"`
		Title         string        `json:"title"`
		SectionsOrder SectionsOrder `json:"sections_order"`
		UnlockAfter   uint          `json:"unlock_after"`
		CreatedAt     int64         `json:"created_at"`
		UpdatedAt     int64         `json:"updated_at"`
	}
)

func (m *Module) BeforeCreate() {
	m.SK = uuid.NewString()
	m.CreatedAt = time.Now().UnixMilli()
	m.UpdatedAt = time.Now().UnixMilli()
}

func (m *Module) BeforeUpdate() {
	m.UpdatedAt = time.Now().UnixMilli()
}
