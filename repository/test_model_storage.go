package repository

import (
	"context"
	"fmt"

	"github.com/Alexandrhub/cli-orm-gen/db/dao"
	"github.com/Alexandrhub/cli-orm-gen/genstorage/models"
	"github.com/Alexandrhub/cli-orm-gen/infrastructure/db/scanner"
	"github.com/Alexandrhub/cli-orm-gen/utils"
)

type TestDTOStorage struct {
	dto *dao.DAO
}

func NewTestDTOStorage(dto *dao.DAO) *TestDTOStorage {
	return &TestDTOStorage{dto: dto}
}

func (t *TestDTOStorage) Create(ctx context.Context, dto models.TestDTO) error {
	return t.dto.Create(ctx, &dto)
}

func (t *TestDTOStorage) Upsert(ctx context.Context, dto []models.TestDTO) error {
	var entities []scanner.Tabler
	for _, d := range dto {
		entities = append(entities, &d)
	}

	return t.dto.Upsert(
		ctx,
		entities,
	)
}

func (t *TestDTOStorage) GetCount(ctx context.Context, dto models.TestDTO, condition utils.Condition) (uint64, error) {
	var count uint64
	count, err := t.dto.GetCount(ctx, &dto, condition)
	if err != nil {
		return 0, fmt.Errorf("testdto storage: GetCount not found")
	}

	return count, nil
}

func (t *TestDTOStorage) List(ctx context.Context, condition utils.Condition) ([]models.TestDTO, error) {
	var list []models.TestDTO
	var table models.TestDTO
	err := t.dto.List(ctx, &list, &table, condition)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (t *TestDTOStorage) Update(ctx context.Context, dto models.TestDTO, condition utils.Condition) error {
	return t.dto.Update(
		ctx,
		&dto,
		condition,
		"update",
	)
}
