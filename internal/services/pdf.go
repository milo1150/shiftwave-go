package services

import (
	"fmt"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/code"
	"github.com/johnfercher/maroto/v2/pkg/core"
)

func GenerateReviewQRcode(baseUrl string, branchId uint16) core.Maroto {
	m := maroto.New()

	url := fmt.Sprintf("%s/?branch_id=%d", baseUrl, branchId)

	fmt.Println(url)

	m.AddRow(60,
		code.NewQrCol(100, url),
	)

	return m
}
