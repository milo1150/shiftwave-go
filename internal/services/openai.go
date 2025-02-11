package services

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"shiftwave-go/internal/enum"
	"shiftwave-go/internal/types"
	v1repo "shiftwave-go/internal/v1/repository"
	v1types "shiftwave-go/internal/v1/types"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	openaiV2 "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
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
	reviews, err := v1repo.RetrieveReviewsByLang(app.DB, *app.ENV.LocalTimezone, enum.LangMY, -10*time.Hour)
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

func getOpenAIClientV2(authToken string) *openaiV2.Client {
	// Load system certificates
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		log.Printf("Failed to load system certificates: %v", err)
		caCertPool = x509.NewCertPool()
	}

	// Create custom TLS configuration
	tlsConfig := &tls.Config{
		RootCAs: caCertPool, // Use system CA pool for certificate validation
	}

	// Create custom HTTP transport
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	// Create HTTP client
	httpClient := &http.Client{
		Transport: transport,
	}

	defaultConfig := openaiV2.DefaultConfig(authToken)
	defaultConfig.HTTPClient = httpClient

	// Initialize OpenAI client with the custom HTTP client
	client := openaiV2.NewClientWithConfig(defaultConfig)

	return client
}

func TranslateAndUpdateMyanmarReviewsV2(app *types.App) {
	// Get TranslateMessages
	// If query length = 0, then do nothing.
	assistantMessage, userMessage, err := getTranslateMessage(app)
	if err != nil {
		log.Printf("no query data needed for translate to MY")
		return
	}

	// Parse JSON to String
	assistantMessageString, err := json.Marshal(assistantMessage)
	if err != nil {
		log.Printf("Error marshal assistantMessage")
		return
	}
	userMessageString, err := json.Marshal(userMessage)
	if err != nil {
		log.Printf("Error marshal userMessage")
		return
	}

	// Log message
	log.Println("assistantMessageString: ", string(assistantMessageString))
	log.Println("userMessageString: ", string(userMessageString))

	// Initialize OpenAI Client
	client := getOpenAIClientV2(app.ENV.OpenAI)
	// client := openaiV2.NewClient(app.ENV.OpenAI)

	// Define the expected result structure
	type Result struct {
		Results []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"results"`
	}
	var result Result

	// Generate JSON schema for the expected result structure
	schema, err := jsonschema.GenerateSchemaForType(result)
	if err != nil {
		log.Fatalf("GenerateSchemaForType error: %v", err)
	}

	// Send request to OpenAI
	chatCompletion, err := client.CreateChatCompletion(
		app.Context,
		openaiV2.ChatCompletionRequest{
			Model: openaiV2.GPT4oMini,
			Messages: []openaiV2.ChatCompletionMessage{
				{
					Role:    openaiV2.ChatMessageRoleUser,
					Content: string(userMessageString),
				},
				{
					Role:    openaiV2.ChatMessageRoleAssistant,
					Content: string(assistantMessageString),
				},
			},
			ResponseFormat: &openaiV2.ChatCompletionResponseFormat{
				Type: openaiV2.ChatCompletionResponseFormatTypeJSONSchema,
				JSONSchema: &openaiV2.ChatCompletionResponseFormatJSONSchema{
					Name:   "translate_burmese_to_english",
					Strict: true,
					Schema: schema,
				},
			},
		},
	)

	// Prevent nil pointer dereference
	if err != nil {
		log.Printf("Error from OpenAI: %v \n", err.Error())
		return
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
		return
	}

	// Check func is ok
	log.Println("translateAndUpdateMyanmarReviewsV2 Done !!!")
}

func getOpenAIClient(apiKey string) *openai.Client {
	// Load system certificates
	caCertPool, err := x509.SystemCertPool()
	if err != nil || caCertPool == nil {
		caCertPool = x509.NewCertPool()
	}

	// Create custom TLS configuration
	tlsConfig := &tls.Config{
		RootCAs:            caCertPool,
		InsecureSkipVerify: false, // Enable certificate verification
	}

	// Create custom HTTP transport
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	// Create HTTP client
	httpClient := &http.Client{
		Transport: transport,
	}

	// Initialize OpenAI client with the custom HTTP client
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithHTTPClient(httpClient),
	)

	return client
}

// NOTE
// Local - Work fine
// Prod - do not use debian:bullseye-slim otherwise will face this error
// Error from OpenAI: Post "https://api.openai.com/v1/chat/completions": tls: failed to verify certificate: x509: certificate signed by unknown authority
func TranslateAndUpdateMyanmarReviews(app *types.App) {
	// Get TranslateMessages
	// If query length = 0, then do nothing.
	assistantMessage, userMessage, err := getTranslateMessage(app)
	if err != nil {
		log.Printf("no query data needed for translate to MY")
		return
	}

	// Parse JSON to String
	assistantMessageString, err := json.Marshal(assistantMessage)
	if err != nil {
		log.Printf("Error marshal assistantMessage")
		return
	}
	userMessageString, err := json.Marshal(userMessage)
	if err != nil {
		log.Printf("Error marshal userMessage")
		return
	}

	// Log message
	log.Println("assistantMessageString: ", string(assistantMessageString))
	log.Println("userMessageString: ", string(userMessageString))

	// Initialize OpenAI Client
	client := getOpenAIClient(app.ENV.OpenAI)
	// client := openai.NewClient(
	// 	option.WithAPIKey(app.ENV.OpenAI),
	// )

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

	// Prevent nil pointer dereference
	if err != nil {
		log.Printf("Error from OpenAI: %v \n", err.Error())
		return
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
		return
	}

	// Check func is ok
	log.Println("translateAndUpdateMyanmarReviews Done !!!")
}
