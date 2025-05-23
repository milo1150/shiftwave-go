package setup

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"shiftwave-go/internal/auth"
	"shiftwave-go/internal/model"
	"shiftwave-go/internal/types"
	v1repo "shiftwave-go/internal/v1/repository"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

func getMasterDataJson(masterDataByte []byte) *types.MasterDataJson {
	masterDataJson := &types.MasterDataJson{}
	if err := json.Unmarshal(masterDataByte, &masterDataJson); err != nil {
		log.Fatalf("Failed to init masterDataJson: %v", err)
	}

	validate := validator.New()
	if err := validate.Struct(masterDataJson); err != nil {
		log.Fatalf("Invalid master_data.json structure: %v", err)
	}

	return masterDataJson
}

func insertBranchesIntoDB(db *gorm.DB, branchJsons []types.BranchMasterData) {
	branches := &[]model.Branch{}
	if err := db.Find(branches).Error; err != nil {
		log.Fatalf("Error to find branches: %s", err)
	}

	check := make(map[string]string)
	for _, branch := range *branches {
		check[branch.Name] = ""
	}

	for _, branch := range branchJsons {
		if _, existed := check[branch.Name]; !existed {
			if err := db.Create(&model.Branch{Name: branch.Name, Uuid: uuid.New(), IsActive: true}).Error; err != nil {
				log.Fatalf("Failed to inserted %v into Branch table: %v", branch.Name, err)
			}
		}
	}
}

func insertUserIntoDB(db *gorm.DB, userJsons []types.UserMasterData, adminPassword string) {
	users := &[]model.User{}
	if err := db.Find(users).Error; err != nil {
		log.Fatalf("Find users error: %s", err)
	}

	check := map[string]string{}
	for _, user := range *users {
		check[user.Username] = ""
	}

	for _, userJson := range userJsons {
		if _, existed := check[userJson.Username]; !existed {
			// Hash adminpassword before use
			hashPassword, err := auth.HashedPassword(adminPassword)
			if err != nil {
				log.Fatalf("Failed to hash admin password")
			}

			result := db.Create(&model.User{Username: userJson.Username, Role: userJson.Role, Password: hashPassword})
			if err := result.Error; err != nil {
				log.Fatalf("Failed to inserted %v into User table: %v", userJson.Username, err)
			}
		}
	}
}

func MasterDataLoader(app *types.App) {
	masterDataByte := loadMasterDataJsonFile()

	masterDataJson := getMasterDataJson(masterDataByte)

	insertBranchesIntoDB(app.DB, masterDataJson.Branches)

	insertUserIntoDB(app.DB, masterDataJson.Users, app.ENV.AdminPassword)

	if err := v1repo.UpdateAllAdminUserBranches(app.DB); err != nil {
		log.Fatalf("Error UpdateAdminUserBranches: %v", err.Error())
	}
}
