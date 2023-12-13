package migrate

import (
	"context"
	"fmt"
	"strings"
	
	"github.com/Alexandrhub/cli-orm-gen/infrastructure/db/scanner"
	"github.com/Alexandrhub/cli-orm-gen/utils"
	
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/sync/errgroup"
)

//go:generate qtc -dir=./

type Scanner interface {
	OperationFields(tableName, operation string) []string
	Tables() map[string]scanner.Table
}

type Migrator struct {
	db      *sqlx.DB
	dbConf  utils.DB
	scanner Scanner
}

func NewMigrator(db *sqlx.DB, dbConf utils.DB, scanner Scanner) *Migrator {
	return &Migrator{db: db, dbConf: dbConf, scanner: scanner}
}

func (m *Migrator) Migrate() error {
	tables := m.scanner.Tables()
	var err error
	var query string
	var args []interface{}
	var builder sq.StatementBuilderType
	var schema string
	if m.dbConf.Driver == "mysql" {
		builder = sq.StatementBuilder.PlaceholderFormat(sq.Question)
		schema = m.dbConf.Name
	}
	if m.dbConf.Driver == "postgres" {
		builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
		schema = "public"
	}
	errGroup, ctx := errgroup.WithContext(context.Background())
	for name := range tables {
		table := tables[name]
		var tableFields []string
		queryRaw := builder.Select("COLUMN_NAME").From("INFORMATION_SCHEMA.COLUMNS")
		queryRaw = queryRaw.Where(sq.Eq{"TABLE_SCHEMA": schema, "TABLE_NAME": table.Name})
		query, args, err = queryRaw.ToSql()
		err = m.db.Select(&tableFields, query, args...)
		if err != nil {
			return fmt.Errorf("%s, %s", err, query)
		}
		tableFieldsMap := make(map[string]string, len(tableFields))
		for i := range tableFields {
			tableFieldsMap[tableFields[i]] = tableFields[i]
		}
		if len(tableFields) < 1 {
			errGroup.Go(
				func() error {
					createQuery := CreateTable(table, m.dbConf)
					queries := strings.Split(createQuery, ";")
					for i := range queries {
						queries[i] = strings.TrimSpace(queries[i])
						if queries[i] == "" {
							continue
						}
						_, err = m.db.ExecContext(ctx, queries[i])
						if err != nil {
							if v, ok := err.(*pq.Error); ok {
								if v.Code != "42P07" {
									return fmt.Errorf("%s, %s", err, queries[i])
								}
							} else {
								return fmt.Errorf("%s, %s", err, queries[i])
							}
						}
					}
					return nil
				},
			)
		}
		if len(tableFields) > 0 {
			errGroup.Go(
				func() error {
					entityFields := m.scanner.OperationFields(table.Name, scanner.AllFields)
					diff := make(map[string]scanner.Field, len(entityFields))
					for i := range entityFields {
						if _, ok := tableFieldsMap[entityFields[i]]; !ok {
							diff[entityFields[i]] = table.FieldsMap[entityFields[i]]
						}
					}
					for fieldName := range diff {
						if fieldName == "" {
							continue
						}
						alterQuery := AlterTable(diff[fieldName])
						queries := strings.Split(alterQuery, ";")
						for i := range queries {
							queries[i] = strings.TrimSpace(queries[i])
							if queries[i] == "" {
								continue
							}
							_, err = m.db.Queryx(queries[i])
							if err != nil {
								return fmt.Errorf("%s, %s", err, queries[i])
							}
						}
					}
					
					return nil
				},
			)
		}
	}
	
	return errGroup.Wait()
}
