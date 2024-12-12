package types

type CreateAssessmentPayload struct {
	Remark string `json:"remark" validate:"required"`
	Score  uint   `json:"score" validate:"required,min=1,max=10"`
}
