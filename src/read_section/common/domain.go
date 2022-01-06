package common

type (
	Section struct {
		PK        string `json:"module_id"`
		SK        string `json:"section_id"`
		Title     string `json:"title"`
		CreatedAt int64  `json:"created_at"`
		UpdatedAt int64  `json:"updated_at"`
	}
)
