package auth

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
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

	// Register KeyMatch2 function
	// See https://casbin.org/docs/function/ and https://casbin.org/docs/rbac-with-pattern
	enforcer.AddNamedMatchingFunc("g", "KeyMatch2", util.KeyMatch2)

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
	enforcer.AddPolicy("admin", "/generate-pdf")
	enforcer.AddPolicy("admin", "/v1/reviews")
	enforcer.AddPolicy("admin", "/v1/reviews/average-rating")
	enforcer.AddPolicy("admin", "/v1/reviews/sse")
	enforcer.AddPolicy("admin", "/v1/branches")
	enforcer.AddPolicy("admin", "/v1/branch/:id")
	enforcer.AddPolicy("admin", "/v1/user/create-user")
	enforcer.AddPolicy("admin", "/v1/user/get-users")
	enforcer.AddPolicy("admin", "/v1/user/update-users")

	// Role - User
	enforcer.AddPolicy("user", "/v1/reviews")
	enforcer.AddPolicy("user", "/v1/reviews/average-rating")
	enforcer.AddPolicy("user", "/v1/reviews/sse")
}

func removePolicy(enforcer *casbin.Enforcer) {
	// enforcer.RemovePolicy("admin", "/v1/average-rating")
}
