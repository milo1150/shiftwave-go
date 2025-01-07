package repository

import (
	"shiftwave-go/internal/model"
	"shiftwave-go/internal/types"
	v1dto "shiftwave-go/internal/v1/dto"
	v1types "shiftwave-go/internal/v1/types"

	"gorm.io/gorm"
)

func CreateReview(db *gorm.DB, payload *v1types.CreateReviewPayload) error {
	return db.Create(&model.Review{Remark: payload.Remark, Score: payload.Score, BranchID: payload.Branch}).Error
}

func GetReviews(app *types.App, q *v1types.ReviewQueryParams) (*v1types.ReviewsResponse, error) {
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
	reviews := v1dto.TransformGetReviews(*review, app.ENV.LocalTimezone)
	result := &v1types.ReviewsResponse{
		Page:       page,
		PageSize:   pageSize,
		Items:      reviews,
		TotalItems: totalItems,
	}

	return result, nil
}

func GetReview(app *types.App, id int) (*v1types.GetReviewDTO, error) {
	review := &model.Review{}

	dbQuery := app.DB.Where("id = ?", id).First(review)

	if err := dbQuery.Error; err != nil {
		return nil, err
	}

	result, err := v1dto.TransformGetReview(*review, app.ENV.LocalTimezone)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
