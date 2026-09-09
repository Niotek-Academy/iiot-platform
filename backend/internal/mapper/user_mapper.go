package mapper

import (
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db/generated"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/dto"
)

func ToUserResponse(u generated.User) dto.UserResponse {
	return dto.UserResponse{
		UserID:    u.UserID,
		Username:  u.Username,
		Role:      u.Role,
		CreatedAt: u.CreatedAt.Time, // pgtype.Timestamptz -> time.Time
	}
}