package repository

import (
	"context"

	"github.com/Alexandrhub/cli-orm-gen/genstorage/models"
	"github.com/Alexandrhub/cli-orm-gen/utils"
)

type ITestDTO interface {
	Create(ctx context.Context, dto models.TestDTO) error
	Upsert(ctx context.Context, dto []models.TestDTO) error
	GetCount(ctx context.Context, dto models.TestDTO, condition utils.Condition) (uint64, error)
	List(ctx context.Context, condition utils.Condition) ([]models.TestDTO, error)
	Update(ctx context.Context, dto models.TestDTO, condition utils.Condition) error
}
