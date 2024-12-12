package repositories

import (
	"shiftwave-go/internal/models"
	"shiftwave-go/internal/types"
)

func CreateAssessment(app *types.App, payload *types.CreateAssessmentPayload) error {
	return app.DB.Create(&models.Assessment{Remark: payload.Remark, Score: payload.Score}).Error
}
