package config

import (
	"os"
	"strconv"
)

type Config struct {
	OllamaURL  string
	MilvusAddr string
	EmbedModel string
	ChatModel  string
	ChunkSize  int
	Overlap    int
	TopK       int
	Collection string
	EmbedDim   int
}

func Load() Config {
	c := Config{
		OllamaURL:  "http://localhost:11434",
		MilvusAddr: "localhost:19530",
		EmbedModel: "nomic-embed-text",
		ChatModel:  "llama3.2",
		ChunkSize:  512,
		Overlap:    64,
		TopK:       5,
		Collection: "documents",
		EmbedDim:   768,
	}

	if v := os.Getenv("OLLAMA_URL"); v != "" {
		c.OllamaURL = v
	}
	if v := os.Getenv("MILVUS_ADDR"); v != "" {
		c.MilvusAddr = v
	}
	if v := os.Getenv("EMBED_MODEL"); v != "" {
		c.EmbedModel = v
	}
	if v := os.Getenv("CHAT_MODEL"); v != "" {
		c.ChatModel = v
	}
	if v := os.Getenv("CHUNK_SIZE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			c.ChunkSize = n
		}
	}
	if v := os.Getenv("TOP_K"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			c.TopK = n
		}
	}

	return c
}
