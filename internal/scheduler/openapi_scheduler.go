package scheduler

import (
	"encoding/json"
	"fmt"
	"log"
	"shiftwave-go/internal/types"
	v1repo "shiftwave-go/internal/v1/repository"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// TODO: return dto
func MessageBuilder(app *types.App) {
	reviews, err := v1repo.RetrieveReviewsByLang(app.DB, types.LangMY, *app.ENV.LocalTimezone, time.Hour)
	if err != nil {
	}
	fmt.Println(len(*reviews))
}

func start(app *types.App) {
	fmt.Println("Fetch OpenAI Translate")
	// return
	// a := map[string]string{"role": "system", "content": "You will be provided with a sentence in Thailand, and your task is to translate it into English."}
	// b := `{"role": "system", "content": "You will be provided with a sentence in Thailand, and your task is to translate it into English."}`
	// c := `{"role": "user", "content": { "point":"Respond with a valid JSON object.", "values":[{"1","วันนี้กินอะไรดี"},{"2":"พรุ่งนี้รถจะติดไหม"}]}`
	d := `[
		{"role": "system", "content": "You will be provided with a sentence in Thailand, and your task is to translate it into English."},
		{"role": "user", "content": { "point":"Respond with a valid JSON object.", "values":[{"1","วันนี้กินอะไรดี"},{"2":"พรุ่งนี้รถจะติดไหม"}]}
		]`

	client := openai.NewClient(
		option.WithAPIKey(app.ENV.OpenAI),
	)
	chatCompletion, err := client.Chat.Completions.New(app.Context, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			// openai.AssistantMessage(b),
			openai.UserMessage(d),
		}),
		Model: openai.F(openai.ChatModelGPT4oMini),
		ResponseFormat: openai.F(openai.ChatCompletionNewParamsResponseFormatUnion(openai.ChatCompletionNewParamsResponseFormat{
			Type: openai.F(openai.ChatCompletionNewParamsResponseFormatTypeJSONSchema),
			JSONSchema: openai.F(interface{}(map[string]interface{}{
				"name":   "translate",
				"strict": false,
				"schema": map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{
						// "role":    map[string]interface{}{"type": "string"},
						// "content": map[string]interface{}{"type": "string"},
					},
					// "required": []interface{}{"role", "content"},
				},
			})),
		})),
	})
	if err != nil {
		panic(err.Error())
	}

	var parseResponse map[string]interface{}
	json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &parseResponse)

	v := fmt.Sprintf("%v", parseResponse)

	fmt.Println(v)
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
			1*time.Second,
		),
		gocron.NewTask(
			func() {
				count++
				log.Println(count)

				if count == 1 {
					s.Shutdown()
					s.StopJobs()
				}

				start(app)

			},
		),
	)
	if err != nil {
		log.Fatalf("Error init NewJob: %v", err)
	}

	// Start Cron
	s.Start()
}
