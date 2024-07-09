package user

import (
	"github.com/AndreiMartynenko/auth/internal/repository"
	"github.com/AndreiMartynenko/auth/internal/service"
	"github.com/AndreiMartynenko/common/pkg/db"
)

type serv struct {
	userRepository repository.UserRepository
	logRepository  repository.LogRepository
	txManager      db.TxManager
}

// NewService creates new object of service layer.
func NewService(userRepository repository.UserRepository, logRepository repository.LogRepository, txManager db.TxManager) service.UserService {
	return &serv{
		userRepository: userRepository,
		logRepository:  logRepository,
		txManager:      txManager,
	}
}
