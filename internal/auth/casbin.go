package auth

import (
	"log"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

func InitPermission(db *gorm.DB) *casbin.Enforcer {
	// Create Casbin table into existed db
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Fatalf("Error create casbin adapter: %v", err)
	}

	// Load init config
	enforcer, err := casbin.NewEnforcer("rbac_model.conf", adapter)
	if err != nil {
		log.Fatalf("Error load enforcer: %v", err)
	}

	// Add more Route here...
	enforcer.AddPolicy("admin", "/v1/reviews")

	enforcer.LoadPolicy()

	return enforcer
}
