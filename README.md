# cli-orm-gen



## Что это такое?

Генерация файлов хранилища и интерфейса CRUD реализованная для заданной модели.

## Пример вызова

```bash
go build cmd/cli-orm/repogen.go
./repogen -entity="./genstorage/models/test_model.go"
```

Все нужные файлы будут сгенерированы в корневую папку ./storage вашего проекта

## Пример модели для генерации

```go
package models


type BaseDTO struct {
	ID        int            `json:"id" db:"id" db_type:"BIGSERIAL primary key" db_default:"not null" db_ops:"id" mapper:"id"`
	UUID      string         `json:"uuid" db:"uuid" db_ops:"create" db_type:"char(36)" db_default:"not null" db_index:"index,unique" mapper:"uuid"`
	Active    types.NullBool `json:"active" db:"active" db_ops:"create,update" db_type:"boolean" db_default:"null" mapper:"active"`
	CreatedAt time.Time      `json:"created_at" db:"created_at" db_type:"timestamp" db_default:"default (now()) not null" db_index:"index" db_ops:"created_at" mapper:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at" db_ops:"update" db_type:"timestamp" db_default:"default (now()) not null" db_index:"index" mapper:"updated_at"`
	DeletedAt types.NullTime `json:"deleted_at" db:"deleted_at" db_ops:"update" db_type:"timestamp" db_default:"default null" db_index:"index" db_ops:"deleted_at" mapper:"deleted_at"`
}

func (b *BaseDTO) TableName() string {
	return "base"
}

func (b *BaseDTO) OnCreate() []string {
	return []string{}
}

func (b *BaseDTO) SetID(id int) *BaseDTO {
	b.ID = id
	return b
}

func (b *BaseDTO) GetID() int {
	return b.ID
}

func (b *BaseDTO) SetUUID(uuid string) *BaseDTO {
	b.UUID = uuid
	return b
}

func (b *BaseDTO) GetUUID() string {
	return b.UUID
}

func (b *BaseDTO) SetCreatedAt(createdAt time.Time) *BaseDTO {
	b.CreatedAt = createdAt
	return b
}

func (b *BaseDTO) GetCreatedAt() time.Time {
	return b.CreatedAt
}

func (b *BaseDTO) SetUpdatedAt(updatedAt time.Time) *BaseDTO {
	b.UpdatedAt = updatedAt
	return b
}

func (b *BaseDTO) GetUpdatedAt() time.Time {
	return b.UpdatedAt
}

func (b *BaseDTO) SetDeletedAt(deletedAt time.Time) *BaseDTO {
	b.DeletedAt.Time.Time = deletedAt
	b.DeletedAt.Time.Valid = true
	return b
}

func (b *BaseDTO) GetDeletedAt() time.Time {
	return b.DeletedAt.Time.Time
}

func (b *BaseDTO) SetActive(active bool) *BaseDTO {
	b.Active.Bool = active
	b.Active.Valid = true
	
	return b
}

func (b *BaseDTO) GetActive() bool {
	return b.Active.Bool
}

```
