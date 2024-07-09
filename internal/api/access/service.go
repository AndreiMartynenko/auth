package access

import (
	"github.com/AndreiMartynenko/auth/internal/service"
	desc "github.com/AndreiMartynenko/auth/pkg/access_v1"
)

// Implementation structure describes API layer.
type Implementation struct {
	desc.UnimplementedAccessV1Server
	accessService service.AccessService
}

// NewImplementation creates new object of API layer.
func NewImplementation(accessService service.AccessService) *Implementation {
	return &Implementation{
		accessService: accessService,
	}
}
