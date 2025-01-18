package handler

import (
	"fmt"
	"log"
	"net/http"
	"shiftwave-go/internal/middleware"
	"shiftwave-go/internal/types"
	"shiftwave-go/internal/utils"
	v1repository "shiftwave-go/internal/v1/repository"
	v1types "shiftwave-go/internal/v1/types"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CreateReviewHandler(c echo.Context, db *gorm.DB) error {
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

	if result := v1repository.CreateReview(db, payload); result != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error()})
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

	result, err := v1repository.GetReviews(app, q, *app.ENV.LocalTimezone)
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

	result, _ := v1repository.GetReview(app, id)

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

	result, err := v1repository.GetAverageRating(app.DB, q, *app.ENV.LocalTimezone)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

// Single WS connection
func ReviewWsSingleConnection(c echo.Context, app *types.App) error {
	ws, err := middleware.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	done := make(chan bool)

	// Routine for receiving signal from Review table if got any update or create
	// and close routine when ws connection dropped.
	go func() {
		for {
			select {
			case <-middleware.ReviewChannel:
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
	middleware.ActiveChannels.Store(ws, nil)
	log.Println("New WebSocket connection established")

	// Initialize channel for closing goroutine
	done := make(chan bool)

	// Routine for receiving signal from Review table if got any update or create
	// and close routine when ws connection dropped.
	go func() {
		for {
			select {
			case <-middleware.ReviewChannel:
				middleware.ActiveChannels.Range(func(key, value any) bool {
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

		fmt.Printf("Received message: %s\n", msg)
	}

	// Remove the WebSocket connection from the active connections map
	middleware.ActiveChannels.Delete(ws)

	return nil
}
