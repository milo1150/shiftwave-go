package types

import (
	v1dto "shiftwave-go/internal/v1/dto"
)

type CreateReviewPayload struct {
	Remark string `json:"remark" validate:"omitempty"`
	Score  uint   `json:"score" validate:"required,min=1,max=5"`
	Branch uint   `json:"branch" validate:"required,numeric"`
}

type ReviewQueryParams struct {
	Page      *int    `query:"page" validate:"omitempty,numeric"`
	PageSize  *int    `query:"page_size" validate:"omitempty,numeric"`
	Remark    *string `query:"remark" validate:"omitempty"`
	Score     *int    `query:"score" validate:"omitempty,numeric"`
	DateType  *string `query:"date_type" validate:"omitempty,oneof=date date_range month year"`
	StartDate *string `query:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate   *string `query:"end_date" validate:"omitempty,datetime=2006-01-02"`
	Month     *int    `query:"month" validate:"omitempty,min=1,max=12"`
	Year      *int    `query:"year" validate:"omitempty,numeric"`
}

type ReviewsResponse struct {
	Page       int                  `json:"page" validate:"numeric"`
	PageSize   int                  `json:"page_size" validate:"numeric"`
	TotalItems int64                `json:"total_items"`
	Items      []v1dto.GetReviewDTO `json:"items"`
}

type ReviewSummaryResponse struct {
	TotalReviews  int     `json:"total_reviews"`
	AverageRating float32 `json:"average_rating"`
	FiveStar      int     `json:"five_star"`
	FourStar      int     `json:"four_star"`
	ThreeStar     int     `json:"three_star"`
	TwoStar       int     `json:"two_star"`
	OneStar       int     `json:"one_star"`
}
