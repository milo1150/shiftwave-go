package dto

import (
	"shiftwave-go/internal/model"
	mainTypes "shiftwave-go/internal/types"
	"shiftwave-go/internal/utils"
	"time"
)

type GetReviewDTO struct {
	ID        uint           `json:"id"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
	Remark    string         `json:"remark"`
	Score     uint           `json:"score"`
	Lang      mainTypes.Lang `json:"lang"`
	RemarkEn  string         `json:"remark_en"`
}

func TransformGetReviews(reviews []model.Review, timezone *time.Location) []GetReviewDTO {
	transformed := []GetReviewDTO{}
	for _, review := range reviews {
		dto, _ := TransformGetReview(review, timezone)
		transformed = append(transformed, dto)
	}
	return transformed
}

func TransformGetReview(model model.Review, timezone *time.Location) (GetReviewDTO, error) {
	transformed := GetReviewDTO{
		ID:        model.ID,
		CreatedAt: model.CreatedAt.In(timezone).Format("02/01/2006 15:04"),
		UpdatedAt: model.UpdatedAt.In(timezone).Format("02/01/2006 15:04"),
		Remark:    model.Remark,
		Score:     model.Score,
		Lang:      model.Lang,
		RemarkEn:  model.RemarkEn,
	}
	return transformed, nil
}

type AverageRatingDTO struct {
	TotalCount       int     `json:"total_review"`
	AverageRating    float64 `json:"average_rating"`
	FiveStarCount    int     `json:"five_star_count"`
	FiveStarPercent  float64 `json:"five_star_percent"`
	FourStarCount    int     `json:"four_star_count"`
	FourStarPercent  float64 `json:"four_star_percent"`
	ThreeStarCount   int     `json:"three_star_count"`
	ThreeStarPercent float64 `json:"three_star_percent"`
	TwoStarCount     int     `json:"two_star_count"`
	TwoStarPercent   float64 `json:"two_star_percent"`
	OneStarCount     int     `json:"one_star_count"`
	OneStarPercent   float64 `json:"one_star_percent"`
}

func GetAverageRating(reviews []model.Review) *AverageRatingDTO {
	result := &AverageRatingDTO{}

	if len(reviews) <= 0 {
		return result
	}

	var (
		oneScores   []model.Review
		twoScores   []model.Review
		threeScores []model.Review
		fourScores  []model.Review
		fiveScores  []model.Review
	)

	for _, review := range reviews {
		switch review.Score {
		case 1:
			oneScores = append(oneScores, review)

		case 2:
			twoScores = append(twoScores, review)

		case 3:
			threeScores = append(threeScores, review)

		case 4:
			fourScores = append(fourScores, review)

		case 5:
			fiveScores = append(fiveScores, review)
		}
	}

	// Parse
	totalCount := float64(len(reviews))
	fiveScoreCount := float64(len(fiveScores))
	fourScoreCount := float64(len(fourScores))
	threeScoreCount := float64(len(threeScores))
	twoScoreCount := float64(len(twoScores))
	oneScoreCount := float64(len(oneScores))

	// Set counting value
	result.TotalCount = int(totalCount)
	result.OneStarCount = len(oneScores)
	result.TwoStarCount = len(twoScores)
	result.ThreeStarCount = len(threeScores)
	result.FourStarCount = len(fourScores)
	result.FiveStarCount = len(fiveScores)

	// Calculate percentage
	result.FiveStarPercent = utils.RoundFloat64(fiveScoreCount / totalCount * 100)
	result.FourStarPercent = utils.RoundFloat64(fourScoreCount / totalCount * 100)
	result.ThreeStarPercent = utils.RoundFloat64(threeScoreCount / totalCount * 100)
	result.TwoStarPercent = utils.RoundFloat64(twoScoreCount / totalCount * 100)
	result.OneStarPercent = utils.RoundFloat64(oneScoreCount / totalCount * 100)

	// Calculate Average
	averageRating := ((5 * fiveScoreCount) + (4 * fourScoreCount) + (3 * threeScoreCount) + (2 * twoScoreCount) + oneScoreCount) / totalCount
	result.AverageRating = utils.RoundToTwoDecimals(averageRating)

	return result
}
