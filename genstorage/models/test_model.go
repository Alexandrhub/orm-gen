package models

import (
	"time"
	
	"github.com/Alexandrhub/cli-orm-gen/infrastructure/db/types"
)

type TestDTO struct {
	ID        int            `json:"id" db:"id" db_type:"BIGSERIAL primary key" db_default:"not null" db_ops:"id" mapper:"id"`
	UUID      string         `json:"uuid" db:"uuid" db_ops:"create" db_type:"char(36)" db_default:"not null" db_index:"index,unique" mapper:"uuid"`
	Active    types.NullBool `json:"active" db:"active" db_ops:"create,update" db_type:"boolean" db_default:"null" mapper:"active"`
	CreatedAt time.Time      `json:"created_at" db:"created_at" db_type:"timestamp" db_default:"default (now()) not null" db_index:"index" db_ops:"created_at" mapper:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at" db_ops:"update" db_type:"timestamp" db_default:"default (now()) not null" db_index:"index" mapper:"updated_at"`
	DeletedAt types.NullTime `json:"deleted_at" db:"deleted_at" db_ops:"update" db_type:"timestamp" db_default:"default null" db_index:"index" db_ops:"deleted_at" mapper:"deleted_at"`
}

func (b *TestDTO) TableName() string {
	return "my_test"
}

func (b *TestDTO) OnCreate() []string {
	return []string{}
}

func (b *TestDTO) SetID(id int) *TestDTO {
	b.ID = id
	return b
}

func (b *TestDTO) GetID() int {
	return b.ID
}

func (b *TestDTO) SetUUID(uuid string) *TestDTO {
	b.UUID = uuid
	return b
}

func (b *TestDTO) GetUUID() string {
	return b.UUID
}

func (b *TestDTO) SetCreatedAt(createdAt time.Time) *TestDTO {
	b.CreatedAt = createdAt
	return b
}

func (b *TestDTO) GetCreatedAt() time.Time {
	return b.CreatedAt
}

func (b *TestDTO) SetUpdatedAt(updatedAt time.Time) *TestDTO {
	b.UpdatedAt = updatedAt
	return b
}

func (b *TestDTO) GetUpdatedAt() time.Time {
	return b.UpdatedAt
}

func (b *TestDTO) SetDeletedAt(deletedAt time.Time) *TestDTO {
	b.DeletedAt.Time.Time = deletedAt
	b.DeletedAt.Time.Valid = true
	return b
}

func (b *TestDTO) GetDeletedAt() time.Time {
	return b.DeletedAt.Time.Time
}

func (b *TestDTO) SetActive(active bool) *TestDTO {
	b.Active.Bool = active
	b.Active.Valid = true
	
	return b
}

func (b *TestDTO) GetActive() bool {
	return b.Active.Bool
}
