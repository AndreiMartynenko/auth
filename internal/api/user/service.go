package user

import (
	"github.com/AndreiMartynenko/auth/internal/service"
	desc "github.com/AndreiMartynenko/auth/pkg/user_v1"
)

// Implementation structure describes API layer.
type Implementation struct {
	desc.UnimplementedUserV1Server
	userService service.UserService
}

// NewImplementation creates new object of API layer.
func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
