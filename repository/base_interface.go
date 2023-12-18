package repository

import (
	"context"

	"github.com/Alexandrhub/cli-orm-gen/genstorage/models"
	"github.com/Alexandrhub/cli-orm-gen/utils"
)

type IBaseDTO interface {
	Create(ctx context.Context, dto models.BaseDTO) error
	Upsert(ctx context.Context, dto []models.BaseDTO) error
	GetCount(ctx context.Context, dto models.BaseDTO, condition utils.Condition) (uint64, error)
	List(ctx context.Context, condition utils.Condition) ([]models.BaseDTO, error)
	Update(ctx context.Context, dto models.BaseDTO, condition utils.Condition) error
}
