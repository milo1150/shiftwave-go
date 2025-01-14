package scheduler

import (
	"encoding/json"
	"fmt"
	"log"
	"shiftwave-go/internal/types"
	v1repo "shiftwave-go/internal/v1/repository"
	"strconv"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type assistantMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type userMessage struct {
	Role    string              `json:"role"`
	Content []translateResponse `json:"content"`
}

type translateResponse map[string]string

const PROVIDE_LANG = "Burmese"
const TARGET_LANG = "English"

func getTranslateMessage(app *types.App) (*assistantMessage, *userMessage, error) {

	var assistantMessage = &assistantMessage{
		Role: "system",
		Content: fmt.Sprintf(
			"You will be provided with a sentence in %s, and your task is to translate it into %s. id in schema meant to id of item. }", PROVIDE_LANG, TARGET_LANG,
		),
	}

	var userMessage = &userMessage{
		Role:    "user",
		Content: []translateResponse{},
	}

	reviews, err := v1repo.RetrieveReviewsByLang(app.DB, types.LangMY, *app.ENV.LocalTimezone, time.Hour)
	if err != nil {
		return nil, nil, err
	}

	for _, review := range *reviews {
		idString := strconv.FormatUint(uint64(review.ID), 10)
		fmt.Println("ID:", idString)
		fmt.Println("Remark:", review.Remark)
		v := translateResponse{
			idString: review.Remark,
		}
		userMessage.Content = append(userMessage.Content, v)
	}

	fmt.Println(len(*reviews))

	return assistantMessage, userMessage, nil
}

func start(app *types.App) {
	assistantMessage, userMessage, err := getTranslateMessage(app)
	if err != nil {
		fmt.Printf("error query data for translate: %v", err.Error())
	}

	assistantMessageString, _ := json.Marshal(assistantMessage)
	userMessageString, _ := json.Marshal(userMessage)

	log.Println("Fetch OpenAI Translate - assistantMessageString", string(assistantMessageString))
	log.Println("Fetch OpenAI Translate - assistantMessageString", string(userMessageString))

	// return

	client := openai.NewClient(
		option.WithAPIKey(app.ENV.OpenAI),
	)
	chatCompletion, err := client.Chat.Completions.New(app.Context, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.AssistantMessage(string(assistantMessageString)),
			openai.UserMessage(string(userMessageString)),
		}),
		Model: openai.F(openai.ChatModelGPT4oMini),
		ResponseFormat: openai.F(openai.ChatCompletionNewParamsResponseFormatUnion(openai.ChatCompletionNewParamsResponseFormat{
			Type: openai.F(openai.ChatCompletionNewParamsResponseFormatTypeJSONSchema),
			JSONSchema: openai.F(interface{}(map[string]interface{}{
				"name":   "translate",
				"strict": false,
				"schema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"results": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"id": map[string]interface{}{
									"type": "string",
								},
								"text": map[string]interface{}{
									"type": "string",
								},
							},
						},
					},
					"required": []string{"results"},
				},
			})),
		})),
	})
	if err != nil {
		log.Printf("Error from OpenAI: %v \n", err.Error())
	}

	fmt.Println(chatCompletion.Choices[0].Message.Content)

	// var parseResponse map[string]interface{}
	// json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &parseResponse)
}

func OpenAITranslateScheduler(app *types.App) {
	// Create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Error init scheduler: %v", err)
	}

	count := 0

	// add a job to the scheduler
	_, err = s.NewJob(
		gocron.DurationJob(
			5*time.Second,
		),
		gocron.NewTask(
			func() {
				log.Println(count)
				if count == 1 {
					s.Shutdown()
				}

				start(app)

				count++
			},
		),
	)
	if err != nil {
		log.Fatalf("Error init NewJob: %v", err)
	}

	// Start Cron
	s.Start()
}
