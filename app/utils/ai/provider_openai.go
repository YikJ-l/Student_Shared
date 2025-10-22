package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"
	config "student_shared/app/conf"
)

// openAIConfig holds runtime configuration for OpenAI provider.
type openAIConfig struct {
	APIKey  string
	BaseURL string
	Model   string
	Timeout time.Duration
}

// llmSummaryResult is the expected JSON structure returned by the LLM.
type llmSummaryResult struct {
	Summary  string   `json:"summary"`
	Keywords []string `json:"keywords"`
}

// SummarizeWithLLM tries to use an LLM provider to summarize text and extract keywords.
// Returns empty summary/keywords if provider is not configured or on error.
func SummarizeWithLLM(text string) (string, []string, error) {
	cfg := getOpenAIConfig()
	if cfg.APIKey == "" {
		return "", nil, errors.New("openai api key not configured")
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.openai.com/v1"
	}
	if cfg.Model == "" {
		cfg.Model = "gpt-3.5-turbo"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 12 * time.Second
	}

	prompt := buildPrompt(text)
	payload := map[string]any{
		"model": cfg.Model,
		"messages": []map[string]string{
			{"role": "system", "content": "You are a helpful assistant that summarizes text and extracts concise keywords. Respond ONLY with JSON of the form {\"summary\": string, \"keywords\": [string, ...]} without any commentary."},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.2,
	}
	body, _ := json.Marshal(payload)

	client := &http.Client{Timeout: cfg.Timeout}
	var raw struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	// 简单重试：最多2次
	for attempt := 0; attempt < 2; attempt++ {
		req, err := http.NewRequest("POST", cfg.BaseURL+"/chat/completions", bytes.NewReader(body))
		if err != nil { return "", nil, err }
		req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			if attempt == 1 { return "", nil, err }
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			if attempt == 1 { return "", nil, errors.New("openai api error: " + resp.Status) }
			continue
		}
		if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
			if attempt == 1 { return "", nil, err }
			continue
		}
		break
	}
	if len(raw.Choices) == 0 {
		return "", nil, errors.New("openai empty choices")
	}
	content := strings.TrimSpace(raw.Choices[0].Message.Content)

	// Try to extract JSON from content robustly
	jsonStr := extractJSON(content)
	var result llmSummaryResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		// As a fallback, attempt a naive parse: split by lines
		summary := content
		keywords := []string{}
		return summary, keywords, nil
	}

	// Normalize keywords: trim and filter empty
	kw := make([]string, 0, len(result.Keywords))
	for _, k := range result.Keywords {
		k = strings.TrimSpace(k)
		if k != "" { kw = append(kw, k) }
	}
	return strings.TrimSpace(result.Summary), kw, nil
}

func getOpenAIConfig() openAIConfig {
	c := config.Load()
	return openAIConfig{
		APIKey:  c.AI.APIKey,
		BaseURL: c.AI.BaseURL,
		Model:   c.AI.Model,
		Timeout: time.Duration(c.AI.TimeoutS) * time.Second,
	}
}

func buildPrompt(text string) string {
	// Ask the model to produce concise summary and 5-10 keywords.
	return "请为以下文本生成简洁摘要，并提取5-10个代表性关键词。用JSON格式返回：{\"summary\": string, \"keywords\": [string]}. 文本：\n\n" + text
}

func extractJSON(s string) string {
	// Find the first JSON object in the string.
	re := regexp.MustCompile(`\{[\s\S]*\}`)
	m := re.FindString(s)
	if m != "" { return m }
	return s
}