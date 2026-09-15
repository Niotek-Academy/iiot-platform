package dto

import "time"

type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=ADMIN OPERATOR"`
}

type UserResponse struct {
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}