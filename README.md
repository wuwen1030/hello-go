# CMS Service

一个使用 Go 语言开发的内容管理系统 API 服务。

## 项目介绍

本项目是一个简单的内容管理系统（CMS）的后端 API 服务。它提供了用户注册、登录以及文章的创建、读取、更新和删除等基本功能。

## 技术栈

- **Go 1.21+**:  后端开发语言
- **Gin Web Framework**:  轻量级 Web 框架
- **GORM**:  Go 语言 ORM 库
- **MySQL**:  关系型数据库
- **Swagger**:  API 文档生成工具
- **Testify**:  Go 语言测试框架

## 如何运行项目

### 前置要求

- Go 1.21+
- MySQL 5.7+
- Make (可选)

### 安装步骤

1.  **克隆仓库**

```bash
git clone https://github.com/yourusername/cms-service.git
cd cms-service
```

2.  **安装依赖**

```bash
go mod download
```

3.  **配置数据库**

- 修改 `configs/config.yaml` 文件，配置数据库连接信息。

4.  **生成 Swagger 文档**

```bash
swag init -g cmd/server/main.go
```

5.  **运行服务**

```bash
go run cmd/server/main.go
```

### API 文档

启动服务后，访问 Swagger 文档：http://localhost:8080/swagger/index.html

### 运行测试
```
go test -v ./test
```

## 项目结构

```
.
├── cmd/ # 主要的应用程序入口
│ └── server/ # HTTP 服务器
│ └── main.go # 主程序入口
├── configs/ # 配置文件
│ └── config.yaml # 应用配置文件
├── internal/ # 私有应用程序和库代码
│ ├── app/ # 应用程序初始化和配置
│ │ └── app.go # 应用程序核心结构
│ ├── handler/ # HTTP 处理器
│ │ └── article.go # 文章相关处理器
│ │ └── user.go # 用户相关处理器
│ ├── model/ # 数据库模型
│ │ └── article.go # 文章模型定义
│ │ └── user.go # 用户模型定义
│ ├── pkg/ # 内部公共包
│ │ ├── auth/ # 认证相关
│ │ ├── config/ # 配置加载
│ │ ├── database/ # 数据库连接
│ │ └── response/ # 响应处理
│ ├── repository/ # 数据访问层
│ │ └── article.go # 文章数据访问
│ │ └── user.go # 用户数据访问
│ ├── router/ # 路由定义
│ │ ├── api/ # API 路由
│ │ └── router.go # 路由接口
│ └── service/ # 业务逻辑层
│ └── article.go # 文章业务逻辑
│ └── user.go # 用户业务逻辑
├── test/ # 测试代码
│ └── api_test.go # API 测试
├── docs/ # Swagger 文档（自动生成）
├── go.mod # Go 模块定义
├── go.sum # Go 依赖版本锁定
└── README.md # 项目说明文档
```

## 功能特性

- RESTful API
- 用户注册和登录
- 文章管理（CRUD）
- JWT 认证
- 统一的错误处理
- API 文档（Swagger）
- 自动化测试
- 优雅关闭
