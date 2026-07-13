# Gin 博客系统设计方案（当前基线版）

## 1. 文档目的

这份文档描述的是当前仓库里已经落地的真实实现，而不是最初的理想蓝图。  
后续开发、补测试、补文档时，应以这里记录的“当前能力”和“当前边界”为准。

## 2. 项目概览

当前项目是一个前后端同仓库的博客系统：

- 后端：`Go + Gin + GORM + MySQL + Viper + JWT`
- 前端：`Vue 3 + Vite + Pinia + Vue Router + Element Plus`
- 运行方式：后端既可单独提供 API，也可在存在 `web/dist` 时直接托管前端静态资源

系统当前已经覆盖以下核心能力：

- 用户注册、登录、JWT 鉴权、个人信息获取
- 文章发布、编辑、删除、列表、详情、热门、最新
- 分类管理
- 评论发布、评论管理
- 图片上传
- 管理后台仪表盘、文章/分类/评论/用户管理

## 3. 当前架构

后端采用分层结构：

`router -> controller -> service -> dao -> model`

职责划分如下：

- `router`：组装依赖并注册路由
- `controller`：HTTP 参数绑定、响应封装、错误映射
- `service`：业务规则与状态校验
- `dao`：数据库访问
- `model`：数据模型定义
- `dto`：请求与响应结构
- `middleware`：鉴权、管理员校验、日志、恢复
- `pkg`：数据库、JWT、上传、统一响应等基础能力

前端主要目录：

- `web/src/layouts`：前台/后台布局
- `web/src/views/front`：前台页面
- `web/src/views/admin`：后台页面
- `web/src/api`：接口封装
- `web/src/stores`：Pinia 状态管理
- `web/src/router`：路由与权限守卫

## 4. 权限模型

### 4.1 公开访问

无需登录即可访问：

- 分类列表
- 已发布文章列表
- 已发布文章详情
- 文章聚合详情
- 文章评论列表
- 最新文章
- 热门文章

### 4.2 登录用户访问

登录后可访问：

- 当前用户资料
- 发布评论
- 发布文章
- 创建分类
- 上传图片

说明：

- 当前实现允许普通登录用户创建文章、创建分类
- 当前实现也允许普通登录用户发布状态为 `published` 的文章

### 4.3 管理员访问

管理员当前可访问：

- 仪表盘统计
- 文章管理
- 分类管理
- 评论管理
- 用户列表
- 用户启用/禁用

## 5. 数据模型

核心表包括：

- `users`
- `categories`
- `articles`
- `comments`

关键字段说明：

### 5.1 users

- `username`：唯一用户名
- `password`：bcrypt 哈希
- `nickname`：昵称
- `role`：`user` 或 `admin`
- `avatar`：头像地址
- `status`：`1` 启用，`0` 禁用

### 5.2 categories

- `name`：唯一分类名
- `description`：分类描述

### 5.3 articles

- `title`
- `summary`
- `content`
- `cover_image`
- `status`：`draft` 或 `published`
- `view_count`
- `user_id`
- `category_id`

### 5.4 comments

- `article_id`
- `user_id`
- `content`
- `status`

## 6. 接口基线

### 6.1 基础接口

- `GET /health`

### 6.2 公开接口

- `GET /api/categories`
- `GET /api/articles`
- `GET /api/articles/latest`
- `GET /api/articles/popular`
- `GET /api/articles/:id`
- `GET /api/articles/:id/full`
- `GET /api/articles/:id/comments`

### 6.3 认证接口

- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/user/profile`

### 6.4 登录态接口

- `POST /api/comments`
- `POST /api/articles`
- `POST /api/categories`
- `POST /api/upload`

### 6.5 管理员接口

- `GET /api/admin/dashboard`
- `GET /api/admin/articles`
- `POST /api/admin/articles`
- `PUT /api/admin/articles/:id`
- `DELETE /api/admin/articles/:id`
- `POST /api/admin/categories`
- `PUT /api/admin/categories/:id`
- `DELETE /api/admin/categories/:id`
- `GET /api/admin/articles/:id/comments`
- `DELETE /api/admin/comments/:id`
- `GET /api/admin/users`
- `PUT /api/admin/users/:id/status`

## 7. 上传设计

上传接口当前具备以下约束：

- 仅允许：`.jpg`、`.jpeg`、`.png`、`.webp`
- 校验内容类型，防止仅修改扩展名绕过
- 按 `articles/YYYY/MM` 目录存储
- 自动生成唯一文件名
- 返回可访问 URL
- 支持最大体积限制 `upload.max_size_bytes`

仓库默认配置为 `5242880`，即 5 MB。

## 8. 工程化状态

当前已经落地：

- 路由按公开、登录态、管理员分组注册
- 统一响应结构
- JWT 鉴权
- 管理员权限校验
- 上传类型与大小校验
- 自动化测试覆盖 service 与 http 关键路径

当前仍需长期关注：

- 文档与代码保持同步
- 前端剩余页面文案和体验持续清理
- 数据迁移策略目前仍以 `AutoMigrate` 为主，`migrations/init.sql` 作为初始化参考

## 9. 当前阶段结论

项目已经不再属于“从零搭建博客骨架”阶段，而是进入了“在现有系统上做补齐、验证、收口”的阶段。  
因此后续工作应以以下原则推进：

1. 文档先反映真实代码状态。
2. 功能补齐必须配套测试。
3. 体验优化不应和真实能力描述混淆。
4. 新增能力要同时考虑后端接口、前端入口、权限边界和配置说明。
