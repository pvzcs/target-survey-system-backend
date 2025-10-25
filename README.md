# Survey System

目标问卷系统 - 一个专业的问卷管理和数据收集平台

[![Build Status](https://github.com/pvzcs/target-survey-system-backend/workflows/Build/badge.svg)](https://github.com/pvzcs/target-survey-system-backend/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/pvzcs/target-survey-system-backend)](https://goreportcard.com/report/github.com/pvzcs/target-survey-system-backend)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
lic

## 特性

- 🎯 多种题型支持（填空题、单选题、多选题、表格题）
- 🔐 加密链接和预填字段功能
- 🔒 一次性填答机制，防止重复提交
- 📊 数据导出（CSV、Excel）
- 🚀 高性能缓存（Redis）
- 🔑 JWT 认证和授权
- 🛡️ 限流保护
- 📝 完整的 API 文档

## 技术栈

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **ORM**: GORM
- **Database**: MySQL 8.0+
- **Cache**: Redis 7.0+
- **Configuration**: Viper
- **Authentication**: JWT
- **Encryption**: AES-256-GCM

## 项目结构

```
survey-system/
├── cmd/
│   └── server/              # 应用入口
│       └── main.go
├── internal/
│   ├── api/                 # API 层
│   │   ├── handler/         # HTTP 处理器
│   │   ├── middleware/      # 中间件（认证、CORS、限流）
│   │   └── router/          # 路由定义
│   ├── service/             # 业务逻辑层
│   ├── repository/          # 数据访问层
│   ├── model/               # 数据模型
│   ├── dto/                 # 数据传输对象
│   │   ├── request/         # 请求 DTO
│   │   └── response/        # 响应 DTO
│   ├── cache/               # Redis 缓存操作
│   └── config/              # 配置管理
├── pkg/
│   ├── database/            # 数据库工具
│   ├── redis/               # Redis 工具
│   ├── errors/              # 自定义错误
│   ├── utils/               # 工具函数
│   └── constants/           # 常量定义
├── config/                  # 配置文件
│   ├── config.yaml
│   └── config.example.yaml
├── migrations/              # 数据库迁移脚本
├── docs/                    # 文档
│   ├── api.md              # API 文档
│   └── DEPLOYMENT.md       # 部署文档
├── scripts/                 # 工具脚本
├── Dockerfile              # Docker 镜像构建文件
├── docker-compose.yml      # Docker Compose 配置
├── .env.example            # 环境变量示例
├── go.mod
├── go.sum
└── README.md
```

## 快速开始

### 方式一：使用 Docker Compose（推荐）

这是最简单的启动方式，会自动启动 MySQL、Redis 和应用服务。

1. **克隆仓库**

```bash
git clone https://github.com/pvzcs/target-survey-system-backend.git
cd survey-system
```

2. **配置环境变量**

```bash
cp .env.example .env
# 编辑 .env 文件，修改必要的配置（特别是密钥和密码）
```

3. **启动所有服务**

```bash
docker-compose up -d
```

4. **查看日志**

```bash
docker-compose logs -f app
```

5. **访问应用**

```
API: http://localhost:8080
健康检查: http://localhost:8080/health
```

6. **默认管理员账号**

- 用户名: `admin`
- 密码: `admin123`
- **重要**: 首次登录后请立即修改密码！

### 方式二：手动安装

#### 前置要求

- Go 1.21 或更高版本
- MySQL 8.0 或更高版本
- Redis 7.0 或更高版本

#### 安装步骤

1. **克隆仓库**

```bash
git clone https://github.com/pvzcs/target-survey-system-backend.git
cd survey-system
```

2. **安装依赖**

```bash
go mod download
```

3. **配置数据库**

```bash
# 登录 MySQL
mysql -u root -p

# 创建数据库和用户
CREATE DATABASE survey_system CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'survey_user'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON survey_system.* TO 'survey_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;

# 导入数据库结构
mysql -u survey_user -p survey_system < migrations/001_create_tables.sql
mysql -u survey_user -p survey_system < migrations/002_seed_data.sql
```

4. **配置应用**

```bash
# 复制配置文件
cp config/config.example.yaml config/config.yaml

# 或使用环境变量
cp .env.example .env
```

编辑配置文件，设置数据库和 Redis 连接信息。

5. **生成加密密钥**

```bash
# 生成 32 字节的加密密钥
openssl rand -base64 32

# 或使用 Go 脚本
go run scripts/hash_password.go
```

6. **运行应用**

```bash
go run cmd/server/main.go
```

或构建后运行：

```bash
go build -o survey-system ./cmd/server
./survey-system
```

## 配置说明

### 环境变量

应用支持通过环境变量配置，优先级高于配置文件：

```bash
# 服务器配置
SERVER_PORT=8080
SERVER_MODE=release  # debug 或 release

# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=survey_user
DB_PASSWORD=your_password
DB_DATABASE=survey_system

# Redis 配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT 配置
JWT_SECRET=your-secret-key-change-in-production

# 加密配置（必须是 32 字节）
ENCRYPTION_KEY=your-32-byte-encryption-key-here

# CORS 配置
CORS_ALLOWED_ORIGINS=http://localhost:3000

# 限流配置
RATE_LIMIT_REQUESTS_PER_MINUTE=100
```

### 配置文件

也可以使用 YAML 配置文件 `config/config.yaml`，详见 `config/config.example.yaml`。

## API 文档

完整的 API 文档请查看 [docs/api.md](docs/api.md)

### 主要端点

#### 认证

- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/register` - 用户注册

#### 问卷管理（需要认证）

- `POST /api/v1/surveys` - 创建问卷
- `GET /api/v1/surveys` - 获取问卷列表
- `GET /api/v1/surveys/:id` - 获取问卷详情
- `PUT /api/v1/surveys/:id` - 更新问卷
- `DELETE /api/v1/surveys/:id` - 删除问卷
- `POST /api/v1/surveys/:id/publish` - 发布问卷

#### 题目管理（需要认证）

- `POST /api/v1/questions` - 创建题目
- `PUT /api/v1/questions/:id` - 更新题目
- `DELETE /api/v1/questions/:id` - 删除题目
- `PUT /api/v1/surveys/:id/questions/reorder` - 重新排序题目

#### 分享链接（需要认证）

- `POST /api/v1/surveys/:id/share` - 生成分享链接

#### 公开访问（无需认证）

- `GET /api/v1/public/surveys/:id` - 获取问卷（需要 token）
- `POST /api/v1/public/responses` - 提交填答

#### 数据管理（需要认证）

- `GET /api/v1/surveys/:id/responses` - 获取填答记录
- `GET /api/v1/surveys/:id/statistics` - 获取统计信息
- `GET /api/v1/surveys/:id/export` - 导出数据（CSV/Excel）

## 开发

### 构建

```bash
# 开发构建
go build -o survey-system ./cmd/server

# 生产构建（优化）
CGO_ENABLED=1 go build -ldflags="-s -w" -o survey-system ./cmd/server
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行测试并显示覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 代码检查

```bash
# 格式化代码
go fmt ./...

# 静态分析
go vet ./...

# 使用 golangci-lint（推荐）
golangci-lint run
```

## 部署

详细的部署文档请查看 [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)

### Docker 部署

```bash
# 构建镜像
docker build -t survey-system:latest .

# 运行容器
docker run -d \
  --name survey-system \
  -p 8080:8080 \
  -e DB_HOST=mysql \
  -e REDIS_HOST=redis \
  survey-system:latest
```

### 使用 Docker Compose

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

## 安全建议

1. **修改默认密码**: 首次部署后立即修改默认管理员密码
2. **使用强密钥**: 生成强随机密钥用于 JWT 和加密
3. **启用 HTTPS**: 在生产环境使用 HTTPS（配置 Nginx/Caddy）
4. **限制 CORS**: 只允许信任的域名访问 API
5. **定期备份**: 定期备份数据库
6. **更新依赖**: 定期更新 Go 依赖包

## 工具脚本

### 生成密码哈希

```bash
go run scripts/hash_password.go your_password
```

## 监控和维护

### 健康检查

```bash
curl http://localhost:8080/health
```

### 查看日志

```bash
# Docker
docker-compose logs -f app

# Systemd
journalctl -u survey-system -f
```

### 数据库备份

```bash
mysqldump -u survey_user -p survey_system > backup_$(date +%Y%m%d).sql
```

## 故障排查

### 应用无法启动

1. 检查配置文件是否正确
2. 验证数据库和 Redis 连接
3. 查看应用日志
4. 确认端口未被占用

### 数据库连接失败

1. 验证数据库服务运行状态
2. 检查数据库用户权限
3. 确认防火墙规则

详细的故障排查指南请查看 [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)

## 贡献

欢迎贡献代码！请遵循以下步骤：

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 联系方式

- 项目主页: [https://github.com/pvzcs/target-survey-system-backend](https://github.com/pvzcs/target-survey-system-backend)
- 问题反馈: [https://github.com/pvzcs/target-survey-system-backend/issues](https://github.com/pvzcs/target-survey-system-backend/issues)

## 致谢

感谢所有贡献者和开源社区的支持！
