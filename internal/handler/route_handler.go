package handler

import (
	"net/http"
	"shiftwave-go/internal/model"
	"shiftwave-go/internal/types"
	v1dto "shiftwave-go/internal/v1/dto"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupRoutes(e *echo.Echo, app *types.App) {
	e.GET("/generate-pdf", func(c echo.Context) error {
		return GenerateQRCodeHandler(c)
	})

	e.GET("/generate-random-reviews", func(c echo.Context) error {
		return GenerateRandomReviews(c, app)
	})
}

type GenerateRandomReviewsParams struct {
	BranchId uint `query:"branch_id" validate:"required,number"`
}

// Mock fn
func GenerateRandomReviews(c echo.Context, app *types.App) error {
	q := &GenerateRandomReviewsParams{}
	if err := c.Bind(q); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid params"})
	}

	v := validator.New()
	if err := v.Struct(q); err != nil {
		return c.JSON(http.StatusBadRequest, "Check your branch please.")
	}

	if err := app.DB.First(&model.Branch{Model: gorm.Model{ID: q.BranchId}}).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	reviews := []v1dto.GetReviewDTO{}
	for i := 0; i < 15; i++ {
		randomScore := gofakeit.Number(1, 5)
		randomRemark := gofakeit.LoremIpsumSentence(30)
		review := &model.Review{
			Score:    uint(randomScore),
			Remark:   randomRemark,
			BranchID: q.BranchId,
		}
		v, _ := v1dto.TransformGetReview(*review, app.ENV.LocalTimezone)
		reviews = append(reviews, v)
		app.DB.Create(review)
	}

	return c.JSON(http.StatusOK, map[string][]v1dto.GetReviewDTO{"DB gonna be ok...": reviews})
}
