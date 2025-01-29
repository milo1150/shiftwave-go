package dto

import (
	"shiftwave-go/internal/enum"
	"shiftwave-go/internal/model"
)

type UserModelDto struct {
	Id           int         `json:"id"`
	Username     string      `json:"username"`
	ActiveStatus bool        `json:"active_status"`
	Role         enum.Role   `json:"role"`
	Branch       []BranchDto `json:"branches"`
}

func TransformUserModels(users []model.User) []UserModelDto {
	if len(users) == 0 {
		return []UserModelDto{}
	}

	userDtos := make([]UserModelDto, len(users))
	for i, user := range users {
		userDtos[i] = TransformUserModel(user)
	}

	return userDtos
}

func TransformUserModel(user model.User) UserModelDto {
	dto := UserModelDto{
		Id:           int(user.ID),
		Username:     user.Username,
		ActiveStatus: user.ActiveStatus,
		Role:         user.Role,
		Branch:       TransformBranches(user.Branches),
	}
	return dto
}
