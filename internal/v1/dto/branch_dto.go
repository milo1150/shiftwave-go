package dto

import (
	"shiftwave-go/internal/model"

	"github.com/google/uuid"
)

type BranchDto struct {
	Uuid     uuid.UUID `json:"uuid"`
	Name     string    `json:"name"`
	IsActive bool      `json:"is_active"`
}

func TransformBranches(branches []model.Branch) []BranchDto {
	if len(branches) == 0 {
		return []BranchDto{}
	}

	branchesDto := make([]BranchDto, len(branches))
	for i, branch := range branches {
		branchesDto[i] = TransformBranch(branch)
	}

	return branchesDto
}

func TransformBranch(branch model.Branch) BranchDto {
	return BranchDto{
		Uuid:     branch.Uuid,
		Name:     branch.Name,
		IsActive: branch.IsActive,
	}
}
