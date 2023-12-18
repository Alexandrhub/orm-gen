package repository

import (
	"context"
	"fmt"

	"github.com/Alexandrhub/cli-orm-gen/db/dao"
	"github.com/Alexandrhub/cli-orm-gen/genstorage/models"
	"github.com/Alexandrhub/cli-orm-gen/infrastructure/db/scanner"
	"github.com/Alexandrhub/cli-orm-gen/utils"
)

type BaseStorage struct {
	dto *dao.DAO
}

func NewBaseStorage(dto *dao.DAO) *BaseStorage {
	return &BaseStorage{dto: dto}
}

func (b *BaseStorage) Create(ctx context.Context, dto models.BaseDTO) error {
	return b.dto.Create(ctx, &dto)
}

func (b *BaseStorage) Upsert(ctx context.Context, dto []models.BaseDTO) error {
	var entities []scanner.Tabler
	for _, d := range dto {
		entities = append(entities, &d)
	}

	return b.dto.Upsert(
		ctx,
		entities,
	)
}

func (b *BaseStorage) GetCount(ctx context.Context, dto models.BaseDTO, condition utils.Condition) (uint64, error) {
	var count uint64
	count, err := b.dto.GetCount(ctx, &dto, condition)
	if err != nil {
		return 0, fmt.Errorf("base storage: GetCount not found")
	}

	return count, nil
}

func (b *BaseStorage) List(ctx context.Context, condition utils.Condition) ([]models.BaseDTO, error) {
	var list []models.BaseDTO
	var table models.BaseDTO
	err := b.dto.List(ctx, &list, &table, condition)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (b *BaseStorage) Update(ctx context.Context, dto models.BaseDTO, condition utils.Condition) error {
	return b.dto.Update(
		ctx,
		&dto,
		condition,
		"update",
	)
}
