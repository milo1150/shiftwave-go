package handler

import (
	"net/http"
	"shiftwave-go/internal/enum"
	"shiftwave-go/internal/model"
	"shiftwave-go/internal/types"
	"shiftwave-go/internal/utils"
	v1dto "shiftwave-go/internal/v1/dto"
	v1repo "shiftwave-go/internal/v1/repository"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type GenerateRandomReviewsParams struct {
	Lang   string    `query:"lang" validate:"required,oneof=EN TH MY"` // always keep update enum with types.Lang
	Count  int       `query:"count" validate:"required,min=1"`
	Branch uuid.UUID `query:"branch" validate:"required,uuid"`
}

// Meme function
func GenerateRandomReviews(c echo.Context, app *types.App) error {
	q := &GenerateRandomReviewsParams{}
	if err := c.Bind(q); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid params"})
	}

	v := validator.New()
	if err := v.Struct(q); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessagees := utils.ExtractErrorMessages(validationErrors)
		return c.JSON(http.StatusBadRequest, errorMessagees)
	}

	// Check is valid branch uuid
	if _, err := v1repo.FindBranchByUUID(app.DB, q.Branch); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid branch id"})
	}

	// Check is Branch existed
	if err := app.DB.Model(&model.Branch{}).Where("uuid = ?", q.Branch).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// :)
	parseLang, _ := enum.ParseLang(q.Lang)

	reviews := []v1dto.GetReviewDTO{}
	for i := 0; i < q.Count; i++ {
		randomScore := gofakeit.Number(1, 5)

		remark := ""
		switch *parseLang {
		case "TH":
			remark = "กลไกสนใจสตาร์ทอัพ พยายามจิตวิญญาณปัญญาพันธกิจสนใจ ปัญญาผสมผสานกระแสโซเชียลประโยชน์ถ่ายทอด สตาร์ทอัพต้นแบบสุขภาวะสนทนาเอสเอ็มอี สนธิกำลังบริหารในมุ่งมั่นด้วยสำคัญ สุขภาวะสนใจหลากหลายประชารัฐหรือวาทกรรม คัดเลือกซึ่งมวลชนหรือยุคใหม่กระแส พัฒนาสร้างสรรค์สุขภาวะกลไก ผสมผสานต้นแบบจิตวิญญาณมวลชน ของต้นแบบสุขภาวะจิตวิญญาณพยายามล้ำสมัย สนทนาผสมผสานคัดเลือกวิเคราะห์พัฒนาวัยรุ่น หลากหลายสูงสุดคุณธรรมจริยธรรมพันธกิจ เอสเอ็มอีพันธกิจและลงทุนคัดเลือก เอสเอ็มอีพยายามหรือ"
		case "MY":
			remark = "လူ့ဘောင်နှင့်တူသော မူလတန်းကျောင်းတွင် တက္ကသိုလ်ရောက်ကာ ကျောင်းသားများကို ဖိစီးမှုများနှင့် အတူ အသက်မွေးဝမ်းကျောင်းများ၏ အမြတ်များကို ပျက်ကွက်ခြင်းမရှိဘဲ စနစ်ကျကျ စိတ်ချမ်းသာမှု၊ အခြေခံကျတဲ့ သင်ကြားမှုများ၊ လူမှုရေးဆိုင်ရာ ဆက်ဆံရေးများကို အဓိကထားသည်။ ကျောင်းသားတစ်ဦးချင်းစီ၏ အတန်းအစားတိုးတက်မှုများနှင့် အတူ စွမ်းဆောင်ရည်နှင့် ဦးစားပေးမှုများကို ဦးတည်၍ ပညာရေးဝန်ကြီးဌာနမှ လုပ်ဆောင်နေသည်။"
		case "EN":
			remark = gofakeit.LoremIpsumSentence(30)
		}

		review := &model.Review{
			Score:      uint(randomScore),
			Remark:     remark,
			BranchUUID: q.Branch,
			Lang:       *parseLang,
		}

		// Insert to DB
		if err := app.DB.Create(review).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		// Results
		v, _ := v1dto.TransformGetReview(*review, app.ENV.LocalTimezone)
		reviews = append(reviews, v)
	}

	return c.JSON(http.StatusOK, map[string][]v1dto.GetReviewDTO{"DB gonna be ok...": reviews})
}
