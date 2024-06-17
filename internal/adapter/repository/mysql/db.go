package mysql

import (
	"bam/internal/adapter/config"
	m "bam/internal/core/domain"
	"context"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type datebase struct {
	DB *gorm.DB
}

var tables = []interface{}{
	m.LunarDateResponse{},
	m.FormOrdianReq{},
	m.MultiFile{},
}

func NewDatabase(ctx context.Context, config *config.DB, prod bool) (*datebase, error) {
	var logLevel logger.LogLevel
	if prod {
		logLevel = logger.Silent
	} else {
		logLevel = logger.Info
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("Database connected!")

	return &datebase{db}, nil
}

func (db *datebase) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		fmt.Println("Database close error -:", err)
		return err
	}
	sqlDB.Close()
	fmt.Println("Database close!")
	return nil
}

func (db *datebase) Migrate() error {
	tx := db.DB.Begin()
	for _, table := range tables {
		if err := db.DB.AutoMigrate(table); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Error
}

// set max open connections
func (db *datebase) SetMaxOpenConns(n int) error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(n)
	return nil
}
