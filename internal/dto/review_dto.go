package dto

import (
	"shiftwave-go/internal/model"
	"shiftwave-go/internal/types"
)

func TransformGetReviews(reviews []model.Review) []types.GetReviewDTO {
	transformed := []types.GetReviewDTO{}
	for _, v := range reviews {
		transformed = append(transformed, TransformGetReview(v))
	}
	return transformed
}

func TransformGetReview(model model.Review) types.GetReviewDTO {
	transformed := types.GetReviewDTO{
		ID:        model.ID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		Remark:    model.Remark,
		Score:     model.Score,
	}
	return transformed
}
