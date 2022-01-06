package common

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
