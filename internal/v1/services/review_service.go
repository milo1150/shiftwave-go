package services

import (
	"fmt"
	"shiftwave-go/internal/model"
	"shiftwave-go/internal/utils"

	v1types "shiftwave-go/internal/v1/types"
)

func GetAverageRating(reviews []model.Review) (*v1types.AverageRatingResponse, error) {
	if len(reviews) <= 0 {
		return nil, fmt.Errorf("data not found")
	}

	var (
		oneScores   []model.Review
		twoScores   []model.Review
		threeScores []model.Review
		fourScores  []model.Review
		fiveScores  []model.Review
	)
	result := &v1types.AverageRatingResponse{}

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
	result.AverageRating = utils.RoundFloat64(((5 * fiveScoreCount) + (4 * fourScoreCount) + (3 * threeScoreCount) + (2 * twoScoreCount) + oneScoreCount) / totalCount)

	return result, nil
}
