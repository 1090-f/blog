# 从 GitHub 到服务器的 CI/CD 全流程

这篇是学习手册：链路长什么样、要准备什么、怎么完整走一遍、坏了怎么查。

## 一、这条链路在做什么

一句话：**把「在我电脑上能跑」变成「打一个 tag，线上就是那个版本」。**

```
你打 tag  →  GitHub 构建镜像  →  镜像进仓库  →  SSH 让服务器换版本  →  健康检查通过才认账
   ↑                                                                        ↓
   └──────────────────── 失败就切回上一个 tag，服务不停 ───────────────────┘
```

关键设计只有三条：

1. **构建和运行彻底分开**。编译发生在 GitHub 的机器上，服务器上只有 Docker，没有 Go、没有 Node、没有源码也可以。
2. **镜像仓库是唯一的中转站**。GitHub 把产物推上去，服务器从那里拉。两边不直接传文件。
3. **镜像用 commit sha 做 tag**。同一个提交永远对应同一个镜像，所以「上一个版本」是确定存在的，回滚才有意义。

### 三条网络通路，只有一条是 SSH

这条链路里其实跑着**三条互不相干**的通道，别把它们的职责混在一起：

| 通道 | 谁连谁 | 走什么协议 | 传的是什么 |
|---|---|---|---|
| ① 构建产物 | Actions runner → 阿里云 ACR | HTTPS（docker push） | **镜像**（几十 MB） |
| ② 部署指令 | Actions runner → 服务器 | **SSH** | **一条命令字符串**：`bash scripts/deploy.sh <sha>` |
| ③ 取代码 | 服务器 → github.com | HTTPS（git pull） | 仓库文件（主要是 `scripts/` 和 `docker/`） |
| ④ 取镜像 | 服务器 → 阿里云 ACR | HTTPS（docker compose pull） | **镜像** |

**关键：镜像从来不走 SSH。** CI 构建完镜像推到镜像仓库就没事了；服务器是**自己去**那里把镜像拉下来的。
SSH 那条通道上真正"从 GitHub 传到服务器"的业务数据，只有那个 **12 位的 commit sha** —— 也就是"这次该上哪个版本"这句话。

这也解释了为什么服务器上要 clone 一份仓库：**不是为了源码**（源码早就编译进镜像了），是为了拿到
`scripts/deploy.sh` 和 `docker/compose.prod.yaml` 这两个**部署用的脚本和编排文件**。
`git pull` 只是为了让部署脚本自己保持最新（所以 release.yml 里是先 pull 再执行）。

再顺一层：**push 镜像 ≠ 服务更新**。镜像躺在仓库里，容器不会自己动。
必须有人在服务器上执行 `docker compose pull app && up -d app`，服务才会真的换成新版本 —— 这就是 SSH 那一环存在的唯一理由。

## 二、一次性准备

### 2.1 服务器侧

```bash
# 1) 装 Docker（自带 compose 插件）
curl -fsSL https://get.docker.com | sh

# 2) 建一个专用部署用户，不要用 root 跑部署
sudo adduser --disabled-password --gecos "" deploy
sudo usermod -aG docker deploy

# 3) 放一份仓库。服务器上只用它的 docker/ 和 scripts/，不在这里编译
sudo mkdir -p /srv && sudo chown deploy:deploy /srv
sudo -u deploy git clone https://github.com/1090-f/blog.git /srv/blog

# 4) 配置密钥（唯一需要人工填的部分）
cd /srv/blog
cp docker/.env.example docker/.env
chmod 600 docker/.env
openssl rand -base64 32        # 生成 BLOG_JWT_SECRET
# 编辑 docker/.env：换掉所有示例密码、填入 JWT 密钥
# 并把 BLOG_ADMIN_BOOTSTRAP_ENABLED 改成 false
```

**在本地机器上**生成一对部署专用密钥（不是服务器上）：

```bash
# 实际落地位置是 F:\blog\.deploy\（已被 .gitignore 忽略），不在 ~/.ssh 下
ssh-keygen -t ed25519 -C "github-actions-deploy" -f /f/blog/.deploy/blog_deploy -N ""
ssh-copy-id -i /f/blog/.deploy/blog_deploy.pub deploy@<服务器IP>
cat /f/blog/.deploy/blog_deploy     # 私钥全文，粘贴到 GitHub Secret DEPLOY_SSH_KEY
```

> GitHub Secret 是**单向**的：存进去之后在 UI/API 都读不回来，只能覆盖或删除。
> 所以本地 `.deploy\blog_deploy` 是**唯一能取回的副本**，删它之前要确认 CI 已经跑通了。

#### 镜像仓库的地址怎么读

`crpi-lxhi919p6atrn34g.cn-hangzhou.personal.cr.aliyuncs.com/fddz/blog:a1b2c3d4e5f6` 拆开看：

| 段 | 值 | 含义 |
|---|---|---|
| ① registry 域名 | `crpi-lxhi919p6atrn34g.cn-hangzhou.personal.cr.aliyuncs.com` | ACR 个人版实例的专属域名。2024-09-09 之后新建的实例都是 `crpi-<实例ID>.<地域>.personal.cr.aliyuncs.com` 这个格式 |
| ② 命名空间 | `fddz` | 一个命名空间下可以有多个镜像仓库 |
| ③ 仓库名 | `blog` | 就是 Docker 语境里的「镜像仓库」，一个仓库里有多个版本 |
| ④ 标签 | `:a1b2c3d4e5f6` | 版本（commit sha 前 12 位）；打版本 tag 时还会额外多推一个 `:v1.0.0` |

**域名一定要从控制台「访问凭证」页复制，别手打。** 实例 ID 是 16 位随机串，抄错一位的后果是
DNS 直接解析不了（`curl` 返回 `HTTP 000`、`connect=0s`），而报错里**完全不会提「域名写错了」**，
很容易往网络故障方向怀疑 —— 这个坑本项目的迁移过程里真的踩过一次。

#### 为什么用 ACR 而不是 GHCR

纯网络现实逼出来的，不是技术偏好：

| | 直连 `ghcr.io` | 阿里云 ACR |
|---|---|---|
| 实测拉取速度 | **40~50 KB/s** | **25 MB/s** |
| 17MB 应用镜像耗时 | 6~8 分钟 | 1 秒内 |

服务器在武汉（`cn-wuhan-lr`），镜像放杭州 —— 跨地域，但同在国内骨干网上，差两个数量级。

**这不是配置问题。** Docker 的 `registry-mirrors` 在客户端里硬编码**只对 `docker.io` 生效**，
`docker pull ghcr.io/...` 压根不会去查镜像加速器列表。所以给 `daemon.json` 加再多加速器也救不了 ghcr
—— 实测同一台机器走加速器拉 799MB 的 `mysql:8.0` 几分钟就完事，却对 17MB 的 blog 镜像束手无策。

顺带记一笔代价：ACR 个人版**上传下载都免费**，限额 3 个命名空间 / 300 个公开仓库；
但它是共享带宽、无 SLA，官方标注「仅限开发测试使用」。对一个个人博客完全够。

#### 推送要登录，拉取不用

阿里云官方文档原话：

> 无论是公开还是私有类型仓库，**推送镜像都需要先进行登录**。
> 公开：拉取镜像时可以免登录，直接通过网络拉取。

所以链路两端的需求是不一样的：

| 方向 | 谁发起 | 要凭据吗 | 用什么 |
|---|---|---|---|
| 推送 | Actions runner | ✅ 要 | **固定密码**，经 `docker/login-action` 注入 |
| 拉取 | 服务器 | ❌ 不要 | 仓库类型选「公开」，匿名直接拉 |

三处容易搞错的地方：

- **必须用「固定密码」，不能用「临时密码」。** 临时密码有效期只有 1 小时，放进 CI 等于每次流水线
  都得先换一次口令。固定密码在「访问凭证」页点「设置固定密码」生成。
  ⚠️ 它**存进去之后再也读不回来**，只能重置 —— 生成后立刻存进密码管理器。
- **登录名是「阿里云账号名」**，不是邮箱也不是 AccessKey。
  控制台「访问凭证」页会直接给一条填好用户名的 `docker login` 命令，照着抄就行。
  ⚠️ 这个账号名不要写进任何会提交到 public 仓库的文档里（本文件此前就写过，已移除）——
  它是阿里云控制台的登录凭据之一。
- **留意「自动创建仓库」的默认类型。** 如果命名空间开着自动创建、且默认配置是「私有」，
  push 一个不存在的仓库名会自动建出一个**私有**仓库 —— 服务器随后匿名拉取就会失败。
  **稳妥做法是手工创建仓库、类型显式选「公开」**，别依赖自动创建。

#### 怎么验证服务器真的能匿名拉

**先记住一个反直觉的点：不能用「不带 token 请求是不是返回 401」来判断可见性 —— 这个方法不成立。**

`401` 是 Docker Registry v2 的**标准挑战**：任何不带 Bearer token 的请求都会收到它，
**公开镜像也一样**。仓库是在用 401 告诉客户端「去我 `WWW-Authenticate` 头里那个地址换张票」。

正确做法是走完整的换票流程，看**带上 token 之后**的状态码 —— token 端点由仓库在 401 里告诉你，
不要自己猜：

```bash
REG=crpi-lxhi919p6atrn34g.cn-hangzhou.personal.cr.aliyuncs.com

# 1) 问出 token 端点
curl -sI "https://$REG/v2/" | grep -i www-authenticate
# ACR 杭州的实际返回：
#   Bearer realm="https://dockerauth.cn-hangzhou.aliyuncs.com/auth",
#          service="registry.aliyuncs.com:cn-hangzhou:26842"

# 2) 换一张匿名 pull token
AUTH=https://dockerauth.cn-hangzhou.aliyuncs.com/auth
SVC=registry.aliyuncs.com:cn-hangzhou:26842
TOK=$(curl -s "$AUTH?scope=repository:fddz/blog:pull&service=$SVC" \
      | sed -E 's/.*"token":"([^"]+)".*/\1/')

# 3) 带 token 列 tag —— 这个状态码才说明问题
curl -s -o /dev/null -w "%{http_code}\n" -H "Authorization: Bearer $TOK" \
     "https://$REG/v2/fddz/blog/tags/list"
# 200 → 公开，匿名可拉
# 401/403 → 私有或不存在，服务器必须 docker login
```

比 API 更权威的是**直接问 Docker 本身**（在服务器上跑，先确认没有 `~/.docker/config.json`）：

```bash
docker manifest inspect $REG/fddz/blog:<sha12> && echo "匿名可拉"
```

`deploy.sh` 里没有登录步骤，所以**这个前提必须先满足**，否则第 4 步 `compose pull app`
会以 401 失败，而且失败点看起来像"部署脚本坏了"，很容易误判。

### 2.2 GitHub 侧

仓库 → Settings → Secrets and variables → Actions 里配这些：

| 名称 | 类型 | 值 |
|---|---|---|
| `DEPLOY_HOST` | Secret | 服务器 IP 或域名 |
| `DEPLOY_USER` | Secret | `deploy` |
| `DEPLOY_SSH_KEY` | Secret | 上一步 `cat` 出来的私钥全文，含 `BEGIN` / `END` 两行 |
| `ACR_USERNAME` | Secret | 阿里云账号名，ACR「访问凭证」页可查 |
| `ACR_PASSWORD` | Secret | ACR 的**固定密码**（不是阿里云登录密码，也不是 AccessKey） |
| `DEPLOY_ENABLED` | Variable | `true`（不配就只构建、不部署） |
| `DEPLOY_PORT` | Variable | `22` |
| `DEPLOY_PATH` | Variable | `/srv/blog` |
| `SITE_URL` | Variable | 例如 `https://blog.example.com`，暂时没有可留空 |

> 为什么要分 Secret 和 Variable？GitHub 出于安全不允许在 `if` 条件里读 `secrets`，
> 所以「是否执行部署」这个开关必须用 Variable 来表达。
> `ACR_USERNAME` / `ACR_PASSWORD` 只被 `docker/login-action` 读取，不参与任何条件判断，
> 所以放在 Secret 里是安全的。

**只练构建、还不想碰服务器**：把 `DEPLOY_ENABLED` 留空，直接打一个 tag，
第一个 job 会构建并推送镜像，第二个 job 自动跳过。整条链路的 GitHub 半边你就能先跑通。

## 三、完整走一遍

### 演练 A：正常发布

```bash
# 1) 先推代码（只触发测试，不发布）
git add -A
git commit -m "chore: 接入 CI/CD 流水线"
git push

# 2) 打 tag 并推送 —— 这一步才会触发发布
git tag v0.1.0
git push origin v0.1.0

# 3) 打开 GitHub → Actions → Release
#    应该看到 build-and-push 先绿，deploy 后绿

# 4) 到服务器上验证
cd /srv/blog
bash scripts/deploy.sh --status          # 看到当前版本 = 那个 sha
curl -fsS http://127.0.0.1:8080/health   # 应该返回 200
docker compose -f docker/compose.prod.yaml ps
```

第一次成功后，`.deploy_state` 里会记下这个 sha —— 它就是你以后自动回滚的落点。

### 演练 B：手动回滚

想确认「旧版本还能拉起来」这件事是真的，手动切一次：

```bash
docker images | grep crpi-                 # 列出本地已有的版本
bash scripts/deploy.sh <上一个 sha>
```

### 演练 C：故意发一个坏版本，看它自动回滚

这是最值得做的一次演练，因为它验证了整套机制。

```bash
# 在本地临时把 Dockerfile 的启动命令改坏，让容器「起得来但服务不健康」
# 把最后一行的 ENTRYPOINT 改成：
#   ENTRYPOINT ["sleep", "infinity"]
git add docker/Dockerfile
git commit -m "test: 故意造一个坏版本，验证回滚"
git tag v0.1.1 && git push origin v0.1.1
```

然后在服务器上看结果：

```bash
# Actions 里 deploy 这一步会失败变红
# 服务器上会看到脚本探测 30 次、判定失败、自动切回 v0.1.0 的 sha
bash scripts/deploy.sh --status
curl -fsS http://127.0.0.1:8080/health   # 服务依然是好的
```

**验证完记得把 Dockerfile 改回去。** 演练 C 结束后你应该能自己解释：
为什么脚本要在「健康检查通过之后」才写 `.deploy_state`。

### 关于回滚的边界（要理解，别背）

`deploy.sh` 的回滚只覆盖「镜像拉下来了，但容器起不来 / 服务不健康」这一种失败。
如果镜像根本没拉到（tag 写错、registry 拒绝），脚本在 `pull` 阶段就退出了，
不会执行回滚——**这是对的**，因为那时线上还没被动过，没什么可回滚的。
真正的危险场景是另一类：容器起来了、健康检查也过了，但业务逻辑是错的。
这种只能靠灰度、监控和人工判断，自动化回滚救不了你。

## 四、坏了怎么查

| 现象 | 常见原因 | 怎么处理 |
|---|---|---|
| build 步骤失败 | 代码或 Dockerfile 问题 | 本地跑一次 `docker build -f docker/Dockerfile .` 复现 |
| deploy job 被跳过 | `DEPLOY_ENABLED` 不是 `true` | 检查 repository variable |
| SSH 连接失败 | 私钥不对 / 端口不对 / 安全组没放行 | 在本地用同一把私钥手动 `ssh` 一次 |
| `pull` 报 denied / 401 | 仓库是私有的，或根本不存在 | 在 ACR 控制台确认仓库类型是「公开」且名称拼写一致 |
| 容器反复重启 | 环境变量缺失、连不上数据库 | `docker compose -f docker/compose.prod.yaml logs app` |
| 健康检查一直不过 | 应用启动慢，或 `/health` 依赖没就绪 | 调大 `start_period` 和脚本里的等待次数 |
| 部署绿了但页面打不开 | 反向代理没配 | 检查 Nginx / Caddy 的 upstream 是否指向 `127.0.0.1:8080` |

常用命令：

```bash
cd /srv/blog
docker compose --env-file docker/.env -f docker/compose.prod.yaml ps
docker compose --env-file docker/.env -f docker/compose.prod.yaml logs -f app
docker compose --env-file docker/.env -f docker/compose.prod.yaml images
```

## 五、容易踩的坑

1. **ACR 推送必须登录、拉取不用** —— 但这是以「仓库类型选公开」为前提的。
   如果依赖命名空间的「自动创建仓库」、而它的默认配置是私有，push 出来的仓库就是私有的，
   服务器随后匿名拉取会 401，**而报错看起来像部署脚本坏了**。所以手工建仓库、类型显式选公开。
2. **ACR 域名从控制台复制，别手打。** 实例 ID 是 16 位随机串，错一位就是 DNS 解析失败
   （`curl` 返回 `HTTP 000`、`connect=0s`），而报错里**完全不会说「域名写错了」**。
   同一类坑还有：登录名是「阿里云账号名」，不是邮箱、不是 AccessKey。
3. **不要给镜像打 `latest`**。它会被下一次构建覆盖，回滚就失去了目标。
4. **`docker/login-action` 必须用「固定密码」**，临时密码只有 1 小时有效期。
   固定密码存进去就读不回来，生成后立刻存进密码管理器。
5. **不要用「不带 token 是否 401」判断镜像可见性。** 401 是 Registry v2 的标准挑战，
   **公开镜像也返回 401**。要换成 token 之后再看状态码（§2.1 有完整命令）。
6. **卷名前缀跟着 compose 文件所在目录走**。两个 compose 文件都在 `docker/` 下，
   所以它们共用同一对卷（`docker_mysql_data`、`docker_app_uploads`），
   **容器名也因此是 `docker-app-1` / `docker-db-1`**。好处是数据延续，代价是本地和线上的
   数据混在一个项目名下。（想改成 `blog-app-1` 得在 compose 里加 `name: blog`，
   但换了项目名后旧容器不会被复用，会因 8080 端口冲突而启动失败 —— 必须先 `down` 掉旧的。）
7. **永远不要 `down -v`**。`-v` 会连数据卷一起删掉。
8. **国内服务器拉 `ghcr.io` 很慢，本项目已把镜像仓库换成阿里云 ACR**（见 §2.1）。
   这是整条链路里唯一需要按网络环境调整的地方。
9. **服务器上的 `docker/.env` 权限要是 600**，且不能进 git（`.gitignore` 已经覆盖 `.env`）。
10. **首次部署要临时打开 bootstrap 建管理员**。`compose.prod.yaml` 里该项默认 `false`，
    首次部署前在 `.env` 里改成 `true` → 建好管理员 → **改回 `false` 再部署一次**，
    否则每次重启都会去尝试 bootstrap。
11. **Actions 的 `if` 里不能读 secrets**，所以开关用 Variable。
12. **部署类步骤的超时必须显式设置。** `appleboy/ssh-action` 的 `command_timeout` 默认只有 10 分钟；
    超时后 SSH 会话被杀、远端脚本静默中断，服务器上既没有 `.deploy_state` 也没有容器，
    **看起来像「部署根本没执行」**，实际是中途被掐。本项目 v0.1.1 就是这么失败的
    （步骤耗时 10.1 分钟，恰好撞线）。

## 六、学完自测

答得出来，这套流程就算真的掌握了：

1. 为什么服务器上不需要装 Go 和 Node？
2. 为什么镜像 tag 用 commit sha 而不是 `latest`？
3. 新版本健康检查失败时，数据库里的数据发生了什么？为什么应该是这样？
4. 如果要求「发布过程中零中断」，这套方案要改哪三个地方？
5. `.deploy_state` 为什么必须在健康检查通过之后才写？

## 七、还没做的部分（知道边界在哪）

- **反向代理与 HTTPS**：目前对外只到 `127.0.0.1:8080`，需要前面加 Caddy 或 Nginx 提供 443
- **数据库迁移**：`migrations/init.sql` 还没有自动执行的时机，发布流程里也没接迁移步骤
- **备份**：还没有 mysqldump 定时任务和异地存放
- **多副本滚动发布**：现在是单实例重建，有秒级空窗
- **告警**：失败目前只能在 Actions 页面看到，没有主动通知
