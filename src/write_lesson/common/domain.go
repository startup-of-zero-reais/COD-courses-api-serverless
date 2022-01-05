package common

type (
	Lesson struct {
		SK            string            `json:"lesson_id"`
		PK            string            `json:"section_id"`
		Title         string            `json:"title"`
		VideoSource   string            `json:"video_source"`
		DurationTotal uint              `json:"duration_total"`
		ParentCourse  string            `json:"parent_course"`
		ParentModule  string            `json:"parent_module"`
		Artifacts     map[string]string `json:"artifacts"`
	}
)
