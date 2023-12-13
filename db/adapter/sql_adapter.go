package adapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	
	"github.com/Alexandrhub/cli-orm-gen/utils"
	
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const (
	AllFields = "all"
	Create    = "create"
	Upsert    = "upsert"
	Conflict  = "conflict"
)

//go:generate mockgen -source=./sql_adapter.go -destination=../../mock/adapter_mock.go -package=mock
// Adapter интерфейс описаня методов для работы с базой данных
type Adapter interface {
	Create(ctx context.Context, entity utils.Tabler, opts ...interface{}) error
	Upsert(ctx context.Context, entities []utils.Tabler, opts ...interface{}) error
	GetCount(ctx context.Context, entity utils.Tabler, condition utils.Condition, opts ...interface{}) (uint64, error)
	List(ctx context.Context, dest interface{}, tableName string, condition utils.Condition, opts ...interface{}) error
	Update(ctx context.Context, entity utils.Tabler, condition utils.Condition, operation string, opts ...interface{}) error
}

// SQLAdapter структура адаптера
type SQLAdapter struct {
	db         *sqlx.DB
	scanner    utils.Scanner
	sqlBuilder sq.StatementBuilderType
}

// NewSqlAdapter конструктор
func NewSqlAdapter(db *sqlx.DB, dbConf utils.DB, scanner utils.Scanner) *SQLAdapter {
	var builder sq.StatementBuilderType
	if dbConf.Driver != "mysql" {
		builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	}
	
	return &SQLAdapter{db: db, scanner: scanner, sqlBuilder: builder}
}

// Create создание записи в базе
func (s *SQLAdapter) Create(ctx context.Context, entity utils.Tabler, opts ...interface{}) error {
	createFields := s.scanner.OperationFields(entity.TableName(), Create)
	createFieldsPointers := GetFieldsPointers(entity, "create")
	
	queryRaw := s.sqlBuilder.Insert(entity.TableName()).Columns(createFields...).Values(createFieldsPointers...)
	
	query, args, err := queryRaw.ToSql()
	if err != nil {
		return err
	}
	
	if tx := getTransaction(opts); tx != nil {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = s.db.ExecContext(ctx, query, args...)
	}
	
	return err
}

// Upsert обновление записи в базе
func (s *SQLAdapter) Upsert(ctx context.Context, entities []utils.Tabler, opts ...interface{}) error {
	if len(entities) < 1 {
		return fmt.Errorf("SQL adapter: zero entities passed")
	}
	createFields := s.scanner.OperationFields(entities[0].TableName(), Create)
	queryRaw := s.sqlBuilder.Insert(entities[0].TableName()).Columns(createFields...)
	
	for i := range entities {
		createFieldsPointers := GetFieldsPointers(entities[i], "create")
		queryRaw = queryRaw.Values(createFieldsPointers...)
	}
	
	query, args, err := queryRaw.ToSql()
	if err != nil {
		return err
	}
	
	conflictFields := s.scanner.OperationFields(entities[0].TableName(), Conflict)
	if len(conflictFields) > 0 {
		query = query + " ON CONFLICT (%s)"
		query = fmt.Sprintf(query, strings.Join(conflictFields, ","))
		query = query + " DO UPDATE SET"
	}
	upsertFields := s.scanner.OperationFields(entities[0].TableName(), Upsert)
	for _, field := range upsertFields {
		query += fmt.Sprintf(" %s = excluded.%s,", field, field)
	}
	if len(upsertFields) > 0 {
		query = query[0 : len(query)-1]
	}
	
	if tx := getTransaction(opts); tx != nil {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = s.db.ExecContext(ctx, query, args...)
	}
	
	return err
}

// buildSelect собирает запрос
func (s *SQLAdapter) buildSelect(tableName string, condition utils.Condition, fields ...string) (string, []interface{}, error) {
	if condition.ForUpdate {
		temp := []string{"FOR UPDATE"}
		temp = append(temp, fields...)
		fields = temp
	}
	queryRaw := s.sqlBuilder.Select(fields...).From(tableName)
	
	if condition.Equal != nil {
		for field, args := range condition.Equal {
			queryRaw = queryRaw.Where(sq.Eq{field: args})
		}
	}
	
	if condition.NotEqual != nil {
		for field, args := range condition.NotEqual {
			queryRaw = queryRaw.Where(sq.NotEq{field: args})
		}
	}
	
	if condition.Order != nil {
		for _, order := range condition.Order {
			direction := "DESC"
			if order.Asc {
				direction = "ASC"
			}
			queryRaw = queryRaw.OrderBy(fmt.Sprintf("%s %s", order.Field, direction))
		}
	}
	
	if condition.LimitOffset != nil {
		if condition.LimitOffset.Limit > 0 {
			queryRaw.Limit(uint64(condition.LimitOffset.Limit))
		}
		if condition.LimitOffset.Offset > 0 {
			queryRaw.Offset(uint64(condition.LimitOffset.Offset))
		}
	}
	
	return queryRaw.ToSql()
}

// GetCount считает количество записей COUNT(*)
func (s *SQLAdapter) GetCount(ctx context.Context, entity utils.Tabler, condition utils.Condition, opts ...interface{}) (uint64, error) {
	query, args, err := s.buildSelect(entity.TableName(), condition, "COUNT(*)")
	if err != nil {
		return 0, err
	}
	
	var rows *sqlx.Rows
	if tx := getTransaction(opts); tx != nil {
		rows, err = tx.QueryxContext(ctx, query, args...)
	} else {
		rows, err = s.db.QueryxContext(ctx, query, args...)
	}
	
	if err != nil {
		return 0, err
	}
	
	var count uint64
	// iterate over each row
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	// check the error from rows
	err = rows.Err()
	
	return count, err
}

// List возвращает список записей по условию
func (s *SQLAdapter) List(ctx context.Context, dest interface{}, tableName string, condition utils.Condition, opts ...interface{}) error {
	fields := s.scanner.OperationFields(tableName, AllFields)
	query, args, err := s.buildSelect(tableName, condition, fields...)
	if err != nil {
		return err
	}
	
	if tx := getTransaction(opts); tx != nil {
		err = tx.SelectContext(ctx, dest, query, args...)
	} else {
		err = s.db.SelectContext(ctx, dest, query, args...)
	}
	
	return err
}

// Update обновляет запись по условию
func (s *SQLAdapter) Update(ctx context.Context, entity utils.Tabler, condition utils.Condition, operation string, opts ...interface{}) error {
	ent := entity
	updateFields := s.scanner.OperationFields(entity.TableName(), operation)
	
	updateFieldsPointers := GetFieldsPointers(entity, operation)
	
	updateRaw := s.sqlBuilder.Update(ent.TableName())
	
	if condition.Equal != nil {
		for field, args := range condition.Equal {
			updateRaw = updateRaw.Where(sq.Eq{field: args})
		}
	}
	
	if condition.NotEqual != nil {
		for field, args := range condition.NotEqual {
			updateRaw = updateRaw.Where(sq.NotEq{field: args})
		}
	}
	
	for i := range updateFields {
		updateRaw = updateRaw.Set(updateFields[i], updateFieldsPointers[i])
	}
	
	query, args, err := updateRaw.ToSql()
	if err != nil {
		return err
	}
	
	var res sql.Result
	if tx := getTransaction(opts); tx != nil {
		res, err = tx.ExecContext(ctx, query, args...)
	} else {
		res, err = s.db.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	
	return err
}

// getTransaction транзакция с опциями sqlx.Tx
func getTransaction(opts ...interface{}) *sqlx.Tx {
	for _, opt := range opts {
		switch opt.(type) {
		case *sqlx.Tx:
			return opt.(*sqlx.Tx)
		}
	}
	return nil
}
