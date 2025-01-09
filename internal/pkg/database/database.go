package database

import (
	"fmt"

	"github.com/wuwen/hello-go/internal/pkg/config"
	"gorm.io/gorm"
)

// DBClient 数据库客户端接口
type DBClient interface {
	GetDB() (*gorm.DB, error)
}

func NewDBClient(cfg *config.DatabaseConfig) (DBClient, error) {
	switch cfg.Driver {
	case "mysql":
		return newMySQLClient(cfg)
	case "postgres":
		return newPostgresClient(cfg)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}
}
