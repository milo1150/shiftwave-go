package repository

import (
	"fmt"
	"shiftwave-go/internal/model"
	"shiftwave-go/internal/types"
	v1dto "shiftwave-go/internal/v1/dto"
	v1types "shiftwave-go/internal/v1/types"
	"time"

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

	// Handle Remark param
	if q.Remark != nil {
		dbQuery = dbQuery.Where("remark LIKE ?", "%"+*q.Remark+"%")
	}

	// Handle Score param
	if q.Score != nil {
		dbQuery = dbQuery.Where("score = ?", *q.Score)
	}

	// Handle date_type (oneof="date date_range month year")
	if q.DateType != nil {
		switch *q.DateType {
		case "date":
			if q.StartDate == nil {
				return nil, fmt.Errorf("invalid start_date")
			}
			startDate, _ := time.Parse(time.DateOnly, *q.StartDate)
			endDate := startDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			dbQuery = dbQuery.Where("created_at BETWEEN ? AND ?", startDate, endDate)

		case "date_range":
			if q.StartDate == nil || q.EndDate == nil {
				return nil, fmt.Errorf("invalid start_date and end_date")
			}
			startDate, _ := time.Parse(time.DateOnly, *q.StartDate)
			endDate, _ := time.Parse(time.DateOnly, *q.EndDate)
			endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			dbQuery = dbQuery.Where("created_at BETWEEN ? AND ?", startDate, endDate)

		case "month":
			if q.Month == nil {
				return nil, fmt.Errorf("invalid month")
			}
			month := *q.Month
			startOfMonth := time.Date(time.Now().Year(), time.Month(month), 1, 0, 0, 0, 0, time.Now().Location())
			endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-1 * time.Second)
			fmt.Println("startOfMonth", startOfMonth)
			fmt.Println("endOfMonth", endOfMonth)
			dbQuery = dbQuery.Where("created_at BETWEEN ? AND ?", startOfMonth, endOfMonth)

		case "year":
			if q.Year == nil {
				return nil, fmt.Errorf("invalid year")
			}
			year := *q.Year
			startOfYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.Now().Location())
			endOfYear := startOfYear.AddDate(1, 0, 0).Add(-1 * time.Second)
			dbQuery = dbQuery.Where("created_at BETWEEN ? AND ?", startOfYear, endOfYear)
		}
	}

	// Calculate pagination
	offset := (page - 1) * pageSize
	dbQuery = dbQuery.Limit(pageSize).Offset(offset)

	// Count
	var totalItems int64
	dbQuery.Model(&model.Review{}).Count(&totalItems)

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

func GetReview(app *types.App, id int) (*v1dto.GetReviewDTO, error) {
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
