package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/code"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/core/entity"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/johnfercher/maroto/v2/pkg/repository"
)

// not working, got generate() pdf error when use cfg.
func SetupFont() *entity.Config {
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get dir: %v", err)
	}

	fontFile := filepath.Join(basePath, "internal", "assets/Kanit-ExtraBold.ttf")

	customFont := "Kanit"

	customFonts, err := repository.New().
		AddUTF8Font(customFont, fontstyle.Bold, fontFile).
		Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	builder := config.NewBuilder().
		WithCustomFonts(customFonts)

	cfg := builder.WithDefaultFont(&props.Font{Family: customFont}).
		Build()

	return cfg
}

func GenerateReviewQRcode(appUrl string, branchUuid uuid.UUID) core.Maroto {
	m := maroto.New()

	thUrl := fmt.Sprintf("%s/th/reviews?branch_id=%s", appUrl, branchUuid.String())
	enUrl := fmt.Sprintf("%s/en/reviews?branch_id=%s", appUrl, branchUuid.String())
	myUrl := fmt.Sprintf("%s/my/reviews?branch_id=%s", appUrl, branchUuid.String())

	thText := "Please provide feedback (Thai)"
	enText := "Please provide feedback (English)"
	myText := "Please provide feedback (Myanmar)"

	m.AddRow(1,
		text.NewCol(6, thText, props.Text{Left: 0, Size: 16}),
		text.NewCol(6, enText, props.Text{Left: 0, Size: 16}),
	)

	m.AddRow(100,
		code.NewQrCol(5, thUrl, props.Rect{Center: true}),
		text.NewCol(1, ""),
		code.NewQrCol(5, enUrl, props.Rect{Left: 3, Top: 10}),
	)

	m.AddRow(10)

	m.AddRow(10,
		text.NewCol(7, myText, props.Text{Left: 0, Size: 16}),
	)

	m.AddRow(80,
		code.NewQrCol(5, myUrl, props.Rect{Left: 5}),
	)

	return m
}
