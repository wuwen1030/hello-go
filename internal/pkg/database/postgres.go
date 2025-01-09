package database

import (
	"fmt"
	"log"

	"github.com/wuwen/hello-go/internal/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresClient struct {
	db *gorm.DB
}

func (c *postgresClient) GetDB() (*gorm.DB, error) {
	return c.db, nil
}

func newPostgresClient(cfg *config.DatabaseConfig) (DBClient, error) {
	// 先连接 PostgreSQL（不指定数据库）
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%d sslmode=disable",
		cfg.Host, cfg.Username, cfg.Password, cfg.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("failed to connect to PostgreSQL: %v", err)
		return nil, err
	}

	// 创建数据库
	createDB := fmt.Sprintf("CREATE DATABASE %s", cfg.DBName)
	if err := db.Exec(createDB).Error; err != nil {
		log.Printf("failed to create database: %v", err)
		// 如果数据库已存在，则忽略错误
		if err.Error() != fmt.Sprintf(`pq: database "%s" already exists`, cfg.DBName) {
			return nil, err
		}
	}

	// 连接指定的数据库
	dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host, cfg.Username, cfg.Password, cfg.DBName, cfg.Port)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	return &postgresClient{db: db}, nil
}
