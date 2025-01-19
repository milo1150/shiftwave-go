package types

type LoginPayload struct {
	User     string `json:"u" validate:"required"`
	Password string `json:"pwd" validate:"required"`
}
