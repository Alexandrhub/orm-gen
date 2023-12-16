package main

import (
	"time"

	"github.com/Alexandrhub/cli-orm-gen/infrastructure/db/types"
)

//go:generate easytags $GOFILE json,db,db_ops,db_type,db_default,mapper
type Bot struct { // DTO - data transfer object
	ID          int            `json:"id" db:"id" db_type:"INTEGER PRIMARY KEY" db_default:"not null" db_ops:"id" mapper:"id"`
	UserID      int            `json:"user_id" db:"user_id" db_ops:"create" db_type:"int" db_default:"default 1" db_index:"index" mapper:"user_id"`
	Name        string         `json:"name" db:"name" db_ops:"create,update" db_type:"varchar(55)" db_default:"not null" mapper:"name"`
	Description string         `json:"description" db:"description" db_ops:"create,update" db_type:"varchar(144)" db_default:"not null" mapper:"description"`
	OrderCount  int            `json:"order_count" db:"order_count" db_ops:"create,update" db_type:"int" db_default:"default 1" mapper:"order_count"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at" db_type:"timestamp" db_default:"default (date()) not null" db_index:"index" db_ops:"created_at" mapper:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at" db_ops:"update" db_type:"timestamp" db_default:"default (date()) not null" db_index:"index" mapper:"updated_at"`
	DeletedAt   types.NullTime `json:"deleted_at" db:"deleted_at" db_ops:"update" db_type:"timestamp" db_default:"default null" db_index:"index" db_ops:"deleted_at" mapper:"deleted_at"`
}

func (b *Bot) TableName() string {
	return "bots"
}

func (b *Bot) OnCreate() []string {
	return []string{}
}

func (b *Bot) FieldsPointers() []interface{} {
	return []interface{}{
		&b.ID,
		&b.UserID,
		&b.Name,
		&b.Description,
		&b.OrderCount,
		&b.CreatedAt,
		&b.UpdatedAt,
		&b.DeletedAt,
	}
}
