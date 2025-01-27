package handler

import (
	"fmt"
	"net/http"
	"shiftwave-go/internal/utils"
	v1dto "shiftwave-go/internal/v1/dto"
	v1repo "shiftwave-go/internal/v1/repository"
	v1types "shiftwave-go/internal/v1/types"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CreateBranchHandler(c echo.Context, db *gorm.DB) error {
	payload := &v1types.CreateBranchPayload{}
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid payload")
	}

	v := validator.New()
	if err := v.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessagees := utils.ExtractErrorMessages(validationErrors)
		return c.JSON(http.StatusBadRequest, errorMessagees)
	}

	if err := v1repo.CreateBranch(db, payload.BranchName); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, http.StatusOK)
}

func GetBranchesHandler(c echo.Context, db *gorm.DB) error {
	branches, err := v1repo.GetBranches(db)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	branchesDto := v1dto.TransformBranches(*branches)

	return c.JSON(http.StatusOK, branchesDto)
}

func UpdateBranchHandler(c echo.Context, db *gorm.DB) error {
	param := c.Param("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid param")
	}

	payload := &v1types.UpdateBranchPayload{}
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid payload")
	}

	v := validator.New()
	if err := v.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessagees := utils.ExtractErrorMessages(validationErrors)
		return c.JSON(http.StatusBadRequest, errorMessagees)
	}

	if err = v1repo.UpdateBranch(db, id, payload); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v, Unable to Update Branch", err)})
	}

	return c.JSON(http.StatusOK, http.StatusOK)
}
