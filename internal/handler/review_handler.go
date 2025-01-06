package handler

import (
	"net/http"
	"shiftwave-go/internal/model"
	"shiftwave-go/internal/repository"
	"shiftwave-go/internal/types"
	"shiftwave-go/internal/utils"
	"strconv"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CreateReviewHandler(c echo.Context, db *gorm.DB) error {
	payload := new(types.CreateReviewPayload)
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	v := validator.New()
	if err := v.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessagees := utils.ExtractErrorMessages(validationErrors)
		return c.JSON(http.StatusBadRequest, errorMessagees)
	}

	if result := repository.CreateReview(db, payload); result != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error()})
	}

	return c.JSON(http.StatusOK, http.StatusOK)
}

func GetReviewsHandler(c echo.Context, app *types.App) error {
	q := &types.ReviewQueryParams{}
	if err := c.Bind(q); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Query")
	}

	v := validator.New()
	if err := v.Struct(q); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := utils.ExtractErrorMessages(validationErrors)
		return c.JSON(http.StatusBadRequest, errorMessages)
	}

	result, err := repository.GetReviews(app, q)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Query error")
	}

	return c.JSON(http.StatusOK, result)
}

func GetReviewHandler(c echo.Context, app *types.App) error {
	param := c.Param("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid param")
	}

	result, _ := repository.GetReview(app, id)

	return c.JSON(http.StatusOK, result)
}

func GenerateRandomReviews(c echo.Context, db *gorm.DB) error {
	for i := 0; i < 10; i++ {
		randomScore := gofakeit.Number(1, 5)
		randomRemark := gofakeit.LoremIpsumSentence(50)
		review := &model.Review{
			Score:    uint(randomScore),
			Remark:   randomRemark,
			BranchID: 44,
		}
		spew.Dump(review)
		db.Create(review)
	}
	return c.JSON(http.StatusOK, "everything gonna be ok...")
}
