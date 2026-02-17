# DocMind

A fully local RAG (Retrieval Augmented Generation) CLI application built in Go. Ingest markdown files, create embeddings, store them in Milvus, and chat with your documents using Ollama.

## Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)
- [Ollama](https://ollama.ai/)

## Setup

### 1. Start Milvus

```bash
docker compose up -d
```

Wait for Milvus to be healthy:

```bash
docker compose ps
```

### 2. Pull Ollama Models

```bash
ollama pull nomic-embed-text
ollama pull llama3.2
```

### 3. Build

```bash
go build -o docmind .
```

## Usage

### Ingest Documents

```bash
./docmind ingest ./docs
```

This will:
- Scan the directory for `.md` files
- Split them into heading-aware chunks
- Generate embeddings via Ollama
- Store everything in Milvus

### Chat

```bash
./docmind chat
```

Type your questions and get answers grounded in your ingested documents. Type `quit` to exit.

## Configuration

All settings can be overridden with environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `OLLAMA_URL` | `http://localhost:11434` | Ollama API URL |
| `MILVUS_ADDR` | `localhost:19530` | Milvus gRPC address |
| `EMBED_MODEL` | `nomic-embed-text` | Ollama embedding model |
| `CHAT_MODEL` | `llama3.2` | Ollama chat model |
| `CHUNK_SIZE` | `512` | Max chunk size in characters |
| `TOP_K` | `5` | Number of results to retrieve |

## Architecture

```
Ingest: .md files -> Chunker -> Ollama Embedder -> Milvus
Chat:   Question -> Embed -> Milvus Search -> Build Prompt -> Ollama LLM -> Stream Response
```
