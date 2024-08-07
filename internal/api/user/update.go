package user

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/AndreiMartynenko/auth/internal/converter"
	desc "github.com/AndreiMartynenko/auth/pkg/user_v1"
)

// Update is used for updating user info.
func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {
	err := i.userService.Update(ctx, converter.ToUserUpdateFromDesc(req.GetUser()))
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
