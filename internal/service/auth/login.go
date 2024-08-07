package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/AndreiMartynenko/auth/internal/model"
)

func (s *serv) Login(ctx context.Context, creds *model.UserCreds) (string, error) {
	// Get role and hashed password by username from storage
	authInfo, err := s.userRepository.GetAuthInfo(ctx, creds.Username)
	if err != nil {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(authInfo.Password), []byte(creds.Password))
	if err != nil {
		return "", errors.New("wrong password")
	}

	// Get secret key from storage for refresh token HMAC
	refreshTokenSecretKey, err := s.keyRepository.GetKey(ctx, refreshTokenSecretKeyName)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	refreshToken, err := s.tokenOperations.Generate(model.User{
		Name: authInfo.Username,
		Role: authInfo.Role,
	},
		[]byte(refreshTokenSecretKey),
		refreshTokenExpiration,
	)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return refreshToken, nil
}
