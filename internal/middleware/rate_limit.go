package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

func ConfigRateLimiter() echo.MiddlewareFunc {
	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(20), Burst: 20, ExpiresIn: 1 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, "Too many request")
		},
	}

	return middleware.RateLimiterWithConfig(config)
}

func IpRateLimiterMiddleware(rdb *redis.Client, ctx context.Context, limit uint64) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Generate a unique key for this IP and date
			today := time.Now().Format("2006-01-02")
			key := fmt.Sprintf("rate_limit:%s:%s", c.RealIP(), today)

			// Increment request count
			count, err := rdb.Incr(ctx, key).Result()
			if err != nil {
				fmt.Println(err)
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Unable to process request",
				})
			}

			// Set TTL for key if it's the first request today
			if count >= 1 {
				rdb.Expire(ctx, key, 2*time.Second)
			}

			// Check if the IP exceeded the daily limit
			if count > int64(limit) {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "Rate limit exceed. Try again tomorrow.",
				})
			}

			return next(c)
		}
	}
}
