# 部署检查清单

## 上线前准备

- 准备 MySQL 8.x 数据库和运行账户。
- 通过环境变量或部署平台的密钥管理配置 `BLOG_DATABASE_*`、`BLOG_JWT_SECRET`、`BLOG_UPLOAD_DIR` 与 `BLOG_UPLOAD_URL`。
- 将上传目录置于应用发布目录之外或独立挂载，并确认运行用户具有写入权限。
- 确保 `BLOG_SERVER_PORT` 与 `BLOG_ADMIN_SERVER_PORT` 不相同。
- 不在生产环境启用 `BLOG_ADMIN_BOOTSTRAP_ENABLED`；如需首次创建管理员，完成后立即关闭。
- 在发布前运行 `go test ./...`、`npm run lint` 和 `npm run build`。

完整环境变量参见 [配置说明](configuration.md)。

## 构建与发布

在发布机或 CI 中构建前端：

```powershell
cd web
npm ci
npm run lint
npm run build
cd ..
```

再构建后端程序：

```powershell
go build -o blog.exe .
```

发布时必须保留 `web/dist` 与后端程序的相对目录关系；服务会自动托管该目录中的静态资源。生产环境使用 Linux 时，将输出文件名改为适合目标系统的名称即可。

## Docker Compose 部署

仓库将 Docker 相关文件集中在 [`docker/`](../docker/)：`Dockerfile`、`Dockerfile.dockerignore`、`compose.yaml` 与 `.env.example`。该配置会构建前端与后端，并启动 MySQL 8.4。

### 首次启动

```powershell
Copy-Item docker/.env.example docker/.env
# 编辑 docker/.env：替换所有示例密码，并设置长随机 BLOG_JWT_SECRET
docker compose --env-file docker/.env -f docker/compose.yaml up -d --build
```

默认端口为公开服务 `8080` 和管理服务 `8081`。可通过 `.env` 中的 `PUBLIC_PORT` 和 `ADMIN_PORT` 修改主机端口。

Compose 创建两个命名卷：

- `mysql_data`：MySQL 数据
- `app_uploads`：用户上传的图片

停止或更新服务时使用以下命令；不要附加 `-v`，否则会删除这两个持久化卷。查看服务状态与日志：

```powershell
docker compose --env-file docker/.env -f docker/compose.yaml down
docker compose --env-file docker/.env -f docker/compose.yaml ps
docker compose --env-file docker/.env -f docker/compose.yaml logs -f app
```

### 管理员初始化

首次创建管理员时，可暂时在 `.env` 设置：

```dotenv
BLOG_ADMIN_BOOTSTRAP_ENABLED=true
BLOG_ADMIN_BOOTSTRAP_USERNAME=admin
BLOG_ADMIN_BOOTSTRAP_PASSWORD=replace-with-a-strong-password
BLOG_ADMIN_BOOTSTRAP_NICKNAME=管理员
```

应用成功启动并创建管理员后，立即将 `BLOG_ADMIN_BOOTSTRAP_ENABLED` 改回 `false`，然后使用相同的 Compose 命令执行 `up -d` 重新创建应用容器。

## 启动与验证

启动应用：

```powershell
./blog.exe
```

启动后验证：

- `GET http://<host>:<public-port>/health`
- `GET http://<host>:<admin-port>/health`
- 公开首页、文章详情和上传资源可访问
- 管理员可从 `http://<host>:<admin-port>/admin` 登录并访问管理页
- 未登录和非管理员请求无法访问 `/api/admin/*`

公开服务默认端口为 `8080`，管理服务默认端口为 `8081`。

## 反向代理与运维建议

- 在反向代理中分别将公开域名和管理域名（或管理路径）转发到两个端口；仅对外暴露反向代理的 HTTPS 端口。
- 将日志采集、健康检查、数据库备份和上传文件备份纳入部署流程。
- 定期轮换 JWT 密钥和数据库密码；轮换 JWT 密钥会使现有登录令牌失效。
- 为上传目录设置容量监控和备份策略。

## CI

[GitHub Actions 工作流](../.github/workflows/ci.yml) 会在推送与拉取请求中执行：

- `go test ./...`
- `npm ci`
- `npm run lint`
- `npm run build`
