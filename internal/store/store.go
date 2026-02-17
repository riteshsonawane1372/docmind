package store

import (
	"context"
	"fmt"

	"github.com/milvus-io/milvus/client/v2/column"
	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/index"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

type Store struct {
	client     *milvusclient.Client
	collection string
	dim        int
}

type SearchResult struct {
	Content string
	Source  string
	Score   float32
}

func New(ctx context.Context, addr, collection string, dim int) (*Store, error) {
	client, err := milvusclient.New(ctx, &milvusclient.ClientConfig{
		Address: addr,
	})
	if err != nil {
		return nil, fmt.Errorf("connect to milvus: %w", err)
	}
	return &Store{
		client:     client,
		collection: collection,
		dim:        dim,
	}, nil
}

func (s *Store) EnsureCollection(ctx context.Context) error {
	has, err := s.client.HasCollection(ctx, milvusclient.NewHasCollectionOption(s.collection))
	if err != nil {
		return fmt.Errorf("check collection: %w", err)
	}
	if has {
		return nil
	}

	schema := entity.NewSchema().
		WithName(s.collection).
		WithAutoID(true).
		WithField(entity.NewField().
			WithName("id").
			WithDataType(entity.FieldTypeInt64).
			WithIsPrimaryKey(true).
			WithIsAutoID(true)).
		WithField(entity.NewField().
			WithName("content").
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(8192)).
		WithField(entity.NewField().
			WithName("source").
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(512)).
		WithField(entity.NewField().
			WithName("chunk_idx").
			WithDataType(entity.FieldTypeInt64)).
		WithField(entity.NewField().
			WithName("embedding").
			WithDataType(entity.FieldTypeFloatVector).
			WithDim(int64(s.dim)))

	if err := s.client.CreateCollection(ctx, milvusclient.NewCreateCollectionOption(s.collection, schema)); err != nil {
		return fmt.Errorf("create collection: %w", err)
	}

	idx := index.NewIvfFlatIndex(entity.COSINE, 128)
	createIdxTask, err := s.client.CreateIndex(ctx, milvusclient.NewCreateIndexOption(s.collection, "embedding", idx))
	if err != nil {
		return fmt.Errorf("create index: %w", err)
	}
	if err := createIdxTask.Await(ctx); err != nil {
		return fmt.Errorf("await index: %w", err)
	}

	loadTask, err := s.client.LoadCollection(ctx, milvusclient.NewLoadCollectionOption(s.collection))
	if err != nil {
		return fmt.Errorf("load collection: %w", err)
	}
	if err := loadTask.Await(ctx); err != nil {
		return fmt.Errorf("await load: %w", err)
	}

	return nil
}

func (s *Store) Insert(ctx context.Context, contents, sources []string, chunkIdxs []int64, embeddings [][]float32) error {
	contentCol := column.NewColumnVarChar("content", contents)
	sourceCol := column.NewColumnVarChar("source", sources)
	chunkCol := column.NewColumnInt64("chunk_idx", chunkIdxs)
	embeddingCol := column.NewColumnFloatVector("embedding", s.dim, embeddings)

	_, err := s.client.Insert(ctx, milvusclient.NewColumnBasedInsertOption(s.collection, contentCol, sourceCol, chunkCol, embeddingCol))
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}
	return nil
}

func (s *Store) Search(ctx context.Context, queryVec []float32, topK int) ([]SearchResult, error) {
	vectors := []entity.Vector{entity.FloatVector(queryVec)}
	opts := milvusclient.NewSearchOption(s.collection, topK, vectors).
		WithANNSField("embedding").
		WithOutputFields("content", "source")

	results, err := s.client.Search(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	var out []SearchResult
	for _, rs := range results {
		contentCol := rs.GetColumn("content")
		sourceCol := rs.GetColumn("source")
		for i := 0; i < rs.ResultCount; i++ {
			content, _ := contentCol.GetAsString(i)
			source, _ := sourceCol.GetAsString(i)
			out = append(out, SearchResult{
				Content: content,
				Source:  source,
				Score:   rs.Scores[i],
			})
		}
	}
	return out, nil
}

func (s *Store) Close(ctx context.Context) error {
	return s.client.Close(ctx)
}
