package ai

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strings"
	config "student_shared/app/conf"
	"time"
)

// GetTextEmbedding returns an embedding vector for the given text.
// If OpenAI is configured, it uses the embeddings API; otherwise falls back to local hashed bag-of-words embedding.
func GetTextEmbedding(text string) ([]float64, error) {
	clean := strings.TrimSpace(text)
	if clean == "" {
		return nil, errors.New("empty text")
	}
	cfg := getOpenAIEmbeddingConfig()
	if cfg.APIKey != "" {
		if v, err := openAIEmbed(clean, cfg); err == nil && len(v) > 0 {
			return v, nil
		}
	}
	return localEmbed(clean), nil
}

type openAIEmbeddingConfig struct {
	APIKey  string
	BaseURL string
	Model   string
	Timeout time.Duration
}

func getOpenAIEmbeddingConfig() openAIEmbeddingConfig {
	c := config.Load()
	model := c.AI.Model
	return openAIEmbeddingConfig{
		APIKey: c.AI.APIKey,
		BaseURL: func() string {
			if c.AI.BaseURL != "" {
				return c.AI.BaseURL
			}
			return ""
		}(),
		Model:   model,
		Timeout: time.Duration(c.AI.TimeoutS) * time.Second,
	}
}

func openAIEmbed(text string, cfg openAIEmbeddingConfig) ([]float64, error) {
	type reqBody struct {
		Model string   `json:"model"`
		Input []string `json:"input"`
	}
	body, _ := json.Marshal(reqBody{Model: cfg.Model, Input: []string{text}})
	client := &http.Client{Timeout: cfg.Timeout}
	req, err := http.NewRequest("POST", cfg.BaseURL+"/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	var raw struct {
		Data []struct {
			Embedding []float64 `json:"embedding"`
		} `json:"data"`
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("openai embeddings error: " + resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	if len(raw.Data) == 0 {
		return nil, errors.New("no embeddings returned")
	}
	return raw.Data[0].Embedding, nil
}

// localEmbed creates a simple hashed bag-of-words embedding of fixed dimension.
func localEmbed(text string) []float64 {
	dim := 256
	vec := make([]float64, dim)
	for _, tok := range tokenize(text) {
		if tok == "" {
			continue
		}
		idx := hashToDim(tok, dim)
		vec[idx] += 1.0
	}
	// L2 normalize
	norm := 0.0
	for _, v := range vec {
		norm += v * v
	}
	norm = math.Sqrt(norm)
	if norm > 0 {
		for i := range vec {
			vec[i] /= norm
		}
	}
	return vec
}

func hashToDim(s string, dim int) int {
	h := sha1.Sum([]byte(strings.ToLower(s)))
	val := binary.BigEndian.Uint32(h[:4])
	return int(val % uint32(dim))
}

// CosineSimilarity computes cosine similarity between two vectors.
func CosineSimilarity(a, b []float64) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}
	dot := 0.0
	normA := 0.0
	normB := 0.0
	for i := 0; i < minLen; i++ {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

// SimpleHighlighter wraps occurrences of query tokens with <em> tags.
func SimpleHighlighter(text, query string) string {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" || text == "" {
		return text
	}
	toks := tokenize(q)
	res := text
	for _, t := range toks {
		if t == "" {
			continue
		}
		// case-insensitive replace by splitting
		res = highlightToken(res, t)
	}
	return res
}

func highlightToken(text, token string) string {
	lower := strings.ToLower(text)
	t := strings.ToLower(token)
	var out strings.Builder
	start := 0
	for {
		idx := strings.Index(lower[start:], t)
		if idx == -1 {
			out.WriteString(text[start:])
			break
		}
		idx += start
		out.WriteString(text[start:idx])
		out.WriteString("<em>")
		out.WriteString(text[idx : idx+len(token)])
		out.WriteString("</em>")
		start = idx + len(token)
	}
	return out.String()
}
