package dao

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Alexandrhub/cli-orm-gen/infrastructure/db/scanner"
	"github.com/Alexandrhub/cli-orm-gen/utils"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const (
	DriverMysql    = "mysql"
	DriverPostgres = "postgres"
	DriverSqlite3  = "sqlite3"
	DriverRamsql   = "ramsql"
)

//go:generate mockgen -source=./sql_adapter.go -destination=../../mock/adapter_mock.go -package=mock
type DAOFace interface {
	Create(ctx context.Context, entity scanner.Tabler, opts ...interface{}) error
	Upsert(ctx context.Context, entities []scanner.Tabler, opts ...interface{}) error
	GetCount(ctx context.Context, entity scanner.Tabler, condition utils.Condition, opts ...interface{}) (uint64, error)
	List(ctx context.Context, dest interface{}, table scanner.Tabler, condition utils.Condition, opts ...interface{}) error
	Update(ctx context.Context, entity scanner.Tabler, condition utils.Condition, operation string, opts ...interface{}) error
}

type DAO struct {
	db         *sqlx.DB
	scanner    scanner.Scanner
	sqlBuilder sq.StatementBuilderType
}

func NewDAO(db *sqlx.DB, dbConf utils.DB, scanner scanner.Scanner) *DAO {
	var builder sq.StatementBuilderType
	if dbConf.Driver != "mysql" {
		builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	}

	return &DAO{db: db, scanner: scanner, sqlBuilder: builder}
}

func (s *DAO) Create(ctx context.Context, table scanner.Tabler, opts ...interface{}) error {
	createFields, createFieldsPointers := s.getFields(table, scanner.Create)

	queryRaw := s.sqlBuilder.Insert(table.TableName()).Columns(createFields...).Values(createFieldsPointers...)

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

func (s *DAO) getFields(entity scanner.Tabler, operation string) ([]string, []interface{}) {
	fields := s.scanner.OperationFields(entity, operation)
	var fieldsPointers []interface{}
	var fieldsName []string
	for i := range fields {
		fieldsPointers = append(fieldsPointers, fields[i].Pointer)
		fieldsName = append(fieldsName, fields[i].Name)
	}

	return fieldsName, fieldsPointers
}

func (s *DAO) Upsert(ctx context.Context, entities []scanner.Tabler, opts ...interface{}) error {
	if len(entities) < 1 {
		return fmt.Errorf("SQL adapter: zero entities passed")
	}
	createFields, _ := s.getFields(entities[0], scanner.Create)
	queryRaw := s.sqlBuilder.Insert(entities[0].TableName()).Columns(createFields...)

	for i := range entities {
		_, createFieldsPointers := s.getFields(entities[i], "create")
		queryRaw = queryRaw.Values(createFieldsPointers...)
	}

	query, args, err := queryRaw.ToSql()
	if err != nil {
		return err
	}

	conflictFields, _ := s.getFields(entities[0], scanner.Conflict)
	if len(conflictFields) > 0 {
		query = query + " ON CONFLICT (%s)"
		query = fmt.Sprintf(query, strings.Join(conflictFields, ","))
		query = query + " DO UPDATE SET"
	}
	upsertFields, _ := s.getFields(entities[0], scanner.Upsert)
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

func (s *DAO) buildSelect(tableName string, condition utils.Condition, fields ...string) (string, []interface{}, error) {
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

func (s *DAO) GetCount(ctx context.Context, entity scanner.Tabler, condition utils.Condition, opts ...interface{}) (uint64, error) {
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

func (s *DAO) List(ctx context.Context, dest interface{}, table scanner.Tabler, condition utils.Condition, opts ...interface{}) error {
	fields, _ := s.getFields(table, scanner.AllFields)
	query, args, err := s.buildSelect(table.TableName(), condition, fields...)
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

func (s *DAO) Update(ctx context.Context, entity scanner.Tabler, condition utils.Condition, operation string, opts ...interface{}) error {
	ent := entity
	updateFields, updateFieldsPointers := s.getFields(entity, operation)

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

func getTransaction(opts ...interface{}) *sqlx.Tx {
	for _, opt := range opts {
		switch opt.(type) {
		case *sqlx.Tx:
			return opt.(*sqlx.Tx)
		}
	}
	return nil
}
