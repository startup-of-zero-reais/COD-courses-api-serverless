package common

type (
	Artifact struct {
		ID          string `json:"artifact_id"`
		ParentID    string `json:"lesson_id"`
		MediaSource string `json:"media_source"`
	}

	Lesson struct {
		ID            string     `json:"lesson_id"`
		ParentID      string     `json:"section_id"`
		Label         string     `json:"label"`
		MediaSource   string     `json:"video_source"`
		DurationTotal uint       `json:"duration_total"`
		Artifacts     []Artifact `json:"artifacts"`
	}
)
