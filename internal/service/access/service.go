package access

import (
	"time"

	"github.com/AndreiMartynenko/auth/internal/repository"
	"github.com/AndreiMartynenko/auth/internal/service"
	"github.com/AndreiMartynenko/auth/internal/tokens"
)

const (
	accessTokenExpiration    = 30 * time.Minute
	accessTokenSecretKeyName = "access"
)

type serv struct {
	accessRepository repository.AccessRepository
	keyRepository    repository.KeyRepository
	tokenOperations  tokens.TokenOperations
}

// NewService creates new object of service layer.
func NewService(accessRepository repository.AccessRepository, keyRepository repository.KeyRepository, tokenOperations tokens.TokenOperations) service.AccessService {
	return &serv{
		accessRepository: accessRepository,
		keyRepository:    keyRepository,
		tokenOperations:  tokenOperations,
	}
}
