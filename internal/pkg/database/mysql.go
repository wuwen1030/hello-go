package database

import (
	"fmt"
	"log"

	"github.com/wuwen/hello-go/internal/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mysqlClient struct {
	db *gorm.DB
}

func (c *mysqlClient) GetDB() (*gorm.DB, error) {
	return c.db, nil
}

func newMySQLClient(cfg *config.DatabaseConfig) (DBClient, error) {
	// 先连接 MySQL（不指定数据库）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("failed to connect to MySQL: %v", err)
		return nil, err
	}

	// 创建数据库
	createDB := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci", cfg.DBName)
	if err := db.Exec(createDB).Error; err != nil {
		log.Printf("failed to create database: %v", err)
		return nil, err
	}

	// 连接指定的数据库
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
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

	return &mysqlClient{db: db}, nil
}
