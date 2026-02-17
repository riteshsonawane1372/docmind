package main

import (
	"context"
	"fmt"
	"os"

	"github.com/riteshsonawane/docmind/internal/config"
	"github.com/riteshsonawane/docmind/internal/rag"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: docmind <command> [args]\n\nCommands:\n  ingest <dir>   Ingest markdown files from a directory\n  chat           Start interactive RAG chat\n")
		os.Exit(1)
	}

	cfg := config.Load()
	ctx := context.Background()

	switch os.Args[1] {
	case "ingest":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Usage: docmind ingest <dir>\n")
			os.Exit(1)
		}
		if err := rag.Ingest(ctx, cfg, os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "chat":
		if err := rag.Chat(ctx, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
