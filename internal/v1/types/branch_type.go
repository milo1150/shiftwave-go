package types

type CreateBranchPayload struct {
	BranchName string `json:"branch_name" validate:"required,alphanum"`
}

type UpdateBranchPayload struct {
	BranchName string `json:"branch_name" validate:"omitempty,alphanum"`
	IsActive   bool   `json:"is_active" validate:"boolean"`
}
