package repositories

import (
	"shiftwave-go/internal/models"
	"shiftwave-go/internal/types"
)

func CreateAssessment(app *types.App) error {
	return app.DB.Create(&models.Assessment{Remark: "test", Score: 50}).Error
	// result := app.DB.Create(&models.Assessment{Remark: "test remark", Score: 50})
	// if result.Error != nil {
	// 	return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
	// }
	// return c.JSON(http.StatusOK, "Create OK")
}
