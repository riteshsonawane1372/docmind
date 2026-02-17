package rag

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/riteshsonawane/docmind/internal/chunker"
	"github.com/riteshsonawane/docmind/internal/config"
	"github.com/riteshsonawane/docmind/internal/embedder"
	"github.com/riteshsonawane/docmind/internal/store"
)

const batchSize = 32

func Ingest(ctx context.Context, cfg config.Config, dir string) error {
	emb := embedder.New(cfg.OllamaURL, cfg.EmbedModel)

	st, err := store.New(ctx, cfg.MilvusAddr, cfg.Collection, cfg.EmbedDim)
	if err != nil {
		return fmt.Errorf("connect to store: %w", err)
	}
	defer st.Close(ctx)

	if err := st.EnsureCollection(ctx); err != nil {
		return fmt.Errorf("ensure collection: %w", err)
	}

	var allChunks []chunker.Chunk

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		relPath, _ := filepath.Rel(dir, path)
		chunks := chunker.ChunkMarkdown(string(data), relPath, cfg.ChunkSize, cfg.Overlap)
		fmt.Printf("  %s: %d chunks\n", relPath, len(chunks))
		allChunks = append(allChunks, chunks...)
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk dir: %w", err)
	}

	if len(allChunks) == 0 {
		fmt.Println("No markdown files found.")
		return nil
	}

	fmt.Printf("\nTotal chunks: %d\nEmbedding and inserting...\n", len(allChunks))

	for i := 0; i < len(allChunks); i += batchSize {
		end := i + batchSize
		if end > len(allChunks) {
			end = len(allChunks)
		}
		batch := allChunks[i:end]

		texts := make([]string, len(batch))
		sources := make([]string, len(batch))
		idxs := make([]int64, len(batch))
		for j, c := range batch {
			texts[j] = c.Content
			sources[j] = c.Source
			idxs[j] = c.ChunkIdx
		}

		embeddings, err := emb.Embed(texts)
		if err != nil {
			return fmt.Errorf("embed batch %d: %w", i/batchSize, err)
		}

		if err := st.Insert(ctx, texts, sources, idxs, embeddings); err != nil {
			return fmt.Errorf("insert batch %d: %w", i/batchSize, err)
		}

		fmt.Printf("  Inserted %d/%d chunks\n", end, len(allChunks))
	}

	fmt.Println("Ingestion complete.")
	return nil
}
