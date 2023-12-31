package scanner

import (
	"reflect"
	"strings"
	"time"
)

const (
	AllFields = "all"
	Create    = "create"
	Update    = "update"
	Upsert    = "upsert"
	Conflict  = "conflict"
)

// Scanner интерфейс для сканирования таблиц
type Scanner interface {
	RegisterTable(entities ...Tabler)
	OperationFields(table Tabler, operation string) []*Field
	OperationFieldsName(tableName string, operation string) []string
	Table(tableName string) Table
	Tables() map[string]Table
}

// Tabler интерфейс для сущностей таблицы
type Tabler interface {
	TableName() string
	OnCreate() []string
	FieldsPointers() []interface{}
}

// TableUpdater интерфейс для обновления таблиц
type TableUpdater interface {
	SetUpdatedAt(updatedAt time.Time) Tabler
}

// Table структура таблицы
type Table struct {
	Name            string
	Fields          []*Field
	FieldsMap       map[string]*Field
	Constraints     []Constraint
	OperationFields map[string][]*Field
	Entity          Tabler
}

// TableScanner сканер таблиц
type TableScanner struct {
	tables map[string]Table
}

// NewTableScanner конструктор
func NewTableScanner() Scanner {
	return &TableScanner{}
}

// RegisterTable регистрация сущностей
func (t *TableScanner) RegisterTable(entities ...Tabler) {
	tableEntities := make(map[string]Tabler, len(entities))
	t.tables = make(map[string]Table, len(entities))
	for i := range entities {
		tableEntities[entities[i].TableName()] = entities[i]
	}

	for name, entity := range tableEntities {
		table := Table{
			Name:            name,
			FieldsMap:       make(map[string]*Field),
			OperationFields: make(map[string][]*Field),
			Entity:          entity,
		}
		reflected := reflect.TypeOf(entity).Elem()

		for i := 0; i < reflected.NumField(); i++ {
			// Get the field, returns https://golang.org/pkg/reflect/#StructField
			structField := reflected.Field(i)
			// Get the structField tag value
			fieldName := structField.Tag.Get("db")

			if fieldName == "" || fieldName == "-" {
				continue
			}

			field := &Field{
				IDx:     i,
				Name:    fieldName,
				Type:    structField.Tag.Get("db_type"),
				Default: structField.Tag.Get("db_default"),
				Table:   &table,
			}
			constraintRaw := structField.Tag.Get("db_index")
			constraintPieces := strings.Split(constraintRaw, ",")
			if len(constraintPieces) < 1 {
				field.Constraint = Constraint{}
			}
			if len(constraintPieces) > 0 {
				for i := range constraintPieces {
					switch constraintPieces[i] {
					case "index":
						field.Constraint.Index = true
					case "unique":
						field.Constraint.Unique = true
					}
				}
			}
			if field.Constraint.Index {
				field.Constraint.Field = field
				table.Constraints = append(table.Constraints, field.Constraint)
			}
			table.Fields = append(table.Fields, field)
			table.FieldsMap[field.Name] = field

			opsRaw := structField.Tag.Get("db_ops")
			ops := strings.Split(opsRaw, ",")
			if opsRaw != "" {
				for j := range ops {
					table.OperationFields[ops[j]] = append(table.OperationFields[ops[j]], field)
				}
			}

			table.OperationFields[AllFields] = append(table.OperationFields[AllFields], field)
		}

		t.tables[name] = table
	}
}

// OperationFieldsName получение полей для операции над таблицей
func (t *TableScanner) OperationFieldsName(tableName string, operation string) []string {
	fields := t.tables[tableName].OperationFields[operation]
	var fieldsName []string
	for i := range fields {
		fieldsName = append(fieldsName, fields[i].Name)
	}

	return fieldsName
}

// OperationFields получение полей для операции над таблицей
func (t *TableScanner) OperationFields(table Tabler, operation string) []*Field {
	fields := t.tables[table.TableName()].OperationFields[operation]
	pointers := table.FieldsPointers()

	for i := range fields {
		fields[i].Pointer = pointers[fields[i].IDx]
	}

	return fields
}

// Table получение таблицы
func (t *TableScanner) Table(tableName string) Table {
	return t.tables[tableName]
}

// Tables получение таблиц
func (t *TableScanner) Tables() map[string]Table {
	return t.tables
}

// Field структура полей
type Field struct {
	IDx        int
	Name       string
	Type       string
	Default    string
	Constraint Constraint
	Table      *Table
	Pointer    interface{}
}

// Constraint структура ограничения
type Constraint struct {
	Index  bool
	Unique bool
	Field  *Field
}
