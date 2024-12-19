package dto

import (
	"shiftwave-go/internal/model"
	"shiftwave-go/internal/types"
)

func TransformGetAssessments(assessments []model.Assessment) []types.GetAssessmentDTO {
	transformed := []types.GetAssessmentDTO{}
	for _, v := range assessments {
		transformed = append(transformed, TransformGetAssessment(v))
	}
	return transformed
}

func TransformGetAssessment(model model.Assessment) types.GetAssessmentDTO {
	transformed := types.GetAssessmentDTO{
		ID:        model.ID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		Remark:    model.Remark,
		Score:     model.Score,
	}
	return transformed
}
