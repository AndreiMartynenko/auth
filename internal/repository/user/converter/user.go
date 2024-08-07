package converter

import (
	model "github.com/AndreiMartynenko/auth/internal/model"
	modelRepo "github.com/AndreiMartynenko/auth/internal/repository/user/model"
)

// ToUserFromRepo converts repository layer model to structure of service layer.
func ToUserFromRepo(user *modelRepo.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToAuthInfoFromRepo converts repository layer model to structure of service layer.
func ToAuthInfoFromRepo(authInfo *modelRepo.AuthInfo) *model.AuthInfo {
	return &model.AuthInfo{
		Username: authInfo.Username,
		Role:     authInfo.Role,
		Password: authInfo.Password,
	}
}
