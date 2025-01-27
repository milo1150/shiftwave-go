package types

import "github.com/google/uuid"

type GeneratePdfParams struct {
	BranchUuid uuid.UUID `query:"branch_uuid" validate:"required"`
}
