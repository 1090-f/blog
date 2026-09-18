# CI/CD 自动化部署（从提交到线上）

面向第一次接手这个仓库的人：不解释"为什么这么设计"，只讲**每一步具体做什么、怎么执行、怎么验证**。

| 想了解 | 看哪篇 |
|---|---|
| **怎么发版、每一步发生了什么**（本文） | `docs/deployment-pipeline.md` |
| 为什么用 ACR 而不是 GHCR、踩过哪些坑、自测题 | `docs/cicd.md` |
| 不用 CI 时的手工部署、环境变量清单 | `docs/deployment.md`、`docs/configuration.md` |
| 服务器账户与密钥体系（含真实 IP，**不提交仓库**） | `docs/server-ssh-keys.md` |

---

## 0. 一句话概括

**打一个 `v*.*.*` 的 tag，线上就是那个提交的版本。**

```bash
git tag v0.1.6 && git push origin v0.1.6
```

```
                    GitHub Actions runner
   ┌──────────────────────────────────────────────────────────┐
   │ job 1: build-and-push                                     │
   │   checkout → buildx → docker login → 构建并推送镜像        │
   │   产物 = <ACR>/fddz/blog:<commit sha 前 12 位>             │
   └───────────────────────────┬──────────────────────────────┘
                               │  SSH：只传「一条命令 + 那个 12 位 sha」
                               ▼
                            服务器 /srv/blog
   ┌──────────────────────────────────────────────────────────┐
   │ job 2: deploy                                             │
   │   git pull → bash scripts/deploy.sh <sha12>               │
   │                                                           │
   │   deploy.sh: 拉镜像 → 重建 app 容器 → 轮询 /health         │
   │      通过 → 写 .deploy_state = <sha12>，退出 0            │
   │      失败 → 用 .deploy_state 里的旧 sha 重建，退出 1       │
   └──────────────────────────────────────────────────────────┘
```

三个必须记住的点：

1. **镜像从不走 SSH。** CI 把镜像推到 ACR，服务器自己去 ACR 拉。SSH 上真正传的业务数据只有那 12 位 sha。
2. **服务器上不编译任何东西。** 没有 Go、没有 Node、没有源码也能跑；服务器上那份仓库只是为了拿 `scripts/deploy.sh` 和 `docker/compose.prod.yaml`。
3. **镜像 tag 用 commit sha，不用 `latest`。** 因为"上一个版本"必须是一个确定存在的东西，回滚才有落点。

### 链路里的四条通道

| 通道 | 方向 | 协议 | 传什么 |
|---|---|---|---|
| ① 推送镜像 | runner → 阿里云 ACR | HTTPS `docker push` | 镜像（约 30 MB） |
| ② 部署指令 | runner → 服务器 | **SSH** | 一条命令 + 12 位 sha |
| ③ 更新部署脚本 | 服务器 → github.com | HTTPS `git pull` | `scripts/`、`docker/` 等文件 |
| ④ 拉取镜像 | 服务器 → 阿里云 ACR | HTTPS `compose pull` | 镜像 |

---

## 1. 参与这条链路的文件

| 文件 | 角色 | 什么时候被用到 |
|---|---|---|
| `.github/workflows/release.yml` | 发布流水线本体：构建 + 部署两个 job | tag 推送时 |
| `.github/workflows/ci.yml` | 日常 CI：`go test ./...`、`npm run lint/build` | 每次 push / PR |
| `docker/Dockerfile` | 三阶段构建：前端 `npm run build` → 后端 `go build` → alpine 运行 | job 1 在 runner 上执行 |
| `docker/compose.prod.yaml` | 生产编排：app 用**镜像**（不 build）、端口只绑 `127.0.0.1`、带 healthcheck 与资源上限 | 服务器上执行 |
| `scripts/deploy.sh` | 服务器侧部署脚本：拉镜像 → 重建容器 → 健康检查 → 失败回滚 | 由 SSH 调用 |
| `docker/.env` | 服务器本地配置（数据库密码、JWT 密钥等），权限 600，**不进 git** | 容器创建时注入环境变量 |
| `.deploy_state` | 服务器上的单行文件，记录**当前线上真的跑起来的** sha | 每次部署成功/回滚时读 |

---

## 2. 一次性准备（整套只做一次）

### 2.1 服务器

```bash
# 1) Docker（自带 compose 插件）
curl -fsSL https://get.docker.com | sh

# 2) 专用部署账户，不要让 CI 用 root
sudo adduser --disabled-password --gecos "" deploy
sudo usermod -aG docker deploy

# 3) 放一份仓库到 /srv/blog（服务器只用它的 docker/ 和 scripts/）
sudo mkdir -p /srv && sudo chown deploy:deploy /srv
sudo -u deploy git clone https://github.com/1090-f/blog.git /srv/blog

# 4) 填配置
cd /srv/blog
cp docker/.env.example docker/.env
chmod 600 docker/.env
openssl rand -base64 32      # 生成 BLOG_JWT_SECRET
# 编辑 docker/.env：替换全部示例密码、填入 JWT 密钥
```

### 2.2 部署专用密钥

在**本地**生成一对，公钥装到服务器，私钥交给 GitHub Secret：

```bash
ssh-keygen -t ed25519 -C "github-actions-deploy" -f ./deploy_key -N ""
ssh-copy-id -i ./deploy_key.pub deploy@<服务器IP>
cat ./deploy_key            # 全文粘到 GitHub Secret DEPLOY_SSH_KEY
```

> Secret 是单向的：存进去读不回来。所以**本地这份私钥是唯一可取的副本**，删之前先确认流水线跑通。

### 2.3 镜像仓库（阿里云 ACR 个人版）

1. 创建实例 → 记下实例域名（形如 `crpi-xxxxxxxxxxxxxxxx.<地域>.personal.cr.aliyuncs.com`，**从控制台复制，别手打**）
2. 创建命名空间（如 `fddz`）→ 创建镜像仓库（如 `blog`），类型**显式选「公开」**
3. 「访问凭证」页设置**固定密码**（不要用 1 小时有效期的临时密码），记下登录名 = 阿里云账号名

需要两个值：`<registry域名>/<命名空间>/<仓库名>`。

### 2.4 GitHub 仓库配置

Settings → Secrets and variables → Actions：

| 名称 | 类型 | 值 |
|---|---|---|
| `DEPLOY_HOST` | Secret | 服务器 IP 或域名 |
| `DEPLOY_USER` | Secret | `deploy` |
| `DEPLOY_SSH_KEY` | Secret | 私钥全文（含 BEGIN/END 两行） |
| `ACR_USERNAME` | Secret | 阿里云账号名 |
| `ACR_PASSWORD` | Secret | ACR **固定密码** |
| `DEPLOY_ENABLED` | Variable | `true`；不配则只构建、不部署 |
| `DEPLOY_PORT` | Variable | `22` |
| `DEPLOY_PATH` | Variable | `/srv/blog` |
| `SITE_URL` | Variable | 可选。填了才会在部署后从公网复验 `/health` |

**为什么分 Secret 和 Variable**：GitHub 不允许在 `if` 里读 secrets，所以"是否部署"这个开关必须用 Variable 表达。

只想练构建、还不想碰服务器 → `DEPLOY_ENABLED` 留空，打 tag 后 job 1 照跑，job 2 自动跳过。

---

## 3. 一次发布的完整过程

### 3.1 触发

```bash
git add -A && git commit -m "feat: xxx" && git push     # 只触发 ci.yml
git tag v0.1.6 && git push origin v0.1.6                # 触发 release.yml
```

`release.yml` 的触发条件：tag 匹配 `v*.*.*`，或在 Actions 页面手动 `workflow_dispatch`（可勾选是否部署）。权限只有 `contents: read` —— 推镜像是 ACR 凭据在管，不需要 `packages: write`。

### 3.2 job 1：build-and-push（在 runner 上）

| 步骤 | 做什么 | 失败常见原因 |
|---|---|---|
| Checkout | 拉代码 | — |
| 计算镜像名与 tag | 小写化镜像名；生成 `<image>:<sha12>`，如果是 tag 触发再额外加一个 `<image>:v0.1.6` | 镜像名含大写字母 |
| 安装 Buildx | 准备构建器 | — |
| 登录阿里云 ACR | 用 `ACR_USERNAME` / `ACR_PASSWORD` | 密码用了临时密码（已过期）；登录名填了邮箱/AccessKey |
| 构建并推送镜像 | 按 `docker/Dockerfile` 构建，推送到 ACR | 见下方两条注意事项 |

产物：`${{ steps.meta.outputs.sha_tag }}`，即 `GITHUB_SHA` 前 12 位，作为 job 输出交给 job 2。

**两条容易翻车的配置（已写在 workflow 里，改的时候别删）**：

- `provenance: false` + `sbom: false` —— buildx 默认会附带证明清单（以 `application/vnd.oci.empty.v1+json` 为 config、`platform=unknown/unknown` 的额外 manifest）。**GHCR 接受，ACR 直接拒绝整个 push**，报 `denied: unknown manifest class ...`。
- `cache-from/cache-to: type=gha` —— 复用构建层缓存。注意默认 scope 是 git ref，**跨 tag 不命中**，所以要么别假设"只改了 workflow 层就不变"，要么接受每次全量重拉镜像。

### 3.3 job 2：deploy（SSH 到服务器）

`needs: build-and-push` 之后才跑，条件是 `DEPLOY_ENABLED == 'true'`。实际执行的就是一段 shell：

```bash
set -e
cd "$DEPLOY_PATH"                 # /srv/blog
git fetch --tags --prune
git checkout main
git pull --ff-only                # 让部署脚本自己保持最新
bash scripts/deploy.sh "<sha12>"  # ← 唯一的"版本号"来源
```

关键参数 `command_timeout: 15m`：`appleboy/ssh-action` 默认只有 10 分钟。超时的症状极具欺骗性——SSH 会话被杀 → 远端脚本收 SIGHUP 静默中断 → 服务器上既没有 `.deploy_state` 也没有容器，**看起来像"部署压根没执行"**，实际是拉到一半被掐。

### 3.4 `scripts/deploy.sh` 做了什么

```
读取 .deploy_state（= PREV_TAG，回滚落点）
        │
        ▼
compose pull app                 ← 按新 sha 拉镜像；这一步失败就直接退出，不回滚
        │
        ▼
compose up -d app                ← 只重建 app，数据库容器不动
        │
        ▼
wait_health()：最多 30 次 × 2 秒 = 60 秒
        │
   ┌────┴─────┐
 通过        失败
   │          │
   │          ├─ compose logs --tail=80 app（留现场）
   │          ├─ 用 PREV_TAG 再 up -d app
   │          ├─ wait_health() 复查
   │          └─ exit 1（即使回滚成功也退出 1）
   │
   └─ 写 .deploy_state = NEW_TAG，exit 0
```

**为什么只有健康检查通过才写 `.deploy_state`**：这样"当前版本"永远等于"最后一个真的跑起来了的版本"，回滚目标才可信。

**退出码语义**：`exit 0` = 发布成功；`exit 1` = 本次发布失败（Actions 上变红）。**回滚成功也算失败**——新版本确实是坏的，不该显示绿色。

### 3.5 健康检查是怎么判定的

两处，配合使用：

| 层级 | 配置 | 作用 |
|---|---|---|
| 容器级 | `compose.prod.yaml` 的 `healthcheck`：容器内 `wget http://127.0.0.1:8080/health`，`interval 10s` / `retries 3` / `start_period 15s` | Docker 标记 `healthy` / `unhealthy`，也是 `db` 的 `depends_on` 依据 |
| 部署级 | `deploy.sh` 的 `wait_health`：宿主机 `curl http://127.0.0.1:8080/health`，30 次 × 2 秒 | 决定这次发布成功还是失败 |

容器"起来了"不等于发布成功——必须应用真的能响应 `/health`。

---

## 4. 实测耗时基线（用来判断"这次是不是不正常"）

| 版本 | 总耗时 | build-and-push | deploy | 结果 |
|---|---|---|---|---|
| `v0.1.4`（正常发布） | 1m 40s | 1m 17s | **16s** | success |
| `v0.1.5`（故意坏版本） | 3m 39s | 2m 7s | **1m 26s** | failure + 自动回滚成功 |
| `v0.1.2`（迁 ACR 之前，走 ghcr.io） | 9m 30s | — | 约 7m | success，差点撞 10 分钟超时 |

镜像仓库速度对比（同一台服务器、同一份镜像）：

| 仓库 | 实测速度 |
|---|---|
| 阿里云 ACR（杭州，服务器在武汉） | 单流 **4.5 MB/s**，容器 16 秒完成部署 |
| `ghcr.io` 直连 | **40~50 KB/s**，拉 6~8 分钟 |

> Docker 的 `registry-mirrors` 只对 `docker.io` 生效，**给 ghcr.io 配加速器无效**——这是换 ACR 的根本原因。

---

## 5. 日常操作清单

```bash
# —— 发版 ——
git tag v0.1.6 && git push origin v0.1.6
# 然后在 GitHub → Actions → Release 看两个 job

# —— 看服务器当前状态 ——
cd /srv/blog && bash scripts/deploy.sh --status

# —— 手动回滚到任意历史版本 ——
docker images | grep crpi-                 # 先看本地有哪些版本
bash scripts/deploy.sh <上一个 sha12>

# —— 看日志 ——
docker compose --env-file docker/.env -f docker/compose.prod.yaml ps
docker compose --env-file docker/.env -f docker/compose.prod.yaml logs -f app

# —— 直接问服务 ——
curl -fsS http://127.0.0.1:8080/health
```

---

## 6. 验收：怎么确认这次发布真的成功了

部署绿了不等于版本换对了，按顺序确认四个点：

| # | 命令 / 观察点 | 期望 |
|---|---|---|
| 1 | Actions 上两个 job 都绿 | `build-and-push` success、`deploy` success |
| 2 | `cat /srv/blog/.deploy_state` | 等于这次推送的 sha 前 12 位 |
| 3 | `docker ps --format '{{.Names}} {{.Image}} {{.Status}}'` | `docker-app-1` 的镜像 tag = 该 sha，状态 `healthy` |
| 4 | `curl -fsS http://127.0.0.1:8080/health` | `{"code":0,...}` |

**没有 job 日志时，怎么证明"自动回滚真的发生了"？** 看两个时间戳就够了：

| 观察 | 发布成功 | 失败后自动回滚 |
|---|---|---|
| `.deploy_state` | 内容 + mtime 都被改写 | **两者都不变** |
| app 容器 | 镜像 = 新 tag，`Created` 在本次窗口内 | `Created` 也被刷新，但镜像是**旧 tag** |
| 本地镜像列表 | 出现新 tag | 出现**坏 tag**（证明 `pull` 用的是新 tag） |

判据：`.deploy_state` 没动 + 容器镜像是旧 tag + 容器 `Created` 落在部署窗口内 ⇒ 只可能是"先起坏版本、健康检查失败、再重建回旧版本"。

---

## 7. 排错速查

| 现象 | 可能原因 | 处理 |
|---|---|---|
| job 1 构建失败 | 代码或 Dockerfile 问题 | 本地复现：`docker build -f docker/Dockerfile .` |
| push 报 `denied: unknown manifest class` | 构建时又带上了 provenance/SBOM 清单 | 确认 `provenance: false` + `sbom: false` 还在 |
| 登录 ACR 失败 | 用了临时密码 / 登录名不对 | 重置固定密码，登录名用阿里云账号名 |
| job 2 被跳过 | `DEPLOY_ENABLED` 不是 `'true'` | 检查 repository variable |
| SSH 步骤失败 | 私钥、端口、安全组 | 本地用同一把私钥手动 `ssh` 一次 |
| SSH 步骤"跑了一半就没了" | `command_timeout` 不够 | 看步骤耗时是否接近 15 分钟 |
| `compose pull` 报 401/denied | ACR 仓库是私有的，服务器没登录 | 控制台把仓库类型改成「公开」，或给服务器配 docker login |
| 容器反复重启 | 环境变量缺失、连不上数据库 | `docker compose ... logs app` |
| 健康检查一直不过 | 应用启动慢 / 依赖没就绪 | 调大 `start_period` 与 `wait_health` 次数 |
| 部署绿了但页面打不开 | 反向代理没配 | 反向代理 upstream 指向 `127.0.0.1:8080` |

---

## 8. 设计取舍与已知边界

**为什么这样设计**

- **镜像 tag 用 commit sha** —— 可追溯、可回滚；`latest` 会被覆盖，等于没有"上一个版本"。
- **构建与运行分离** —— 服务器只装 Docker，攻击面和运维成本都小；换机器时不用装 Go/Node 工具链。
- **健康检查后才认账** —— "发布成功"的定义是应用真的能响应，而不是容器启动了。
- **端口只绑 `127.0.0.1`** —— 应用与数据库都不直接暴露公网，由同机反向代理提供 HTTPS。

**回滚的边界（要理解，别期待它万能）**

`deploy.sh` 只覆盖"镜像拉到了，但容器起不来 / 服务不健康"这一类失败。以下情况它救不了：

- 镜像根本没拉到（tag 写错、registry 拒绝）→ 脚本在 `pull` 阶段就退出，**线上没被动过，不需要回滚**，这是对的
- 容器起来了、健康检查也过了，但**业务逻辑是错的** → 只能靠监控、灰度与人工判断

**还没做的部分**

- 反向代理与 HTTPS 尚未配置，目前对外只到 `127.0.0.1:8080`
- 数据库迁移脚本（`migrations/init.sql`）没有接入发布流程
- 没有数据库 / 上传目录的定时备份
- 单实例重建，发布时有秒级空窗，没有滚动发布
- 失败目前只能在 Actions 页面看到，没有主动告警

**两点运维纪律**

- 永远不要 `docker compose down -v`，`-v` 会连数据卷（`mysql_data`、`app_uploads`）一起删掉。
- 首次部署为了建管理员会临时把 `BLOG_ADMIN_BOOTSTRAP_ENABLED` 设为 `true`，**建完必须改回 `false` 并重新部署一次**。
