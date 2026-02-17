# Introduction to Retrieval Augmented Generation

Retrieval Augmented Generation (RAG) is a technique that enhances large language models by providing them with relevant context retrieved from a knowledge base. Instead of relying solely on the model's training data, RAG systems fetch relevant documents at query time and include them in the prompt.

## How RAG Works

The RAG pipeline consists of two main phases:

### Ingestion Phase

During ingestion, documents are processed and stored for later retrieval:

1. **Document Loading**: Raw documents (PDF, Markdown, HTML, etc.) are loaded into the system.
2. **Chunking**: Documents are split into smaller, semantically meaningful chunks. This is important because embedding models have token limits and smaller chunks tend to produce more focused embeddings.
3. **Embedding**: Each chunk is converted into a dense vector representation using an embedding model. These vectors capture the semantic meaning of the text.
4. **Storage**: The vectors, along with their original text and metadata, are stored in a vector database like Milvus, Pinecone, or Weaviate.

### Query Phase

When a user asks a question:

1. **Query Embedding**: The question is converted into a vector using the same embedding model.
2. **Similarity Search**: The query vector is compared against all stored vectors to find the most relevant chunks. Common similarity metrics include cosine similarity, dot product, and L2 distance.
3. **Context Assembly**: The top-K most relevant chunks are assembled into a context block.
4. **LLM Generation**: The context and question are sent to an LLM, which generates an answer grounded in the retrieved information.

## Benefits of RAG

- **Reduced Hallucination**: By grounding responses in actual documents, RAG significantly reduces the tendency of LLMs to generate plausible but incorrect information.
- **Up-to-date Information**: The knowledge base can be updated independently of the model, ensuring responses reflect the latest information.
- **Source Attribution**: RAG systems can cite their sources, enabling users to verify the information.
- **Domain Specificity**: Organizations can build RAG systems over their proprietary documents, creating specialized assistants without fine-tuning.

## Vector Databases

Vector databases are purpose-built to store and search high-dimensional vectors efficiently. Key features include:

- **Approximate Nearest Neighbor (ANN) Search**: Algorithms like IVF, HNSW, and DiskANN enable fast similarity search over millions of vectors.
- **Hybrid Search**: Many vector databases support combining vector similarity with traditional filtering.
- **Scalability**: Production vector databases can handle billions of vectors with sub-second query latency.

Milvus is an open-source vector database that supports multiple index types, GPU acceleration, and distributed deployment. It uses etcd for metadata storage and MinIO for object storage.

## Embedding Models

Embedding models convert text into dense vector representations. Popular choices include:

- **nomic-embed-text**: A high-quality open-source embedding model with 768 dimensions, available through Ollama.
- **OpenAI text-embedding-3-small**: A commercial embedding model with configurable dimensions.
- **BGE (BAAI General Embedding)**: A family of open-source embedding models in various sizes.

The choice of embedding model affects retrieval quality. Models with higher dimensions capture more nuance but require more storage and compute.

## Chunking Strategies

Effective chunking is critical for RAG performance:

- **Fixed-size chunking**: Split text every N characters or tokens. Simple but may break semantic units.
- **Paragraph-based chunking**: Split on paragraph boundaries. Preserves natural text structure.
- **Heading-aware chunking**: Use document structure (headings) to create contextually rich chunks.
- **Semantic chunking**: Use embedding similarity to determine chunk boundaries.
- **Overlap**: Including overlapping text between chunks ensures that information at chunk boundaries isn't lost.
