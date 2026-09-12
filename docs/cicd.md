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
| ① 构建产物 | Actions runner → ghcr.io | HTTPS（docker push） | **镜像**（几百 MB） |
| ② 部署指令 | Actions runner → 服务器 | **SSH** | **一条命令字符串**：`bash scripts/deploy.sh <sha>` |
| ③ 取代码 | 服务器 → github.com | HTTPS（git pull） | 仓库文件（主要是 `scripts/` 和 `docker/`） |
| ④ 取镜像 | 服务器 → ghcr.io | HTTPS（docker compose pull） | **镜像** |

**关键：镜像从来不走 SSH。** CI 构建完镜像推到 ghcr.io 就没事了；服务器是**自己去** ghcr.io 把镜像拉下来的。
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

#### GHCR 的「包」（package）是什么

`ghcr.io/1090-f/blog:a1b2c3d4e5f6` 拆成四段就清楚了：

| 段 | 值 | 含义 |
|---|---|---|
| ① registry 域名 | `ghcr.io` | GitHub Container Registry —— GitHub Packages 产品线下专存容器镜像的那一支 |
| ② owner | `1090-f` | 归属的账号或组织，**不是仓库名** |
| ③ 包名 | `blog` | **「包」指的就是这一段** —— 等价于 Docker 语境里的「镜像仓库」，一个包里有多个版本 |
| ④ 标签 | `:a1b2c3d4e5f6` | 版本（commit sha 前 12 位）；打版本 tag 时还会额外多推一个 `:v1.0.0` |

包名是怎么来的：workflow 里 `IMAGE_NAME: ${{ github.repository }}` 得到 `1090-f/blog`，
GHCR 取最后一段作为包名，所以是 `ghcr.io/1090-f/blog`。

两点容易搞混的：

- **包挂在账号下，不在仓库的文件树里**。仓库页面翻不到它，要看 `github.com/1090-f?tab=packages`。
  包和仓库之间是**关联**（linked）关系：workflow 用 `GITHUB_TOKEN` 推送时会自动建立这个关联。
  关联带来的只是**权限继承**（谁有仓库写权限谁就能管这个包），**可见性不会跟着仓库走**。
- **默认可见性是 private**。官方文档《Configuring a package's access control and visibility》的原话是
  "When you first publish a package, the default visibility is private and only you can see the package."
  所以别假设"仓库是 public，镜像就是 public" —— 用下面的命令匿名探一次才算数。

> ⚠️ **改为 public 不可逆。** 官方原文："Once you make your package public, you cannot make it private again."
> 博客镜像公开没有实际风险（最终镜像里只有 Go 二进制和前端产物，不含任何 `.env` / 配置），
> 但这个包以后收不回来。如果打算往同一个包里放别的东西，就选私有 + `docker login`。

**让服务器能拉镜像**。二选一：

- 把包改成公开：GitHub 仓库页 → 右侧 Packages → 进入该包 → Package settings → Change visibility（镜像里没有敏感信息时最省事）
- 保持私有：在服务器上登录一次，凭证会存在 `~/.docker/config.json` 里

```bash
echo "<带 read:packages 权限的 Classic PAT>" | docker login ghcr.io -u 1090-f --password-stdin
```

**验证包到底是私有还是公开**（在服务器上跑，故意不带任何凭证）：

```bash
# 把 <sha12> 换成任意一个已存在的 tag
curl -s -o /dev/null -w "%{http_code}\n" \
  -H "Accept: application/vnd.oci.image.manifest.v1+json" \
  https://ghcr.io/v2/1090-f/blog/manifests/<sha12>
# 401 → 私有，必须 docker login
# 200 → 公开，匿名可拉
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
| `DEPLOY_ENABLED` | Variable | `true`（不配就只构建、不部署） |
| `DEPLOY_PORT` | Variable | `22` |
| `DEPLOY_PATH` | Variable | `/srv/blog` |
| `SITE_URL` | Variable | 例如 `https://blog.example.com`，暂时没有可留空 |

> 为什么要分 Secret 和 Variable？GitHub 出于安全不允许在 `if` 条件里读 `secrets`，
> 所以「是否执行部署」这个开关必须用 Variable 来表达。

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
docker images ghcr.io/1090-f/blog     # 列出本地已有的版本
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
| `pull` 报 denied | 包是私有的，服务器没登录 | 改包可见性或 `docker login ghcr.io` |
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

1. **GHCR 包默认私有**，服务器上不登录就拉不动。这是新手卡住最多的地方。
2. **不要给镜像打 `latest`**。它会被下一次构建覆盖，回滚就失去了目标。
3. **别用命令行手动 push 这个包。** 官方文档写得很明确：从命令行 push 的镜像**默认不关联仓库**，
   而 workflow 里的 `GITHUB_TOKEN` **对「已存在但未关联到本仓库」的包没有推送权限**。
   也就是说你本地手滑跑一次 `docker push ghcr.io/1090-f/blog:test`，可能把自己 CI 的推送权限堵掉，
   得先去包设置里手动连接回 `blog` 仓库才能恢复。真要手动推，就另起一个包名（如 `blog-test`）。
4. **卷名前缀跟着 compose 文件所在目录走**。这里两个 compose 文件都在 `docker/` 下，
   所以 `compose.yaml` 和 `compose.prod.yaml` 用的是同一对卷（`docker_mysql_data`、`docker_app_uploads`）。
   好处是数据延续，代价是本地和线上的数据混在一个项目名下，心里要有数。
4. **永远不要 `down -v`**。`-v` 会连数据卷一起删掉。
5. **国内服务器拉 `ghcr.io` 可能很慢**。可选方案：换成腾讯云 TCR / 阿里云 ACR（改 workflow 里的
   registry 和登录步骤），或给 Docker 配置代理。这也是唯一一处需要按网络环境调整的地方。
6. **服务器上的 `docker/.env` 权限要是 600**，且不能进 git（`.gitignore` 已经覆盖 `.env`）。
7. **首次部署前先把管理员建好**。生产编排里 `BLOG_ADMIN_BOOTSTRAP_ENABLED` 被写死为 `false`，
   否则你连后台都进不去。临时开启 → 创建管理员 → 改回 `false` → 再部署一次。
8. **Actions 的 `if` 里不能读 secrets**，所以开关用 Variable。

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
