package common

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
