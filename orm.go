package cli_orm_gen

import (
	"github.com/Alexandrhub/cli-orm-gen/db"
	"github.com/Alexandrhub/cli-orm-gen/db/adapter"
	"github.com/Alexandrhub/cli-orm-gen/infrastructure/db/migrate"
	"github.com/Alexandrhub/cli-orm-gen/utils"
	
	"go.uber.org/zap"
)

func NewOrm(dbConf utils.DB, scanner utils.Scanner, logger *zap.Logger) (*adapter.SQLAdapter, error) {
	var (
		orm *db.SqlDB
		err error
	)
	
	orm, err = db.NewSqlDB(dbConf, scanner, logger)
	if err != nil {
		logger.Fatal("error init db", zap.Error(err))
	}
	migrator := migrate.NewMigrator(orm.DB, dbConf, scanner)
	err = migrator.Migrate()
	if err != nil {
		logger.Fatal("migrator err", zap.Error(err))
	}
	return orm.SqlAdapter, nil
}
