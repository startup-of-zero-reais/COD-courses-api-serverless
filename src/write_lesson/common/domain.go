package common

import (
	"github.com/google/uuid"
	"time"
)

type (
	Lesson struct {
		SK            string            `json:"lesson_id"`
		PK            string            `json:"section_id"`
		Title         string            `json:"title"`
		Thumb         string            `json:"thumb"`
		VideoSource   string            `json:"video_source"`
		DurationTotal uint              `json:"duration_total"`
		ParentCourse  string            `json:"parent_course"`
		ParentModule  string            `json:"parent_module"`
		Artifacts     map[string]string `json:"artifacts"`
		CreatedAt     int64             `json:"created_at"`
		UpdatedAt     int64             `json:"updated_at"`
	}
)

func (l *Lesson) BeforeCreate() {
	l.SK = uuid.NewString()
	l.CreatedAt = time.Now().UnixMilli()
	l.UpdatedAt = time.Now().UnixMilli()
}

func (l *Lesson) BeforeUpdate() {
	l.UpdatedAt = time.Now().UnixMilli()
}
