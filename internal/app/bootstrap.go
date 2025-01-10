package app

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/wuwen/hello-go/internal/model"
	"github.com/wuwen/hello-go/internal/pkg/config"
	"github.com/wuwen/hello-go/internal/pkg/database"
	"gorm.io/gorm"
)

func (a *App) Initialize() error {
	if err := a.loadConfig(); err != nil {
		return err
	}

	db, err := a.initDatabase()
	if err != nil {
		return err
	}

	if err := a.initRoles(db); err != nil {
		return err
	}

	if err := a.initAdmin(db); err != nil {
		return err
	}

	if err := a.initCasbin(db); err != nil {
		return err
	}

	a.setupDependencies(db)

	return nil
}

func (a *App) loadConfig() error {
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}
	a.config = cfg

	gin.SetMode(cfg.Server.Mode)
	return nil
}

func (a *App) initDatabase() (*gorm.DB, error) {
	dbClient, err := database.NewDBClient(&a.config.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	db, err := dbClient.GetDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %v", err)
	}

	// 自动迁移数据库表
	if err := db.AutoMigrate(&model.Role{}, &model.User{}, &model.Article{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	return db, nil
}

func (a *App) initRoles(db *gorm.DB) error {
	// 创建管理员角色
	adminRole := model.Role{
		Name: "admin",
	}
	if err := db.Where("name = ?", adminRole.Name).FirstOrCreate(&adminRole).Error; err != nil {
		return fmt.Errorf("failed to create admin role: %v", err)
	}
	// 创建用户角色
	userRole := model.Role{
		Name: "user",
	}
	if err := db.Where("name = ?", userRole.Name).FirstOrCreate(&userRole).Error; err != nil {
		return fmt.Errorf("failed to create user role: %v", err)
	}

	return nil
}

func (a *App) initAdmin(db *gorm.DB) error {
	adminUser := model.User{
		Username: "admin",
		Email:    "admin@example.com",
		Status:   model.UserStatusActive,
	}
	if err := adminUser.SetPassword("admin123"); err != nil {
		return fmt.Errorf("failed to set admin password: %v", err)
	}

	if err := db.Where("username = ?", adminUser.Username).FirstOrCreate(&adminUser).Error; err != nil {
		return fmt.Errorf("failed to create admin user: %v", err)
	}

	return nil
}

func (a *App) initCasbin(db *gorm.DB) error {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return fmt.Errorf("failed to create casbin adapter: %v", err)
	}

	enforcer, err := casbin.NewEnforcer("configs/model.conf", adapter)
	if err != nil {
		return fmt.Errorf("failed to create casbin enforcer: %v", err)
	}
	a.enforcer = enforcer

	if err := enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("failed to load casbin policy: %v", err)
	}

	if err := a.initializeDefaultPolicies(); err != nil {
		return fmt.Errorf("failed to initialize default policies: %v", err)
	}

	return nil
}
