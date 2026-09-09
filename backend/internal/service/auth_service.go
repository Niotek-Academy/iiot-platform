package service

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Niotek-Academy/iiot-platform/backend/internal/apperr"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/db/generated"
	"github.com/Niotek-Academy/iiot-platform/backend/internal/jwtutil"
)

type AuthService struct {
	store     *db.Store
	jwtSecret []byte
	jwtExpiry time.Duration
}

func NewAuthService(store *db.Store, jwtSecret string, jwtExpiryHours int) *AuthService {
	return &AuthService{
		store:     store,
		jwtSecret: []byte(jwtSecret),
		jwtExpiry: time.Duration(jwtExpiryHours) * time.Hour,
	}
}

// Register creates a new user.
//   - If no users exist yet, this is the bootstrap case: registration is
//     allowed with no actor, and the new user is always ADMIN regardless
//     of what role was requested.
//   - Otherwise, actorRole must be "ADMIN" (an empty actorRole means no/
//     invalid token was presented at all).
func (s *AuthService) Register(ctx context.Context, actorRole, username, password, requestedRole string) (generated.User, error) {
	count, err := s.store.CountUsers(ctx)
	if err != nil {
		return generated.User{}, apperr.NewInternal("could not check existing users")
	}

	role := requestedRole
	if role == "" {
		role = "OPERATOR"
	}

	if count == 0 {
		role = "ADMIN" // first user in the system is always the admin
	} else if actorRole != "ADMIN" {
		return generated.User{}, apperr.NewForbidden("only an admin can register new users")
	}

	if _, err := s.store.GetUserByUsername(ctx, username); err == nil {
		return generated.User{}, apperr.NewConflict("username already taken")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return generated.User{}, apperr.NewInternal("could not hash password")
	}

	user, err := s.store.CreateUser(ctx, generated.CreateUserParams{
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
	})
	if err != nil {
		return generated.User{}, apperr.NewInternal("could not create user")
	}
	return user, nil
}

// Login verifies credentials and issues a JWT. Both "user not found" and
// "wrong password" return the exact same error message, on purpose — never
// reveal which part was wrong to an attacker.
func (s *AuthService) Login(ctx context.Context, username, password string) (string, generated.User, error) {
	user, err := s.store.GetUserByUsername(ctx, username)
	if err != nil {
		return "", generated.User{}, apperr.NewUnauthorized("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", generated.User{}, apperr.NewUnauthorized("invalid username or password")
	}

	token, err := jwtutil.Generate(s.jwtSecret, s.jwtExpiry, user.UserID, user.Username, user.Role)
	if err != nil {
		return "", generated.User{}, apperr.NewInternal("could not generate token")
	}

	return token, user, nil
}
