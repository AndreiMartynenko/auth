package access

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	desc "github.com/AndreiMartynenko/auth/pkg/access_v1"
)

// Check performs user authorization.
func (i *Implementation) Check(ctx context.Context, req *desc.CheckRequest) (*empty.Empty, error) {
	err := i.accessService.Check(ctx, req.GetEndpoint())
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
