package types

type CreateBranch struct {
	BranchName string `json:"branch_name" validate:"required,alphanum"`
}

type UpdateBranch struct {
	BranchName string `json:"branch_name" validate:"omitempty,alphanum"`
	IsActive   bool   `json:"is_active" validate:"boolean"`
}
