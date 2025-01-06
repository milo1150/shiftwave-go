package repository

import (
	"shiftwave-go/internal/dto"
	"shiftwave-go/internal/model"
	"shiftwave-go/internal/types"

	"gorm.io/gorm"
)

func CreateReview(db *gorm.DB, payload *types.CreateReviewPayload) error {
	return db.Create(&model.Review{Remark: payload.Remark, Score: payload.Score, BranchID: payload.Branch}).Error
}

func GetReviews(app *types.App, q *types.ReviewQueryParams) (*types.ReviewsResponse, error) {
	review := &[]model.Review{}

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
	dbQuery.Model(&model.Review{}).Count(&totalItems)

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

	// Preload Branch and execute query
	dbQuery.Preload("Branch").Order("id DESC").Find(review)

	// Execute
	if err := dbQuery.Error; err != nil {
		return nil, err
	}

	// Transform result
	reviews := dto.TransformGetReviews(*review, app.ENV.LocalTimezone)
	result := &types.ReviewsResponse{
		Page:       page,
		PageSize:   pageSize,
		Items:      reviews,
		TotalItems: totalItems,
	}

	return result, nil
}

func GetReview(app *types.App, id int) (*types.GetReviewDTO, error) {
	review := &model.Review{}

	dbQuery := app.DB.Where("id = ?", id).First(review)

	if err := dbQuery.Error; err != nil {
		return nil, err
	}

	result, err := dto.TransformGetReview(*review, app.ENV.LocalTimezone)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
