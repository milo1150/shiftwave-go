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

	// Add Route Guard
	addPolicy(enforcer)

	// Remove unuse policy
	removePolicy(enforcer)

	// LoadPolicy reloads the policy from database.
	enforcer.LoadPolicy()

	return enforcer
}

func addPolicy(enforcer *casbin.Enforcer) {
	// Role - Admin
	enforcer.AddPolicy("admin", "/v1/reviews")
	enforcer.AddPolicy("admin", "/v1/reviews/average-rating")
	enforcer.AddPolicy("admin", "/v1/reviews/sse")

	// Role - User
	enforcer.AddPolicy("user", "/v1/reviews")
	enforcer.AddPolicy("user", "/v1/reviews/average-rating")
	enforcer.AddPolicy("user", "/v1/reviews/sse")
}

func removePolicy(enforcer *casbin.Enforcer) {
	enforcer.RemovePolicy("admin", "/v1/average-rating")
}
