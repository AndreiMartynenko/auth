package converter

import (
	"github.com/AndreiMartynenko/auth/internal/model"
	desc "github.com/AndreiMartynenko/auth/pkg/auth_v1"
)

// ToUserLoginFromDesc converts structure of API layer to service layer model.
func ToUserLoginFromDesc(creds *desc.Creds) *model.UserCreds {
	return &model.UserCreds{
		Username: creds.Username,
		Password: creds.Password,
	}
}
