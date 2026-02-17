# DocMind - Development Context

> **📝 IMPORTANT NOTE FOR DEVELOPMENT:**
> This file tracks the development history and context of the project. Whenever changes are made during a session, they MUST be documented here under the appropriate session date. This ensures continuity across sessions and helps both developers and AI assistants understand what has been done and why.
>
> **When updating this file:**
> - Add new session entries with the date
> - Document what was changed and why
> - Include relevant file paths
> - Note any issues fixed or features added
> - Update TODO items as needed

## Project Overview
DocMind is a Go-based RAG (Retrieval-Augmented Generation) system that ingests markdown documents, creates embeddings using Ollama, stores them in Milvus vector database, and enables chat-based querying.

## Architecture
- **Embedder**: Uses Ollama API to generate embeddings from text
- **Chunker**: Splits markdown documents into manageable chunks
- **Store**: Manages Milvus vector database operations
- **RAG**: Orchestrates ingestion and chat functionality
- **LLM**: Handles chat completions via Ollama

## Recent Changes

### Session: 2026-02-17

#### 1. Build System Setup
**Files Created:**
- `Makefile` - Build automation with targets:
  - `make build` - Compiles the binary
  - `make install` - Installs binary to `~/bin`
  - `make clean` - Removes build artifacts
  - `make test` - Runs tests
  - `make fmt` - Formats code
  - `make vet` - Runs static analysis
  - `make check` - Runs fmt, vet, and test
  - `make uninstall` - Removes installed binary

- `.gitignore` - Standard Go project ignores including:
  - Binary outputs
  - Test coverage files
  - IDE files
  - Environment files
  - Log files

#### 2. Fixed Ollama Embeddings API Integration
**Problem:**
- Getting 404 error when calling Ollama embeddings endpoint
- Error occurred during document ingestion: "embed request returned status 404"

**Root Cause:**
- Incorrect API endpoint: was using `/api/embed` instead of `/api/embeddings`
- Wrong request format: was sending batch `Input []string` field instead of single `Prompt string`
- Wrong response format: was expecting `Embeddings [][]float32` instead of `Embedding []float32`

**Solution Applied:**
Modified `internal/embedder/embedder.go`:
- Changed endpoint from `/api/embed` to `/api/embeddings`
- Updated request structure to match Ollama API:
  ```go
  type embedRequest struct {
      Model  string `json:"model"`
      Prompt string `json:"prompt"`  // Changed from Input []string
  }
  ```
- Updated response structure:
  ```go
  type embedResponse struct {
      Embedding []float32 `json:"embedding"`  // Changed from Embeddings [][]float32
  }
  ```
- Modified `Embed()` function to loop through texts one at a time (Ollama processes embeddings individually)

**Files Modified:**
- `internal/embedder/embedder.go` - Fixed API endpoint and request/response structures

#### 3. Added Progress Indicators for Chat
**Problem:**
- No visual feedback during chat interactions
- Users couldn't tell if the system was processing their query or if something was stuck
- No indication of what stage the process was in (embedding, searching, generating)

**Solution Applied:**
Added progress logging to `internal/rag/chat.go`:
- 🔍 "Embedding query..." - Shows when query is being embedded
- 📚 "Searching knowledge base..." - Shows when searching Milvus
- Shows number of results found
- 🤖 "Generating response..." - Shows when LLM is being called
- Streams tokens in real-time as they're generated

Added error logging to `internal/llm/llm.go`:
- Captures and displays Ollama error messages for debugging
- Shows HTTP status codes when requests fail

**Files Modified:**
- `internal/rag/chat.go` - Added progress indicators and status messages
- `internal/llm/llm.go` - Added error response logging and required imports

#### 4. Installed Required Ollama Model
**Problem:**
- `nomic-embed-text` model was not installed, causing 404 errors

**Solution:**
- Ran `ollama pull nomic-embed-text` to download the model
- Verified model outputs 768-dimensional embeddings (matches config)

## Current Configuration

### Default Settings (config.go)
- **Ollama URL**: `http://localhost:11434`
- **Milvus Address**: `localhost:19530`
- **Embed Model**: `nomic-embed-text`
- **Chat Model**: `llama3.2`
- **Chunk Size**: 512 tokens
- **Overlap**: 64 tokens
- **Top K**: 5 results
- **Collection**: "documents"
- **Embed Dimension**: 768

### Environment Variables
Can override defaults with:
- `OLLAMA_URL`
- `MILVUS_ADDR`
- `EMBED_MODEL`
- `CHAT_MODEL`
- `CHUNK_SIZE`
- `TOP_K`

## Docker Services
Running via `docker-compose.yml`:
- **etcd** - Distributed key-value store for Milvus
- **MinIO** - Object storage for Milvus
- **Milvus** - Vector database (port 19530)

## Next Steps / TODO
- [ ] Test the fixed embedding integration with `./docmind ingest <directory>`
- [ ] Verify embeddings are correctly stored in Milvus
- [ ] Test chat functionality
- [ ] Consider adding progress bars for large ingestions
- [ ] Add error handling for Ollama connectivity issues
- [ ] Document required Ollama models setup

## Known Issues
- None currently (previous 404 error resolved)

## Development Commands

### Build & Install
```bash
make build          # Build binary
make install        # Install to ~/bin
make clean          # Clean build artifacts
```

### Running the Application
```bash
./docmind ingest <directory>    # Ingest markdown files
./docmind chat                   # Start chat interface
```

### Docker
```bash
docker-compose up -d    # Start services
docker-compose down     # Stop services
```

### Ollama Setup
Ensure the embedding model is pulled:
```bash
ollama pull nomic-embed-text
ollama pull llama3.2
```

## Project Structure
```
docmind/
├── internal/
│   ├── chunker/      # Document chunking logic
│   ├── config/       # Configuration management
│   ├── embedder/     # Ollama embedding integration
│   ├── llm/          # LLM chat integration
│   ├── rag/          # RAG orchestration (ingest, chat)
│   └── store/        # Milvus vector store
├── docs/             # Documentation to ingest
├── main.go           # Entry point
├── Makefile          # Build automation
├── docker-compose.yml # Docker services
└── CONTEXT.md        # This file
```
