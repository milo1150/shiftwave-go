package handler

import (
	"net/http"
	"shiftwave-go/internal/model"
	"shiftwave-go/internal/types"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupRoutes(e *echo.Echo, app *types.App) {
	e.GET("/generate-pdf", func(c echo.Context) error {
		return GenerateQRCodeHandler(c)
	})

	e.GET("/generate-random-reviews", func(c echo.Context) error {
		return GenerateRandomReviews(c, app.DB)
	})
}

// Mock fn
func GenerateRandomReviews(c echo.Context, db *gorm.DB) error {
	for i := 0; i < 10; i++ {
		randomScore := gofakeit.Number(1, 5)
		randomRemark := gofakeit.LoremIpsumSentence(50)
		review := &model.Review{
			Score:    uint(randomScore),
			Remark:   randomRemark,
			BranchID: 44, // TODO: dynamic value, this crash in production test.
		}
		spew.Dump(review)
		db.Create(review)
	}
	return c.JSON(http.StatusOK, "everything gonna be ok...")
}
