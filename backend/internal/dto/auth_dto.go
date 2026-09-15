package dto


type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=8"`
	Role string `json:"role" binding:"omitempty,oneof=ADMIN OPERATOR"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}


type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}