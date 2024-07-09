package user

import (
	"context"

	"github.com/AndreiMartynenko/auth/internal/converter"
	desc "github.com/AndreiMartynenko/auth/pkg/user_v1"
)

// Get is used to obtain user info.
func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	user, err := i.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &desc.GetResponse{
		User: converter.ToUserFromService(user),
	}, nil
}
