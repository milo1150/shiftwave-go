package types

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Env struct {
	LocalTimezone *time.Location
	OpenAI        string
	JWT           string
	AdminPassword string
}

type App struct {
	DB      *gorm.DB
	ENV     Env
	RDB     *redis.Client
	Context context.Context
}

type AnyObject map[string]interface{}
