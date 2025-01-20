package types

type LoginPayload struct {
	Username string `json:"u" validate:"required"`
	Password string `json:"pwd" validate:"required"`
}

type CreateUserPayload struct {
	Username string `json:"u" validate:"required" example:"johndoe"`
	Password string `json:"pwd" validate:"required" example:"johndoe"`
}
