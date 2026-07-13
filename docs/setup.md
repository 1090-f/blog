# 本地运行说明

## 环境要求

- Go 1.22+
- MySQL 8.x
- Node.js 20+
- npm 10+

## 配置与数据库

从配置示例创建本地配置：

```powershell
Copy-Item config/config.example.yaml config/config.yaml
```

编辑 `config/config.yaml`，至少填写 `database.*` 和一个安全的 `jwt.secret`。配置文件和环境变量的优先级及完整字段请见 [配置说明](configuration.md)。

服务会用 GORM `AutoMigrate` 自动迁移表结构；如需手动初始化，可执行：

```powershell
mysql -u root -p gin_blog < migrations/init.sql
```

## 启动后端

在项目根目录运行：

```powershell
go run .
```

应用会启动两个服务：

| 服务 | 地址 | 说明 |
| --- | --- | --- |
| 公开服务 | `http://localhost:8080` | 公开站点、基础 API、上传资源 |
| 管理服务 | `http://localhost:8081` | 管理站点和管理 API |

端口可通过 `server.port` 和 `admin_server.port` 修改。两个服务均可通过 `GET /health` 检查状态。

## 启动前端

在 `web` 目录安装依赖：

```powershell
cd web
npm install
```

开发时可按需运行：

```powershell
npm run dev        # 公开站点：http://localhost:3000
npm run dev:admin  # 管理站点：http://localhost:3001/admin
npm run dev:all    # 同时启动以上两个服务
```

Vite 会把 `/api` 和 `/uploads` 代理到公开服务 `8080`，并把 `/admin-api` 代理到管理服务 `8081`。

## 构建一体化站点

```powershell
cd web
npm run build
cd ..
go run .
```

`web/dist` 存在时，Go 服务会托管构建后的资源，并为前端路由回退至 `index.html`。此时可访问：

- `http://localhost:8080`
- `http://localhost:8081/admin`

## 管理员与上传

要在启动时创建管理员，请在配置中设置 `admin_bootstrap.enabled: true`，并提供用户名、密码与昵称。创建成功后请关闭该选项。

上传功能要求同时配置 `upload.dir` 与 `upload.url`。支持 JPEG、PNG 和 WebP 图片，默认单文件最大 5 MiB。具体接口见 [API 文档](api.md)。

## 验证

```powershell
go test ./...
cd web
npm run lint
npm run build
```
