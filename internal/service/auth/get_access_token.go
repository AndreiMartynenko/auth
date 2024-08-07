package auth

import (
	"context"
	"errors"

	"github.com/AndreiMartynenko/auth/internal/model"
)

func (s *serv) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	// Get secret key from storage for refresh token HMAC
	refreshTokenSecretKey, err := s.keyRepository.GetKey(ctx, refreshTokenSecretKeyName)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	// Get secret key from storage for access token HMAC
	accessTokenSecretKey, err := s.keyRepository.GetKey(ctx, accessTokenSecretKeyName)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	claims, err := s.tokenOperations.Verify(refreshToken, []byte(refreshTokenSecretKey))
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	accessToken, err := s.tokenOperations.Generate(model.User{
		Name: claims.Username,
		Role: claims.Role,
	},
		[]byte(accessTokenSecretKey),
		accessTokenExpiration,
	)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return accessToken, nil
}
