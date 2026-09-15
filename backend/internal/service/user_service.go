package service

import (
	"context"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/apperr"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db/generated"
)

type UserService struct {
	store *db.Store
}

func NewUserService(store *db.Store) *UserService {
	return &UserService{store: store}
}

func (s *UserService) List(ctx context.Context) ([]generated.User, error) {
	users, err := s.store.ListUsers(ctx)
	if err != nil {
		return nil, apperr.NewInternal("could not list users")
	}
	return users, nil
}


func (s *UserService) UpdateRole(ctx context.Context, actorUserID, targetUserID int64, role string) (generated.User, error) {
	if actorUserID == targetUserID {
		return generated.User{}, apperr.NewBadRequest("cannot change your own role")
	}

	if _, err := s.store.GetUserByID(ctx, targetUserID); err != nil {
		return generated.User{}, apperr.NewNotFound("user not found")
	}

	user, err := s.store.UpdateUserRole(ctx, generated.UpdateUserRoleParams{
		UserID: targetUserID,
		Role:   role,
	})
	if err != nil {
		return generated.User{}, apperr.NewInternal("could not update user role")
	}
	return user, nil
}