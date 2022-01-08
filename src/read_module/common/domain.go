package common

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
