# Blog

一个前后端一体的个人博客系统。项目提供公开博客站点、独立的管理后台，以及基于 REST API 的 Go 服务端。

- 后端：Go 1.22、Gin、GORM、MySQL、JWT、Viper
- 前端：Vue 3、Vite、Vue Router、Pinia、Axios、Marked
- 部署方式：构建前端后，由 Go 服务直接托管 `web/dist` 中的静态资源

## 功能

### 公开站点

- 文章列表、详情、阅读量、最新文章与热门文章
- 分类、标签和归档浏览
- 站点统计与按月文章活动数据
- Markdown 文章内容渲染
- 用户注册、登录与会话恢复
- 游客可填写昵称、邮箱和可选网站发表评论；登录用户以账户身份发表评论并可删除自己的评论
- 文章、分类和标签仅可由管理员在后台管理
- 图片上传（JPEG、PNG、WebP）

### 管理后台

- 管理仪表盘
- 文章、分类、标签的新增、编辑和删除
- 评论审核与删除
- 用户列表与用户状态管理
- 管理员权限校验

## 项目结构

```text
.
├── cmd/server/       # 后端启动入口
├── config/           # 配置加载与配置示例
├── docs/             # API、配置、部署与本地运行文档
├── internal/         # 控制器、服务、DAO、模型、路由和中间件
├── migrations/       # 数据库初始化 SQL
├── pkg/              # 数据库、JWT、响应与上传等基础组件
├── scripts/          # 辅助脚本
├── tests/            # 后端测试
├── uploads/          # 本地上传文件目录
└── web/              # Vue 前端
```

## 快速开始

### 环境要求

- Go 1.22+
- MySQL 8.x
- Node.js 20+
- npm 10+

### 1. 配置后端

复制并修改 [配置示例](config/config.example.yaml)：

```powershell
Copy-Item config/config.example.yaml config/config.yaml
```

至少配置数据库连接和一个安全的 JWT 密钥：

```yaml
database:
  host: "127.0.0.1"
  port: 3306
  user: "root"
  pass: "your-password"
  name: "gin_blog"

jwt:
  secret: "replace-with-a-long-random-secret"
```

配置也可由环境变量覆盖。例如 `BLOG_DATABASE_HOST`、`BLOG_DATABASE_PASS`、`BLOG_JWT_SECRET`；通过 `BLOG_CONFIG_FILE` 可指定其他 YAML 配置文件。完整字段见 [配置说明](docs/configuration.md)。

如需在首次启动时创建管理员，设置：

```yaml
admin_bootstrap:
  enabled: true
  username: "admin"
  password: "change-me"
  nickname: "管理员"
```

首次创建成功后，建议将 `enabled` 改回 `false`，并妥善保管密码与 JWT 密钥。

### 2. 初始化数据库

服务启动时会使用 GORM 自动迁移表结构。也可以先手动执行 [migrations/init.sql](migrations/init.sql)：

```powershell
mysql -u root -p gin_blog < migrations/init.sql
```

### 3. 启动后端

在项目根目录运行：

```powershell
go run .
```

后端会同时启动两个 HTTP 服务：

| 服务 | 默认地址 | 用途 |
| --- | --- | --- |
| 公开服务 | `http://localhost:8080` | 公开 API、上传资源和博客站点 |
| 管理服务 | `http://localhost:8081` | 管理 API 和管理后台 |

两个服务均提供 `GET /health` 健康检查。端口可通过 `server.port` 与 `admin_server.port` 修改，且不能相同。

### 4. 启动前端开发服务器

在 `web` 目录安装依赖：

```powershell
cd web
npm install
```

分别启动公开站点和管理后台：

```powershell
npm run dev        # http://localhost:3000
npm run dev:admin  # http://localhost:3001/admin
```

也可以一次启动两个开发服务器：

```powershell
npm run dev:all
```

Vite 在公开模式下会将 `/api` 代理到 `8080`；在管理模式下会将 `/api` 代理到 `8081`。`/uploads` 始终代理到公开服务 `8080`。

### 5. 构建并使用一体化服务

```powershell
cd web
npm run build
cd ..
go run .
```

构建产物生成在 `web/dist`。目录存在时，两个 Go 服务会回退到前端的 `index.html`，因此可分别从 `http://localhost:8080` 与 `http://localhost:8081/admin` 访问已构建的站点。

## API 概览

所有业务接口位于 `/api` 下。

- 公开接口：分类、标签、文章、评论、站点统计与文章活动
- 认证接口：`POST /api/auth/register`、`POST /api/auth/login`
- 登录接口：会话查询、删除自己的评论
- 管理接口：`/api/admin/*`，需要管理员 JWT

具体请求参数、响应格式与权限要求见 [API 文档](docs/api.md)。

## Docker 部署

项目提供多阶段构建镜像和 Docker Compose 编排，可同时启动应用与 MySQL，并使用命名卷持久化数据库和上传文件：

```powershell
Copy-Item docker/.env.example docker/.env
# 编辑 docker/.env，替换数据库密码和 BLOG_JWT_SECRET
docker compose --env-file docker/.env -f docker/compose.yaml up -d --build
```

Docker 会重新构建前后端；公开端口仅用于浏览和评论，文章发布、分类标签维护及图片上传请使用管理员后台端口。

公开站点默认访问 `http://localhost:8080`，管理后台默认访问 `http://localhost:8081/admin`。更多操作、数据持久化与生产注意事项见 [Docker 部署说明](docs/deployment.md#docker-compose-部署)。

## 上传限制

上传功能仅在同时配置 `upload.dir` 和 `upload.url` 时启用。默认单文件上限为 5 MiB，接受 `.jpg`、`.jpeg`、`.png` 和 `.webp` 图片；服务会校验实际内容类型，并在写入失败时清理不完整文件。

## 测试与检查

在根目录运行后端测试：

```powershell
go test ./...
```

在 `web` 目录运行前端检查：

```powershell
npm run lint
npm run build
```

GitHub Actions 会在推送和拉取请求中执行上述后端测试、前端 lint 与生产构建。

## 相关文档

- [本地运行说明](docs/setup.md)
- [配置说明](docs/configuration.md)
- [部署检查清单](docs/deployment.md)
- [API 文档](docs/api.md)
