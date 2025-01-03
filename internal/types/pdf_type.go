package types

type GeneratePdfParams struct {
	BranchId uint16 `query:"branch_id" validate:"required,min=1"`
}
