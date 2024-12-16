package repositories

import (
	"shiftwave-go/internal/dto"
	"shiftwave-go/internal/models"
	"shiftwave-go/internal/types"
)

func CreateAssessment(app *types.App, payload *types.CreateAssessmentPayload) error {
	return app.DB.Create(&models.Assessment{Remark: payload.Remark, Score: payload.Score}).Error
}

func GetAssessments(app *types.App, q *types.AssessmentQueryParams) (*types.AssessmentsResponse, error) {
	assessment := &[]models.Assessment{}

	page := 1
	if q.Page != nil {
		page = *q.Page
	}

	pageSize := 10
	if q.PageSize != nil {
		pageSize = *q.PageSize
	}

	dbQuery := app.DB

	if q.Remark != nil {
		dbQuery = dbQuery.Where("remark LIKE ?", "%"+*q.Remark+"%")
	}

	if q.Score != nil {
		dbQuery = dbQuery.Where("score = ?", *q.Score)
	}

	// Add pagination
	offset := (page - 1) * pageSize
	dbQuery = dbQuery.Limit(pageSize).Offset(offset)

	var totalItems int64
	dbQuery.Model(&models.Assessment{}).Count(&totalItems)

	if err := dbQuery.Find(assessment).Error; err != nil {
		return nil, err
	}

	assessments := dto.TransformGetAssessments(*assessment)

	result := &types.AssessmentsResponse{
		Page:       page,
		PageSize:   pageSize,
		Items:      assessments,
		TotalItems: totalItems,
	}

	return result, nil
}
