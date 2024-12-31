package model

// resources/master_data.json
type MasterDataJson struct {
	Branches []BranchMasterData `json:"branches" validate:"required,dive"`
}

type BranchMasterData struct {
	Name string `json:"name" validate:"required"`
}
