package dto

import (
	"shiftwave-go/internal/model"
	"time"
)

type GetReviewDTO struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Remark    string `json:"remark"`
	Score     uint   `json:"score"`
}

func TransformGetReviews(reviews []model.Review, timezone *time.Location) []GetReviewDTO {
	transformed := []GetReviewDTO{}
	for _, review := range reviews {
		dto, _ := TransformGetReview(review, timezone)
		transformed = append(transformed, dto)
	}
	return transformed
}

func TransformGetReview(model model.Review, timezone *time.Location) (GetReviewDTO, error) {
	transformed := GetReviewDTO{
		ID:        model.ID,
		CreatedAt: model.CreatedAt.In(timezone).Format("02/01/2006 15:04"),
		UpdatedAt: model.UpdatedAt.In(timezone).Format("02/01/2006 15:04"),
		Remark:    model.Remark,
		Score:     model.Score,
	}
	return transformed, nil
}
