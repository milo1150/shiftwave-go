package database

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"shiftwave-go/internal/model"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func loadMasterDataJsonFile() []byte {
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get dir: %v", err)
	}

	masterDataFilePath := filepath.Join(basePath, "internal", "resources", "master_data.json")

	masterDataByte, err := os.ReadFile(masterDataFilePath)
	if err != nil {
		log.Fatalf("Failed to load master_data.json file: %v", err)
	}

	return masterDataByte
}

func getMasterDataJson(masterDataByte []byte) *model.MasterDataJson {
	masterDataJson := &model.MasterDataJson{}
	if err := json.Unmarshal(masterDataByte, &masterDataJson); err != nil {
		log.Fatalf("Failed to init masterDataJson: %v", err)
	}

	validate := validator.New()
	if err := validate.Struct(masterDataJson); err != nil {
		log.Fatalf("Invalid master_data.json structure: %v", err)
	}

	return masterDataJson
}

func insertBranchesIntoDB(db *gorm.DB, branchJsons []model.BranchMasterData) {
	branches := &[]model.Branch{}
	if err := db.Find(branches).Error; err != nil {
		log.Fatalf("Branches not found")
	}

	check := make(map[string]string)
	for _, branch := range *branches {
		check[branch.Name] = ""
	}

	for _, branch := range branchJsons {
		if _, existed := check[branch.Name]; !existed {
			if err := db.Create(&model.Branch{Name: branch.Name}).Error; err != nil {
				log.Fatalf("Failed to inserted %v into Branch table: %v", branch.Name, err)
			} else {
				log.Printf("Inserted %v into Branch table.", branch.Name)
			}
		}
	}
}

func MasterDataLoader(db *gorm.DB) {
	masterDataByte := loadMasterDataJsonFile()

	masterDataJson := getMasterDataJson(masterDataByte)

	insertBranchesIntoDB(db, masterDataJson.Branches)
}
