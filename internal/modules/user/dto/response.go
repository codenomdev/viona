package dto

import "time"

type (
	ResponseUser struct {
		ID              int64      `json:"user_id"`
		Username        string     `json:"username"`
		FullName        string     `json:"full_name"`
		Email           string     `json:"email"`
		Password        string     `json:"-"`
		Avatar          string     `json:"avatar_url"`
		EmailVerifiedAt *time.Time `json:"verified_at,omitempty"`
		IsActive        int        `json:"is_active"`
		CreatedAt       time.Time  `json:"created_at"`
		UpdatedAt       time.Time  `json:"updated_at"`
	}
)
