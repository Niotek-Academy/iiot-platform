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

func ToUserResponseList(list []generated.User) []dto.UserResponse {
	out := make([]dto.UserResponse, 0, len(list))
	for _, u := range list {
		out = append(out, ToUserResponse(u))
	}
	return out
}