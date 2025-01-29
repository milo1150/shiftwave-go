package types

import "github.com/google/uuid"

type LoginPayload struct {
	Username string `json:"u" validate:"required"`
	Password string `json:"pwd" validate:"required"`
}

type CreateUserPayload struct {
	Username string      `json:"u" validate:"required" example:"johndoe_username"`
	Password string      `json:"pwd" validate:"required" example:"johndoe_pwd"`
	Role     string      `json:"role" validate:"required,userRole"`
	Branches []uuid.UUID `json:"branches" validate:"required,min=1"`
}

type UpdateUserPayload struct {
	Role         string      `json:"role" validate:"userRole"`
	Branches     []uuid.UUID `json:"branches" validate:"min=1"`
	ActiveStatus bool        `json:"active_status" validate:"boolean"`
}

type UpdateUserPasswordPayload struct {
	Password string `json:"pwd" validate:"omitempty,ascii"`
}
