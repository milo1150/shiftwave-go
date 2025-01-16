package scheduler

import (
	"encoding/json"
	"fmt"
	"log"
	"shiftwave-go/internal/types"
	v1repo "shiftwave-go/internal/v1/repository"
	v1types "shiftwave-go/internal/v1/types"
	"strconv"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const PROVIDE_LANG = "Burmese"
const TARGET_LANG = "English"

func getTranslateMessage(app *types.App) (*v1types.AssistantMessage, *v1types.UserMessage, error) {
	// Generate Assistant message
	assistantMessage := &v1types.AssistantMessage{
		Role: "system",
		Content: fmt.Sprintf(
			"You will be provided with a sentence in %s, and your task is to translate it into %s. id in schema meant to id of item. }", PROVIDE_LANG, TARGET_LANG,
		),
	}

	// Generate User message
	userMessage := &v1types.UserMessage{
		Role:    "user",
		Content: []v1types.TargetTranslateResponse{},
	}

	// Query review list
	reviews, err := v1repo.RetrieveReviewsByLang(app.DB, *app.ENV.LocalTimezone, types.LangMY, -10*time.Hour)
	if err != nil {
		return nil, nil, err
	}

	// If length = 0, no need to translate.
	if len(*reviews) == 0 {
		return nil, nil, fmt.Errorf("no need execution")
	}

	// Update user message content
	for _, review := range *reviews {
		idString := strconv.FormatUint(uint64(review.ID), 10)
		v := v1types.TargetTranslateResponse{
			idString: review.Remark,
		}
		userMessage.Content = append(userMessage.Content, v)
	}

	return assistantMessage, userMessage, nil
}

func translateAndUpdateMyanmarReviews(app *types.App) {
	// Get TranslateMessages
	assistantMessage, userMessage, err := getTranslateMessage(app)
	if err != nil {
		log.Printf("no query data needed for translate to MY")
		return
	}

	// Parse JSON to String
	assistantMessageString, _ := json.Marshal(assistantMessage)
	userMessageString, _ := json.Marshal(userMessage)

	// Log message
	log.Println("assistantMessageString: ", string(assistantMessageString))
	log.Println("userMessageString: ", string(userMessageString))

	// Initialize OpenAI Client
	client := openai.NewClient(
		option.WithAPIKey(app.ENV.OpenAI),
	)

	// Send request to OpenAI
	chatCompletion, err := client.Chat.Completions.New(app.Context, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.AssistantMessage(string(assistantMessageString)),
			openai.UserMessage(string(userMessageString)),
		}),
		Model: openai.F(openai.ChatModelGPT4oMini),
		ResponseFormat: openai.F(openai.ChatCompletionNewParamsResponseFormatUnion(openai.ChatCompletionNewParamsResponseFormat{
			Type: openai.F(openai.ChatCompletionNewParamsResponseFormatTypeJSONSchema),
			JSONSchema: openai.F(interface{}(types.AnyObject{
				"name":   "translate",
				"strict": false,
				"schema": types.AnyObject{
					"type": "object",
					"properties": types.AnyObject{
						"results": types.AnyObject{
							"type": "array",
							"items": types.AnyObject{
								"id": types.AnyObject{
									"type": "string",
								},
								"text": types.AnyObject{
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

	// Log result
	log.Printf("TotalTokens: %v, PromptTokens: %v, CompletionTokens: %v \n", chatCompletion.Usage.TotalTokens, chatCompletion.Usage.PromptTokens, chatCompletion.Usage.CompletionTokens)
	log.Println(chatCompletion.Choices[0].Message.Content)

	// Parse results to struct
	var parseResponse v1types.TranslateResponse
	json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &parseResponse)

	// Update data in Review table with the translated result
	if err := v1repo.UpdateReviewsFromTranslateResult(app.DB, parseResponse.Results); err != nil {
		log.Printf("Error update review table: %v. \n", err)
	}

	// Check func is ok
	log.Println("translateAndUpdateMyanmarReviews Done !!!")
}

func InitializeOpenAiTranslateScheduler(app *types.App) {
	// Create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Error init scheduler: %v", err)
	}

	// add a job to the scheduler
	_, err = s.NewJob(
		// TODO: change to run every hour
		// gocron.OneTimeJob(gocron.OneTimeJobStartImmediately()),
		gocron.DurationJob(
			5*time.Second,
		),
		gocron.NewTask(
			func() {
				translateAndUpdateMyanmarReviews(app)
			},
		),
	)
	if err != nil {
		log.Fatalf("Error init NewJob: %v", err)
	}

	// Start Cron
	s.Start()
}
