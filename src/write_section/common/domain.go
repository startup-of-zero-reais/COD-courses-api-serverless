package common

import (
	"github.com/google/uuid"
	"time"
)

type (
	Section struct {
		PK        string `json:"module_id"`
		SK        string `json:"section_id"`
		Title     string `json:"title"`
		CreatedAt int64  `json:"created_at"`
		UpdatedAt int64  `json:"updated_at"`
	}
)

func (s *Section) BeforeCreate() {
	s.SK = uuid.NewString()
	s.CreatedAt = time.Now().UnixMilli()
	s.UpdatedAt = time.Now().UnixMilli()
}

func (s *Section) BeforeUpdate() {
	s.UpdatedAt = time.Now().UnixMilli()
}
