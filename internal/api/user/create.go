package user

import (
	"context"

	"github.com/AndreiMartynenko/auth/internal/converter"
	desc "github.com/AndreiMartynenko/auth/pkg/user_v1"
)

// Create is used for creating new user.
func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.userService.Create(ctx, converter.ToUserCreateFromDesc(req.GetUser()))
	if err != nil {
		return nil, err
	}

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
