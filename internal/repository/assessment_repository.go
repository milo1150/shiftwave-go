package repository

import (
	"shiftwave-go/internal/dto"
	"shiftwave-go/internal/model"
	"shiftwave-go/internal/types"
)

func CreateAssessment(app *types.App, payload *types.CreateAssessmentPayload) error {
	return app.DB.Create(&model.Assessment{Remark: payload.Remark, Score: payload.Score}).Error
}

func GetAssessments(app *types.App, q *types.AssessmentQueryParams) (*types.AssessmentsResponse, error) {
	assessment := &[]model.Assessment{}

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
	dbQuery.Model(&model.Assessment{}).Count(&totalItems)

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
	dbQuery.Order("id DESC").Find(assessment)

	// Execute
	if err := dbQuery.Error; err != nil {
		return nil, err
	}

	// Transform result
	assessments := dto.TransformGetAssessments(*assessment)
	result := &types.AssessmentsResponse{
		Page:       page,
		PageSize:   pageSize,
		Items:      assessments,
		TotalItems: totalItems,
	}

	return result, nil
}

func GetAssessment(app *types.App, id int) (*types.GetAssessmentDTO, error) {
	assessment := &model.Assessment{}

	dbQuery := app.DB.Where("id = ?", id).First(assessment)

	if err := dbQuery.Error; err != nil {
		return nil, err
	}

	result := dto.TransformGetAssessment(*assessment)

	return &result, nil
}
