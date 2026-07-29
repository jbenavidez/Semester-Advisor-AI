package weaviate

import (
	"semester-advisor-ai/internal/ports"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

type WeaviateDBRepo struct {
	DB *weaviate.Client
}

func NewWeaviateDBRepo(db *weaviate.Client) ports.UploadedFileRepository {
	return &WeaviateDBRepo{
		DB: db,
	}

}

func (m *WeaviateDBRepo) SaveFile() {
	// to be continue
}
