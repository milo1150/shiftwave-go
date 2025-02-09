package middleware

import (
	"errors"
	"net/http"
	"shiftwave-go/internal/database"
	"strconv"
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

func IpRateLimiterMiddleware(rdb *redis.Client, limit uint64) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Generate a unique key for this IP and date
			today := time.Now().Format(time.DateOnly)
			ip := c.RealIP()
			key := database.GetRateLimitKey(ip, today)

			countStr, _ := rdb.Get(c.Request().Context(), key).Result()
			count, _ := strconv.ParseInt(countStr, 10, 64)

			if count > 0 {
				// Set TTL (Time To Live) for key if it's the first request today
				if count >= 1 {
					rdb.Expire(c.Request().Context(), key, 1*time.Hour)
				}

				// Check if the IP exceeded the daily limit
				if count >= int64(limit) {
					return c.JSON(http.StatusTooManyRequests, map[string]string{
						"error": "Rate limit exceed. Try again tomorrow.",
					})
				}
			}

			return next(c)
		}
	}
}

func IncreaseIpCounting(c echo.Context, rdb *redis.Client) error {
	today := time.Now().Format(time.DateOnly)
	ip := c.RealIP()
	key := database.GetRateLimitKey(ip, today)

	_, err := rdb.Incr(c.Request().Context(), key).Result()
	if err != nil {
		return errors.New("unable to process request")
	}

	return nil
}
