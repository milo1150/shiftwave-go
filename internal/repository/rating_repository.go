package repository

import (
	"shiftwave-go/internal/dto"
	"shiftwave-go/internal/model"
	"shiftwave-go/internal/types"
)

func CreateRating(app *types.App, payload *types.CreateRatingPayload) error {
	return app.DB.Create(&model.Rating{Remark: payload.Remark, Score: payload.Score}).Error
}

func GetRatings(app *types.App, q *types.RatingQueryParams) (*types.RatingsResponse, error) {
	rating := &[]model.Rating{}

	page := 1
	if q.Page != nil {
		page = *q.Page
	}

	pageSize := 10
	if q.PageSize != nil {
		pageSize = *q.PageSize
	}

	dbQuery := app.DB

	// Count
	var totalItems int64
	dbQuery.Model(&model.Rating{}).Count(&totalItems)

	// Handle Remark param
	if q.Remark != nil {
		dbQuery = dbQuery.Where("remark LIKE ?", "%"+*q.Remark+"%")
	}

	// Handle Score param
	if q.Score != nil {
		dbQuery = dbQuery.Where("score = ?", *q.Score)
	}

	// Calculate pagination
	offset := (page - 1) * pageSize
	dbQuery = dbQuery.Limit(pageSize).Offset(offset)

	// Order then Find all
	dbQuery.Order("id DESC").Find(rating)

	// Execute
	if err := dbQuery.Error; err != nil {
		return nil, err
	}

	// Transform result
	ratings := dto.TransformGetRatings(*rating)
	result := &types.RatingsResponse{
		Page:       page,
		PageSize:   pageSize,
		Items:      ratings,
		TotalItems: totalItems,
	}

	return result, nil
}

func GetRating(app *types.App, id int) (*types.GetRatingDTO, error) {
	rating := &model.Rating{}

	dbQuery := app.DB.Where("id = ?", id).First(rating)

	if err := dbQuery.Error; err != nil {
		return nil, err
	}

	result := dto.TransformGetRating(*rating)

	return &result, nil
}
