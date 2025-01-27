package types

import "shiftwave-go/internal/enum"

// Use in resources/master_data.json
type MasterDataJson struct {
	Branches []BranchMasterData `json:"branches" validate:"required,dive"`
	Users    []UserMasterData   `json:"users" validate:"required,dive"`
}

type BranchMasterData struct {
	Name string `json:"name" validate:"required"`
}

type UserMasterData struct {
	Username string    `json:"username" validate:"required"`
	Role     enum.Role `json:"role"`
}
