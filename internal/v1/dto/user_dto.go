package dto

import (
	"shiftwave-go/internal/enum"
	"shiftwave-go/internal/model"

	"github.com/google/uuid"
)

type UserModelDto struct {
	Uuid         uuid.UUID   `json:"user_uuid"`
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
		Username:     user.Username,
		ActiveStatus: user.ActiveStatus,
		Role:         user.Role,
		Branch:       TransformBranches(user.Branches),
		Uuid:         user.Uuid,
	}
	return dto
}
