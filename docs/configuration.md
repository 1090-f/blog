# 配置说明

项目启动时会按下面的优先级加载配置：

1. 环境变量
2. `BLOG_CONFIG_FILE` 指定的 YAML 文件
3. 默认的 `config/config.yaml`

这意味着本地开发仍然可以继续使用 `config/config.yaml`，而测试、CI 和生产环境可以直接通过环境变量注入配置。

## 关键环境变量

- `BLOG_SERVER_PORT`
- `BLOG_ADMIN_SERVER_PORT`
- `BLOG_DATABASE_HOST`
- `BLOG_DATABASE_PORT`
- `BLOG_DATABASE_USER`
- `BLOG_DATABASE_PASS`
- `BLOG_DATABASE_NAME`
- `BLOG_JWT_SECRET`
- `BLOG_JWT_EXPIRE`
- `BLOG_UPLOAD_DIR`
- `BLOG_UPLOAD_URL`
- `BLOG_UPLOAD_MAX_SIZE_BYTES`
- `BLOG_ADMIN_BOOTSTRAP_ENABLED`
- `BLOG_ADMIN_BOOTSTRAP_USERNAME`
- `BLOG_ADMIN_BOOTSTRAP_PASSWORD`
- `BLOG_ADMIN_BOOTSTRAP_NICKNAME`

环境变量名称由配置路径转换而来：配置键中的 `.` 会替换为 `_`，并统一使用 `BLOG_` 前缀。例如 `upload.max_size_bytes` 对应 `BLOG_UPLOAD_MAX_SIZE_BYTES`。

## 配置字段

| 字段 | 说明 | 默认值 |
| --- | --- | --- |
| `server.port` | 公开服务端口 | `8080` |
| `admin_server.port` | 管理服务端口 | `8081` |
| `database.*` | MySQL 连接信息 | 无 |
| `jwt.secret` | JWT 签名密钥，必须配置 | 无 |
| `jwt.expire` | Token 有效期，单位为秒 | `7200` |
| `upload.dir` | 上传文件的本地存储目录 | 无 |
| `upload.url` | 上传文件的 URL 前缀 | 无 |
| `upload.max_size_bytes` | 单文件上传上限 | `5242880`（5 MiB） |
| `admin_bootstrap.*` | 首次启动时创建管理员的配置 | `enabled: false` |

`server.port` 与 `admin_server.port` 必须同时存在且不能相同。上传接口和静态上传目录仅在 `upload.dir`、`upload.url` 都非空时启用。

## 自定义配置文件

如果不想使用默认的 `config/config.yaml`，可以在启动前设置：

```powershell
$env:BLOG_CONFIG_FILE="C:\path\to\blog\config\config.yaml"
go run .
```

如果同时提供了配置文件和环境变量，环境变量会覆盖配置文件中的同名字段。

## 最小生产配置示例

```powershell
$env:BLOG_SERVER_PORT="8080"
$env:BLOG_ADMIN_SERVER_PORT="8081"
$env:BLOG_DATABASE_HOST="127.0.0.1"
$env:BLOG_DATABASE_PORT="3306"
$env:BLOG_DATABASE_USER="blog"
$env:BLOG_DATABASE_PASS="replace-me"
$env:BLOG_DATABASE_NAME="blog"
$env:BLOG_JWT_SECRET="replace-with-a-random-secret"
$env:BLOG_UPLOAD_DIR="./uploads"
$env:BLOG_UPLOAD_URL="/uploads"
go run .
```

## 默认值

以下字段在未显式提供时会使用默认值：

- `server.port`: `8080`
- `admin_server.port`: `8081`
- `jwt.expire`: `7200`
- `upload.max_size_bytes`: `5242880`

## 安全建议

- 生产环境通过部署平台的密钥管理或环境变量提供 `BLOG_JWT_SECRET` 和数据库密码，不要将真实密钥提交到仓库。
- 使用足够长的随机 JWT 密钥，并在泄露后及时轮换。
- 管理员首次创建完成后，将 `admin_bootstrap.enabled` 设为 `false`。
- 确保 `upload.dir` 对运行服务的用户可写，同时不允许上传目录覆盖应用源码。
