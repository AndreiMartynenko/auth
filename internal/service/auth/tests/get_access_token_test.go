package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/AndreiMartynenko/auth/internal/model"
	"github.com/AndreiMartynenko/auth/internal/repository"
	repositoryMocks "github.com/AndreiMartynenko/auth/internal/repository/mocks"
	authService "github.com/AndreiMartynenko/auth/internal/service/auth"
	"github.com/AndreiMartynenko/auth/internal/tokens"
	tokenMocks "github.com/AndreiMartynenko/auth/internal/tokens/mocks"
)

func TestGetAccessToken(t *testing.T) {
	t.Parallel()

	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type keyRepositoryMockFunc func(mc *minimock.Controller) repository.KeyRepository
	type tokenOperationsMockFunc func(mc *minimock.Controller) tokens.TokenOperations

	type args struct {
		ctx context.Context
		req string
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		accessKeyName   = "access"
		accessKey       = "access_key"
		accessKeyBytes  = []byte("access_key")
		refreshKeyName  = "refresh"
		refreshKey      = "refresh_key"
		refreshKeyBytes = []byte("refresh_key")

		refreshToken          = "refresh_token"
		accessToken           = "access_token"
		accessTokenExpiration = 30 * time.Minute

		username = "username"
		role     = "USER"

		claims = &model.UserClaims{
			Username: username,
			Role:     role,
		}

		user = model.User{
			Name: username,
			Role: role,
		}

		repositoryErr   = fmt.Errorf("failed to generate token")
		tokenInvalidErr = fmt.Errorf("invalid refresh token")

		req = refreshToken
		res = accessToken
	)

	tests := []struct {
		name                string
		args                args
		want                string
		err                 error
		userRepositoryMock  userRepositoryMockFunc
		keyRepositoryMock   keyRepositoryMockFunc
		tokenOperationsMock tokenOperationsMockFunc
	}{
		{
			name: "refresh key repository error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: "",
			err:  repositoryErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				return mock
			},
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				mock.GetKeyMock.Expect(minimock.AnyContext, refreshKeyName).Return("", repositoryErr)
				return mock
			},
			tokenOperationsMock: func(mc *minimock.Controller) tokens.TokenOperations {
				mock := tokenMocks.NewTokenOperationsMock(mc)
				return mock
			},
		},
		{
			name: "access key repository error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: "",
			err:  repositoryErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				return mock
			},
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				mock.GetKeyMock.When(minimock.AnyContext, refreshKeyName).Then(refreshKey, nil)
				mock.GetKeyMock.When(minimock.AnyContext, accessKeyName).Then("", repositoryErr)
				return mock
			},
			tokenOperationsMock: func(mc *minimock.Controller) tokens.TokenOperations {
				mock := tokenMocks.NewTokenOperationsMock(mc)
				return mock
			},
		},
		{
			name: "token verify error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: "",
			err:  tokenInvalidErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				return mock
			},
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				mock.GetKeyMock.When(minimock.AnyContext, refreshKeyName).Then(refreshKey, nil)
				mock.GetKeyMock.When(minimock.AnyContext, accessKeyName).Then(accessKey, nil)
				return mock
			},
			tokenOperationsMock: func(mc *minimock.Controller) tokens.TokenOperations {
				mock := tokenMocks.NewTokenOperationsMock(mc)
				mock.VerifyMock.Expect(refreshToken, refreshKeyBytes).Return(nil, tokenInvalidErr)
				return mock
			},
		},
		{
			name: "token generate error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: "",
			err:  repositoryErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				return mock
			},
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				mock.GetKeyMock.When(minimock.AnyContext, refreshKeyName).Then(refreshKey, nil)
				mock.GetKeyMock.When(minimock.AnyContext, accessKeyName).Then(accessKey, nil)
				return mock
			},
			tokenOperationsMock: func(mc *minimock.Controller) tokens.TokenOperations {
				mock := tokenMocks.NewTokenOperationsMock(mc)
				mock.VerifyMock.Expect(refreshToken, refreshKeyBytes).Return(claims, nil)
				mock.GenerateMock.Expect(user, accessKeyBytes, accessTokenExpiration).Return("", repositoryErr)
				return mock
			},
		},
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				return mock
			},
			keyRepositoryMock: func(mc *minimock.Controller) repository.KeyRepository {
				mock := repositoryMocks.NewKeyRepositoryMock(mc)
				mock.GetKeyMock.When(minimock.AnyContext, refreshKeyName).Then(refreshKey, nil)
				mock.GetKeyMock.When(minimock.AnyContext, accessKeyName).Then(accessKey, nil)
				return mock
			},
			tokenOperationsMock: func(mc *minimock.Controller) tokens.TokenOperations {
				mock := tokenMocks.NewTokenOperationsMock(mc)
				mock.VerifyMock.Expect(refreshToken, refreshKeyBytes).Return(claims, nil)
				mock.GenerateMock.Expect(user, accessKeyBytes, accessTokenExpiration).Return(accessToken, nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepositoryMock := tt.userRepositoryMock(mc)
			keyRepositoryMock := tt.keyRepositoryMock(mc)
			tokenOperationsMock := tt.tokenOperationsMock(mc)
			srv := authService.NewService(userRepositoryMock, keyRepositoryMock, tokenOperationsMock)

			res, err := srv.GetAccessToken(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
