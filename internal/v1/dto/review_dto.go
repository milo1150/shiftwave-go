package dto

import (
	"shiftwave-go/internal/model"
	"shiftwave-go/internal/v1/types"
	"time"
)

func TransformGetReviews(reviews []model.Review, timezone *time.Location) []types.GetReviewDTO {
	transformed := []types.GetReviewDTO{}
	for _, review := range reviews {
		dto, _ := TransformGetReview(review, timezone)
		transformed = append(transformed, dto)
	}
	return transformed
}

func TransformGetReview(model model.Review, timezone *time.Location) (types.GetReviewDTO, error) {
	transformed := types.GetReviewDTO{
		ID:        model.ID,
		CreatedAt: model.CreatedAt.In(timezone).Format("02/01/2006 15:04"),
		UpdatedAt: model.UpdatedAt.In(timezone).Format("02/01/2006 15:04"),
		Remark:    model.Remark,
		Score:     model.Score,
	}
	return transformed, nil
}
