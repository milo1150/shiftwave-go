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
	parseLang, err := types.ParseLang(payload.Lang)
	if err != nil {
		return err
	}

	return db.Create(&model.Review{Remark: payload.Remark, Score: payload.Score, BranchID: payload.Branch, Lang: *parseLang}).Error
}

func GetReviews(app *types.App, q *v1types.ReviewQueryParams, loc time.Location) (*v1types.ReviewsResponse, error) {
	reviews := &[]model.Review{}

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

			dateQuery, _ := getDateReviewsQuery(app.DB, *q.StartDate, loc)
			dbQuery = dateQuery

		case "date_range":
			if q.StartDate == nil || q.EndDate == nil {
				return nil, fmt.Errorf("invalid start_date and end_date")
			}

			dateRangeQuery, _ := getDateRangeReviewsQuery(app.DB, *q.StartDate, *q.EndDate, loc)
			dbQuery = dateRangeQuery

		case "month":
			if q.Month == nil || q.Year == nil {
				return nil, fmt.Errorf("invalid month and year")
			}

			monthQuery, _ := getMonthReviewsQuery(app.DB, *q.Month, *q.Year, *app.ENV.LocalTimezone)
			dbQuery = monthQuery

		case "year":
			if q.Year == nil {
				return nil, fmt.Errorf("invalid year")
			}

			yearQuery, _ := getYearReviewsQuery(app.DB, *q.Year, loc)
			dbQuery = yearQuery
		}
	}

	// Count
	var totalItems int64
	dbQuery.Model(&model.Review{}).Count(&totalItems)

	// Calculate pagination
	offset := (page - 1) * pageSize
	dbQuery = dbQuery.Limit(pageSize).Offset(offset)

	// Preload Branch and execute query
	if err := dbQuery.Preload("Branch").Order("id DESC").Find(reviews).Error; err != nil {
		return nil, err
	}

	// Transform result
	reviewsDto := v1dto.TransformGetReviews(*reviews, app.ENV.LocalTimezone)
	result := &v1types.ReviewsResponse{
		Page:       page,
		PageSize:   pageSize,
		Items:      reviewsDto,
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

func GetAverageRating(db *gorm.DB, q *v1types.ReviewQueryParams, location time.Location) (*v1dto.AverageRatingDTO, error) {
	if q.DateType == nil {
		return nil, fmt.Errorf("date_type is required")
	}

	reviews := &[]model.Review{}

	switch *q.DateType {
	case "date":
		if q.StartDate == nil {
			return nil, fmt.Errorf("start_date is required")
		}

		dateQuery, _ := getDateReviewsQuery(db, *q.StartDate, location)
		if err := dateQuery.Find(reviews).Error; err != nil {
			return nil, err
		}

	case "date_range":
		if q.StartDate == nil || q.EndDate == nil {
			return nil, fmt.Errorf("start_date and end_date are required")
		}

		dateRangeQuery, _ := getDateRangeReviewsQuery(db, *q.StartDate, *q.EndDate, location)
		if err := dateRangeQuery.Find(reviews).Error; err != nil {
			return nil, err
		}

	case "month":
		if q.Month == nil || q.Year == nil {
			return nil, fmt.Errorf("month and year are required")
		}

		monthQuery, _ := getMonthReviewsQuery(db, *q.Month, *q.Year, location)
		if err := monthQuery.Find(reviews).Error; err != nil {
			return nil, err
		}

	case "year":
		if q.Year == nil {
			return nil, fmt.Errorf("year is required")
		}

		yearQuery, _ := getYearReviewsQuery(db, *q.Year, location)
		if err := yearQuery.Find(reviews).Error; err != nil {
			return nil, err
		}
	}

	result := v1dto.GetAverageRating(*reviews)

	return result, nil
}

func getYearReviewsQuery(db *gorm.DB, year int, loc time.Location) (*gorm.DB, error) {
	location, err := time.LoadLocation(loc.String())
	if err != nil {
		return nil, fmt.Errorf("invalid location")
	}

	startOfYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.Now().In(location).Location())
	endOfYear := startOfYear.AddDate(1, 0, 0).Add(-1 * time.Second)

	query := db.Where("created_at BETWEEN ? AND ?", startOfYear, endOfYear)

	return query, nil
}

func getMonthReviewsQuery(db *gorm.DB, month int, year int, loc time.Location) (*gorm.DB, error) {
	location, err := time.LoadLocation(loc.String())
	if err != nil {
		return nil, fmt.Errorf("invalid location")
	}

	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, location)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-1 * time.Second)

	query := db.Where("created_at BETWEEN ? AND ?", startOfMonth, endOfMonth)

	return query, nil
}

func getDateRangeReviewsQuery(db *gorm.DB, start string, end string, loc time.Location) (*gorm.DB, error) {
	location, err := time.LoadLocation(loc.String())
	if err != nil {
		return nil, fmt.Errorf("invalid location")
	}

	startDate, _ := time.ParseInLocation(time.DateOnly, start, location)
	endDate, _ := time.ParseInLocation(time.DateOnly, end, location)
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	query := db.Where("created_at BETWEEN ? AND ?", startDate, endDate)

	return query, nil
}

func getDateReviewsQuery(db *gorm.DB, start string, loc time.Location) (*gorm.DB, error) {
	location, err := time.LoadLocation(loc.String())
	if err != nil {
		return nil, fmt.Errorf("invalid location")
	}

	startDate, _ := time.ParseInLocation(time.DateOnly, start, location)
	endDate := startDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	query := db.Where("created_at BETWEEN ? AND ?", startDate, endDate)

	return query, nil
}

func RetrieveReviewsByLang(db *gorm.DB, lang types.Lang, loc time.Location, duration time.Duration) (*[]model.Review, error) {
	location, err := time.LoadLocation(loc.String())
	if err != nil {
		return nil, fmt.Errorf("invalid location")
	}

	reviews := &[]model.Review{}
	currentTime := time.Now().In(location)
	startTime := currentTime.Add(-10 * time.Hour)

	query := db.Where("created_at BETWEEN ? AND ? AND lang = ?", startTime, currentTime, lang)

	if err := query.Find(reviews).Error; err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	return reviews, nil
}
