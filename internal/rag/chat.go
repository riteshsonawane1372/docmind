package rag

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/riteshsonawane/docmind/internal/config"
	"github.com/riteshsonawane/docmind/internal/embedder"
	"github.com/riteshsonawane/docmind/internal/llm"
	"github.com/riteshsonawane/docmind/internal/store"
)

const maxHistory = 5

const systemPrompt = `
You are a CLI-based documentation assistant.

You answer questions strictly using ONLY the provided context.
The context comes from retrieved markdown files.

Rules:
1. Do NOT use prior knowledge.
2. Do NOT guess or infer beyond the provided context.
3. If the answer is not clearly present, say:
   "The provided context does not contain enough information to answer this."
4. Be concise but complete.
5. If possible, cite the relevant section or file name from the context.
6. If the question is ambiguous, ask for clarification.

Answer format:
- Direct answer first.
- Then optional short supporting explanation from context.
- Include citations if available.
`

func Chat(ctx context.Context, cfg config.Config) error {
	emb := embedder.New(cfg.OllamaURL, cfg.EmbedModel)
	llmClient := llm.New(cfg.OllamaURL, cfg.ChatModel)

	st, err := store.New(ctx, cfg.MilvusAddr, cfg.Collection, cfg.EmbedDim)
	if err != nil {
		return fmt.Errorf("connect to store: %w", err)
	}
	defer st.Close(ctx)

	var history []llm.Message

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("DocMind RAG Chat (type 'quit' to exit)")
	fmt.Println(strings.Repeat("-", 40))

	for {
		fmt.Print("\nYou: ")
		if !scanner.Scan() {
			break
		}
		question := strings.TrimSpace(scanner.Text())
		if question == "" {
			continue
		}
		if question == "quit" || question == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		fmt.Print("🔍 Embedding query... ")
		queryVec, err := emb.EmbedSingle(question)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError embedding question: %v\n", err)
			continue
		}
		fmt.Println("✓")

		fmt.Print("📚 Searching knowledge base... ")
		results, err := st.Search(ctx, queryVec, cfg.TopK)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError searching: %v\n", err)
			continue
		}
		fmt.Printf("✓ (found %d results)\n", len(results))

		contextBlock := buildContext(results)

		messages := []llm.Message{
			{Role: "system", Content: systemPrompt},
		}
		// Add conversation history (last N exchanges).
		if len(history) > maxHistory*2 {
			history = history[len(history)-maxHistory*2:]
		}
		messages = append(messages, history...)

		userMsg := fmt.Sprintf("Context:\n%s\n\nQuestion: %s", contextBlock, question)
		messages = append(messages, llm.Message{Role: "user", Content: userMsg})

		fmt.Print("🤖 Generating response...\n\nA: ")
		response, err := llmClient.ChatStream(messages, func(token string) {
			fmt.Print(token)
		})
		fmt.Println()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error from LLM: %v\n", err)
			continue
		}

		history = append(history,
			llm.Message{Role: "user", Content: question},
			llm.Message{Role: "assistant", Content: response},
		)
	}

	return nil
}

func buildContext(results []store.SearchResult) string {
	var b strings.Builder
	for i, r := range results {
		fmt.Fprintf(&b, "[%d] (source: %s, score: %.4f)\n%s\n\n", i+1, r.Source, r.Score, r.Content)
	}
	return b.String()
}
