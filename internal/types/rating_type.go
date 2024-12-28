package types

import "time"

type GetRatingDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Remark    string    `json:"remark"`
	Score     uint      `json:"score"`
}

type RatingQueryParams struct {
	Page     *int    `query:"page" validate:"omitempty,numeric"`
	PageSize *int    `query:"page_size" validate:"omitempty,numeric"`
	Remark   *string `query:"remark" validate:"omitempty"`
	Score    *int    `query:"score" validate:"omitempty,numeric"`
}

type RatingsResponse struct {
	Page       int            `json:"page" validate:"omitempty,numeric"`
	PageSize   int            `json:"page_size" validate:"omitempty,numeric"`
	TotalItems int64          `json:"total_items"`
	Items      []GetRatingDTO `json:"items"`
}

type CreateRatingPayload struct {
	Remark string `json:"remark" validate:"required"`
	Score  uint   `json:"score" validate:"required,min=1,max=10"`
}
