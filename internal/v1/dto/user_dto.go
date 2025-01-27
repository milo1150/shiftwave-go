package dto

import (
	"shiftwave-go/internal/enum"
	"shiftwave-go/internal/model"
)

type UserModelDto struct {
	Id           int       `json:"id"`
	Username     string    `json:"username"`
	ActiveStatus bool      `json:"active_status"`
	Role         enum.Role `json:"role"`
	Branch       []int     `json:"user_branches"` // TODO:
}

func TransformUserModel(user model.User) UserModelDto {
	dto := UserModelDto{
		Id:           int(user.ID),
		Username:     user.Username,
		ActiveStatus: user.ActiveStatus,
		Role:         user.Role,
	}
	return dto
}
