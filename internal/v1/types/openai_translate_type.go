package types

type AssistantMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type UserMessage struct {
	Role    string                    `json:"role"`
	Content []TargetTranslateResponse `json:"content"`
}

type TargetTranslateResponse map[string]string

type TranslateResponse struct {
	Results []TranslateResult `json:"results"`
}

type TranslateResult struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}
