# cli-orm-gen

## Что это такое?

- Генерация методов для модели структуры.
- Генерация файлов хранилища и интерфейса CRUD реализованная для заданной модели.

## Как работает?
```bash
$ go build main.go
$ ./main -entity="./путь до модельки" -output="./директория для вывода"
```
Все нужные методы для дальнейшей генерации CRUD будут по указанному адресу (default=./repository)
Также можно скачать бинарник в releases

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

```
