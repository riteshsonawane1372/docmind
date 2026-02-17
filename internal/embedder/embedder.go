package embedder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Embedder struct {
	baseURL string
	model   string
	client  *http.Client
}

type embedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type embedResponse struct {
	Embedding []float32 `json:"embedding"`
}

func New(baseURL, model string) *Embedder {
	return &Embedder{
		baseURL: baseURL,
		model:   model,
		client:  &http.Client{},
	}
}

// Embed sends a batch of texts to Ollama /api/embeddings and returns embeddings.
// Note: Ollama's API processes one text at a time, so we loop through the batch.
func (e *Embedder) Embed(texts []string) ([][]float32, error) {
	embeddings := make([][]float32, len(texts))

	for i, text := range texts {
		body, err := json.Marshal(embedRequest{
			Model:  e.model,
			Prompt: text,
		})
		if err != nil {
			return nil, fmt.Errorf("marshal embed request: %w", err)
		}

		resp, err := e.client.Post(e.baseURL+"/api/embeddings", "application/json", bytes.NewReader(body))
		if err != nil {
			return nil, fmt.Errorf("embed request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("embed request returned status %d", resp.StatusCode)
		}

		var result embedResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("decode embed response: %w", err)
		}

		embeddings[i] = result.Embedding
	}

	return embeddings, nil
}

// EmbedSingle embeds a single text string.
func (e *Embedder) EmbedSingle(text string) ([]float32, error) {
	embeddings, err := e.Embed([]string{text})
	if err != nil {
		return nil, err
	}
	if len(embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}
	return embeddings[0], nil
}
