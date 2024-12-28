package dto

import (
	"shiftwave-go/internal/model"
	"shiftwave-go/internal/types"
)

func TransformGetRatings(ratings []model.Rating) []types.GetRatingDTO {
	transformed := []types.GetRatingDTO{}
	for _, v := range ratings {
		transformed = append(transformed, TransformGetRating(v))
	}
	return transformed
}

func TransformGetRating(model model.Rating) types.GetRatingDTO {
	transformed := types.GetRatingDTO{
		ID:        model.ID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		Remark:    model.Remark,
		Score:     model.Score,
	}
	return transformed
}
