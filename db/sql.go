package db

import (
	"database/sql"
	"fmt"
	"time"
	
	"github.com/Alexandrhub/cli-orm-gen/db/adapter"
	"github.com/Alexandrhub/cli-orm-gen/utils"
	
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// SqlDB структура для работы с базой данных sqlx.DB и адаптером
type SqlDB struct {
	DB         *sqlx.DB
	SqlAdapter *adapter.SQLAdapter
}

// NewSqlDB конструктор
func NewSqlDB(dbConf utils.DB, scanner utils.Scanner, logger *zap.Logger) (*SqlDB, error) {
	var dsn string
	var err error
	var dbRaw *sql.DB
	
	switch dbConf.Driver {
	case "postgres":
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbConf.Host,
			dbConf.Port,
			dbConf.User,
			dbConf.Password,
			dbConf.Name,
		)
	case "mysql":
		cfg := mysql.NewConfig()
		cfg.Net = dbConf.Net
		cfg.Addr = dbConf.Host
		cfg.User = dbConf.User
		cfg.Passwd = dbConf.Password
		cfg.DBName = dbConf.Name
		cfg.ParseTime = true
		cfg.Timeout = time.Duration(dbConf.Timeout) * time.Second
		dsn = cfg.FormatDSN()
	case "ramsql":
		dsn = "Testing"
	}
	
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	timeoutExceeded := time.After(time.Second * time.Duration(dbConf.Timeout))
	
	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("db connection failed after %d timeout %s", dbConf.Timeout, err)
		case <-ticker.C:
			dbRaw, err = sql.Open(dbConf.Driver, dsn)
			if err != nil {
				return nil, err
			}
			err = dbRaw.Ping()
			if err == nil {
				db := sqlx.NewDb(dbRaw, dbConf.Driver)
				db.SetMaxOpenConns(50)
				db.SetMaxIdleConns(50)
				sqlAdapter := adapter.NewSqlAdapter(db, dbConf, scanner)
				
				return &SqlDB{db, sqlAdapter}, nil
			}
			logger.Error("failed to connect to the database", zap.String("dsn", dsn), zap.Error(err))
		}
	}
}

// Close закрытие соединения с базой данных sqlx.DB
func (s *SqlDB) Close() {
	s.DB.Close()
}
