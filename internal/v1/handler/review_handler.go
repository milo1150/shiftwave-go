package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"shiftwave-go/internal/connection"
	"shiftwave-go/internal/database"
	"shiftwave-go/internal/middleware"
	"shiftwave-go/internal/types"
	"shiftwave-go/internal/utils"
	v1repo "shiftwave-go/internal/v1/repository"
	v1types "shiftwave-go/internal/v1/types"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func CreateReviewHandler(c echo.Context, db *gorm.DB, rdb *redis.Client) error {
	payload := new(v1types.CreateReviewPayload)
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	v := validator.New()
	if err := v.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessagees := utils.ExtractErrorMessages(validationErrors)
		return c.JSON(http.StatusBadRequest, errorMessagees)
	}

	// Check is valid branch uuid
	if _, err := v1repo.FindBranchByUUID(db, payload.Branch); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid branch id"})
	}

	// Create Review
	if err := v1repo.CreateReview(db, payload); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Increment Redis rate limit counting
	if err := middleware.IncreaseIpCounting(c, rdb); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, http.StatusOK)
}

func GetReviewsHandler(c echo.Context, app *types.App) error {
	q := &v1types.ReviewQueryParams{}
	if err := c.Bind(q); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Query")
	}

	v := validator.New()
	if err := v.Struct(q); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := utils.ExtractErrorMessages(validationErrors)
		return c.JSON(http.StatusBadRequest, errorMessages)
	}

	result, err := v1repo.GetReviews(app, q, *app.ENV.LocalTimezone)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func GetReviewHandler(c echo.Context, app *types.App) error {
	param := c.Param("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid param")
	}

	result, _ := v1repo.GetReview(app, id)

	return c.JSON(http.StatusOK, result)
}

func GetAverageRatingHandler(c echo.Context, app *types.App) error {
	q := &v1types.ReviewQueryParams{}
	if err := c.Bind(q); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid param")
	}

	v := validator.New()
	if err := v.Struct(q); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := utils.ExtractErrorMessages(validationErrors)
		return c.JSON(http.StatusBadRequest, errorMessages)
	}

	result, err := v1repo.GetAverageRating(app.DB, q, *app.ENV.LocalTimezone)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

// Single WS connection
func ReviewWsSingleConnection(c echo.Context, app *types.App) error {
	log.Println("check in")
	ws, err := middleware.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("Error upgrade ws:", err) // TODO: use zap
		return err
	}
	defer ws.Close()

	done := make(chan bool)

	// Routine for receiving signal from Review table if got any update or create
	// and close routine when ws connection dropped.
	go func() {
		for {
			select {
			case <-connection.ReviewChannelWs:
				ws.WriteJSON(map[string]interface{}{
					"update": true,
				})
			case <-done:
				return
			}
		}
	}()

	for {
		// Write a message to the client
		err := ws.WriteMessage(websocket.TextMessage, []byte("Single, Client!"))
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Println("WebSocket connection closed")
				close(done)
				break
			}
			c.Logger().Error("Error writing to WebSocket:", err)
			close(done)
			break
		}

		// Read a message from the client
		_, msg, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Printf("WebSocket connection closed: %v \n.", err)
				close(done)
				break
			}
			c.Logger().Error("Error reading from WebSocket:", err)
			close(done)
			break
		}

		fmt.Printf("Received message: %s\n", msg)
	}

	return nil
}

// Concurrent WS connection
func ReviewWsMultipleConnection(c echo.Context, app *types.App) error {
	// Upgrade the request to WebSocket connection
	ws, err := middleware.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	// Store websocket connection as a pointer
	// because WebSocket connections are generally large objects, so need to avoid unnecessary copying by passing pointers
	connection.ActiveWsChannels.Store(ws, nil)
	log.Println("New WebSocket connection established")

	// Initialize channel for closing goroutine
	done := make(chan bool)

	// Routine for receiving signal from Review table if got any update or create
	// and close routine when ws connection dropped.
	go func() {
		for {
			select {
			case <-connection.ReviewChannelWs:
				connection.ActiveWsChannels.Range(func(key, value any) bool {
					conn := key.(*websocket.Conn)

					err := conn.WriteJSON(map[string]interface{}{
						"update": true,
					})

					if err != nil {
						log.Printf("Error sending to connection: %v", err)
					}

					return true
				})

			case <-done:
				return
			}
		}
	}()

	// Handle reading and writing from the WebSocket connection
	for {
		// Write a message to the client
		err := ws.WriteMessage(websocket.TextMessage, []byte("Multiple, Client!"))
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Println("WebSocket connection closed")
				close(done)
				break
			}
			c.Logger().Error("Error writing to WebSocket:", err)
			close(done)
			break
		}

		// Read a message from the client
		_, msg, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Printf("WebSocket connection closed: %v \n.", err)
				close(done)
				break
			}
			c.Logger().Error("Error reading from WebSocket:", err)
			close(done)
			break
		}

		fmt.Printf("Received message: %s\n", msg) // TODO: use zap or remove
	}

	// Remove the WebSocket connection from the active connections map
	connection.ActiveWsChannels.Delete(ws)

	return nil
}

func ReviewSse(c echo.Context) error {
	w := c.Response()
	w.Header().Set(echo.HeaderContentType, "text/event-stream")
	w.Header().Set(echo.HeaderCacheControl, "no-cache")
	w.Header().Set(echo.HeaderConnection, "keep-alive")

	connection.ActiveSseChannels.Store(w, nil)

	done := make(chan bool)

	go func() {
		for {
			select {
			case <-connection.ReviewChannelSse:
				connection.ActiveSseChannels.Range(func(key, value any) bool {
					sse := key.(*echo.Response)

					// Create payload
					res := map[string]interface{}{"update": true}
					parseRes, _ := json.Marshal(res)
					event := types.Event{
						Data: parseRes,
					}

					// Handle error
					if err := event.MarshalTo(sse); err != nil { // make sure marshalto sse not w (race condition)
						log.Printf("Error flush: %v \n", err)
						return false
					}

					// Flush the message
					sse.Flush()

					return true
				})

			case <-done:
				return
			}
		}
	}()

	// Handle SSE client lifecycle
	<-c.Request().Context().Done()
	log.Println("SSE client disconnected")

	// Remove active sse connection from the active connections map
	connection.ActiveSseChannels.Delete(w)

	// Close go routine
	close(done)

	return nil
}

func CheckDailyLimit(c echo.Context, rdb *redis.Client) error {
	today := time.Now().Format(time.DateOnly)
	ip := c.RealIP()
	key := database.GetRateLimitKey(ip, today)

	// Find key value
	val, _ := rdb.Get(c.Request().Context(), key).Result()
	if val != "" {
		return c.JSON(http.StatusTooManyRequests, "Rate limit exceed. Try again tomorrow.")
	}

	return c.JSON(http.StatusOK, http.StatusOK)
}
