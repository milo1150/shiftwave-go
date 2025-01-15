package dto

import "shiftwave-go/internal/model"

type BranchDto struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
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
		Id:       int(branch.ID),
		Name:     branch.Name,
		IsActive: branch.IsActive,
	}
}
