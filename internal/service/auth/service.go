package auth

import (
	"time"

	"github.com/AndreiMartynenko/auth/internal/repository"
	"github.com/AndreiMartynenko/auth/internal/service"
	"github.com/AndreiMartynenko/auth/internal/tokens"
)

const (
	refreshTokenExpiration    = 360 * time.Minute
	accessTokenExpiration     = 30 * time.Minute
	refreshTokenSecretKeyName = "refresh"
	accessTokenSecretKeyName  = "access"
)

type serv struct {
	userRepository  repository.UserRepository
	keyRepository   repository.KeyRepository
	tokenOperations tokens.TokenOperations
}

// NewService creates new object of service layer.
func NewService(userRepository repository.UserRepository, keyRepository repository.KeyRepository, tokenOperations tokens.TokenOperations) service.AuthService {
	return &serv{
		userRepository:  userRepository,
		keyRepository:   keyRepository,
		tokenOperations: tokenOperations,
	}
}
