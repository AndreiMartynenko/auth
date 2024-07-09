package auth

import (
	"context"

	"github.com/AndreiMartynenko/auth/internal/converter"
	desc "github.com/AndreiMartynenko/auth/pkg/auth_v1"
)

// Login user and return refresh token.
func (i *Implementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	refreshToken, err := i.authService.Login(ctx, converter.ToUserLoginFromDesc(req.GetCreds()))
	if err != nil {
		return nil, err
	}

	return &desc.LoginResponse{
		RefreshToken: refreshToken,
	}, nil
}
