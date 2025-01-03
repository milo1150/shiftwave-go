package types

import "time"

type GetReviewDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Remark    string    `json:"remark"`
	Score     uint      `json:"score"`
}

type ReviewQueryParams struct {
	Page     *int    `query:"page" validate:"omitempty,numeric"`
	PageSize *int    `query:"page_size" validate:"omitempty,numeric"`
	Remark   *string `query:"remark" validate:"omitempty"`
	Score    *int    `query:"score" validate:"omitempty,numeric"`
}

type ReviewsResponse struct {
	Page       int            `json:"page" validate:"omitempty,numeric"`
	PageSize   int            `json:"page_size" validate:"omitempty,numeric"`
	TotalItems int64          `json:"total_items"`
	Items      []GetReviewDTO `json:"items"`
}

type CreateReviewPayload struct {
	Remark string `json:"remark" validate:"omitempty"`
	Score  uint   `json:"score" validate:"required,min=1,max=5"`
	Branch uint   `json:"branch" validate:"required,numeric"`
}
