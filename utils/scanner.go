package utils

import (
	"github.com/Alexandrhub/cli-orm-gen/infrastructure/db/scanner"
)

const (
	Update = "update"
	Create = "create"
)

// NewTableScanner создание сканера
func NewTableScanner() Scanner {
	return &scanner.TableScanner{}
}

// Scanner интерфейс для сканирования таблиц
type Scanner interface {
	RegisterTable(entities ...scanner.Tabler)
	OperationFields(tableName, operation string) []string
	Table(tableName string) scanner.Table
	Tables() map[string]scanner.Table
}

// Tabler интерфейс для сущностей таблицы
type Tabler interface {
	TableName() string
	OnCreate() []string
}
