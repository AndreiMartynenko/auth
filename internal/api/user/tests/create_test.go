package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	userAPI "github.com/AndreiMartynenko/auth/internal/api/user"
	"github.com/AndreiMartynenko/auth/internal/model"
	"github.com/AndreiMartynenko/auth/internal/service"
	serviceMocks "github.com/AndreiMartynenko/auth/internal/service/mocks"
	desc "github.com/AndreiMartynenko/auth/pkg/user_v1"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id              = int64(1)
		name            = "name"
		email           = "email"
		password        = "password"
		passwordConfirm = "passwordConfirm"
		role            = desc.Role_USER
		roleName        = "USER"

		serviceErr = fmt.Errorf("service error")

		req = &desc.CreateRequest{
			User: &desc.UserCreate{
				Name:            name,
				Email:           email,
				Password:        password,
				PasswordConfirm: passwordConfirm,
				Role:            role,
			},
		}

		userCreate = &model.UserCreate{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: passwordConfirm,
			Role:            roleName,
		}

		res = &desc.CreateResponse{
			Id: id,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(minimock.AnyContext, userCreate).Return(id, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(minimock.AnyContext, userCreate).Return(0, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := userAPI.NewImplementation(userServiceMock)

			res, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
