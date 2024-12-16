package dto

import (
	"shiftwave-go/internal/models"
	"shiftwave-go/internal/types"
)

func TransformGetAssessments(assessments []models.Assessment) []types.GetAssessmentsDTO {
	transformed := []types.GetAssessmentsDTO{}
	for _, v := range assessments {
		transformed = append(transformed, types.GetAssessmentsDTO{
			ID:        v.ID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			Remark:    v.Remark,
			Score:     v.Score,
		})
	}
	return transformed
}
