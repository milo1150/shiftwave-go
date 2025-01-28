package types

type LoginPayload struct {
	Username string `json:"u" validate:"required"`
	Password string `json:"pwd" validate:"required"`
}

type CreateUserPayload struct {
	Username string   `json:"u" validate:"required" example:"johndoe_username"`
	Password string   `json:"pwd" validate:"required" example:"johndoe_pwd"`
	Role     string   `json:"role" validate:"required,userRole"`
	Branches []string `json:"branches" validate:"required,branches"` // uuid
}
