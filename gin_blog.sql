-- MySQL dump 10.13  Distrib 5.7.44, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: gin_blog
-- ------------------------------------------------------
-- Server version	5.7.44-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `article_tags`
--

DROP TABLE IF EXISTS `article_tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `article_tags` (
  `article_id` bigint(20) unsigned NOT NULL,
  `tag_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`article_id`,`tag_id`),
  KEY `fk_article_tags_tag` (`tag_id`),
  CONSTRAINT `fk_article_tags_article` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_article_tags_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `article_tags`
--

LOCK TABLES `article_tags` WRITE;
/*!40000 ALTER TABLE `article_tags` DISABLE KEYS */;
INSERT INTO `article_tags` VALUES (16,1),(17,2),(18,2),(17,3),(19,4),(18,7),(17,8),(19,8),(20,8),(16,9),(20,9);
/*!40000 ALTER TABLE `article_tags` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `articles`
--

DROP TABLE IF EXISTS `articles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `articles` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(150) COLLATE utf8mb4_unicode_ci NOT NULL,
  `summary` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `content` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `cover_image` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'draft',
  `view_count` bigint(20) NOT NULL DEFAULT '0',
  `user_id` bigint(20) unsigned NOT NULL,
  `category_id` bigint(20) unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_articles_status` (`status`),
  KEY `idx_articles_user_id` (`user_id`),
  KEY `idx_articles_category_id` (`category_id`),
  CONSTRAINT `fk_articles_category` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON UPDATE CASCADE,
  CONSTRAINT `fk_articles_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `articles`
--

LOCK TABLES `articles` WRITE;
/*!40000 ALTER TABLE `articles` DISABLE KEYS */;
INSERT INTO `articles` VALUES (1,'Go 语言入门指南','本文介绍 Go 语言的基础知识，包括为什么选择 Go 以及如何编写第一个程序。','# Go 语言入门指南\n\nGo 是一门简洁而强大的编程语言，由 Google 开发。\n\n## 为什么选择 Go？\n\n- **简洁性**：语法简单，易于学习\n- **高性能**：编译型语言，执行效率高\n- **并发支持**：内置 goroutine 和 channel\n- **跨平台**：支持多种操作系统\n\n## 第一个程序\n\n```go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}\n```\n\n## 总结\n\nGo 语言非常适合构建高性能的网络服务和微服务架构。','','published',0,4,1,'2026-06-24 09:30:18.135','2026-06-24 09:30:18.135'),(2,'使用 Gin 框架构建 Web 应用','学习如何使用 Gin 框架快速构建高性能的 Web 应用。','# 使用 Gin 框架构建 Web 应用\n\nGin 是一个高性能的 Go Web 框架，提供了丰富的功能。\n\n## 安装 Gin\n\n```bash\ngo get -u github.com/gin-gonic/gin\n```\n\n## 基本用法\n\n```go\npackage main\n\nimport \"github.com/gin-gonic/gin\"\n\nfunc main() {\n    r := gin.Default()\n    r.GET(\"/ping\", func(c *gin.Context) {\n        c.JSON(200, gin.H{\n            \"message\": \"pong\",\n        })\n    })\n    r.Run(\":8080\")\n}\n```\n\n## 总结\n\nGin 框架简单易用，适合快速开发 Web 应用。','','published',0,4,1,'2026-06-24 09:30:18.145','2026-06-24 09:30:18.145'),(3,'MySQL 数据库最佳实践','掌握 MySQL 数据库的索引优化、查询优化和安全实践。','# MySQL 数据库最佳实践\n\nMySQL 是最流行的关系型数据库之一。\n\n## 索引优化\n\n- 为常用查询字段创建索引\n- 避免过度索引\n- 使用复合索引提高查询效率\n\n## 查询优化\n\n- 使用 EXPLAIN 分析查询\n- 避免 SELECT *\n- 合理使用 JOIN\n\n## 总结\n\n遵循最佳实践可以提高 MySQL 的性能和安全性。','','published',0,4,1,'2026-06-24 09:30:18.150','2026-06-24 09:30:18.150'),(4,'Docker 容器化部署指南','学习 Docker 的基本概念、编写 Dockerfile 和常用命令。','# Docker 容器化部署指南\n\nDocker 让应用部署变得简单而一致。\n\n## 基本概念\n\n- **镜像**：应用的只读模板\n- **容器**：镜像的运行实例\n- **Dockerfile**：构建镜像的脚本\n\n## 常用命令\n\n```bash\ndocker build -t myapp .\ndocker run -p 8080:8080 myapp\ndocker ps\n```\n\n## 总结\n\nDocker 简化了开发和部署流程，是现代开发的必备技能。','','published',1,4,1,'2026-06-24 09:30:18.154','2026-06-24 09:30:18.154'),(5,'Vue 3 组合式 API 入门','学习 Vue 3 的组合式 API，提升代码复用性。','# Vue 3 组合式 API 入门\n\nVue 3 引入了组合式 API，让代码组织更灵活。\n\n## setup 函数\n\n```vue\n<script setup>\nimport { ref, computed } from \'vue\'\n\nconst count = ref(0)\nconst double = computed(() => count.value * 2)\n</script>\n```\n\n## 总结\n\n组合式 API 让逻辑复用更简单。','','published',2,4,2,'2026-06-24 09:32:16.355','2026-06-24 09:32:16.355'),(6,'React Hooks 最佳实践','掌握 React Hooks 的使用技巧和常见模式。','# React Hooks 最佳实践\n\nHooks 让函数组件拥有了状态和生命周期。\n\n## useState\n\n```jsx\nconst [count, setCount] = useState(0)\n```\n\n## useEffect\n\n```jsx\nuseEffect(() => {\n  document.title = `Count: ${count}`\n}, [count])\n```\n\n## 总结\n\n合理使用 Hooks 可以写出更清晰的代码。','','published',1,4,2,'2026-06-24 09:32:16.361','2026-06-24 09:32:16.361'),(7,'微服务架构设计原则','了解微服务架构的核心设计原则和实践。','# 微服务架构设计原则\n\n微服务让大型系统可以独立开发和部署。\n\n## 核心原则\n\n- 单一职责\n- 服务自治\n- 去中心化治理\n- 容错设计\n\n## 服务拆分\n\n- 按业务领域拆分\n- 保持服务小而专\n\n## 总结\n\n好的微服务设计是成功的关键。','','published',0,4,3,'2026-06-24 09:32:16.365','2026-06-24 09:32:16.365'),(8,'RESTful API 设计规范','设计清晰、一致的 RESTful API。','# RESTful API 设计规范\n\n良好的 API 设计提升开发体验。\n\n## URL 设计\n\n- 使用名词复数: `/users`, `/articles`\n- 层级关系: `/users/:id/articles`\n\n## HTTP 方法\n\n- GET: 查询\n- POST: 创建\n- PUT: 更新\n- DELETE: 删除\n\n## 状态码\n\n- 200: 成功\n- 201: 创建成功\n- 400: 请求错误\n- 404: 未找到\n\n## 总结\n\n遵循规范让 API 更易用。','','published',0,4,3,'2026-06-24 09:32:16.369','2026-06-24 09:32:16.369'),(9,'Nginx 反向代理配置','配置 Nginx 作为反向代理服务器。','# Nginx 反向代理配置\n\nNginx 是最流行的 Web 服务器之一。\n\n## 基本配置\n\n```nginx\nserver {\n    listen 80;\n    server_name example.com;\n    \n    location / {\n        proxy_pass http://localhost:8080;\n        proxy_set_header Host $host;\n        proxy_set_header X-Real-IP $remote_addr;\n    }\n}\n```\n\n## SSL 配置\n\n```nginx\nserver {\n    listen 443 ssl;\n    ssl_certificate /path/to/cert.pem;\n    ssl_certificate_key /path/to/key.pem;\n}\n```\n\n## 总结\n\nNginx 是生产环境的必备工具。','','published',2,4,4,'2026-06-24 09:32:16.374','2026-06-24 09:32:16.374'),(10,'GitHub Actions CI/CD 实践','使用 GitHub Actions 实现自动化部署。','# GitHub Actions CI/CD 实践\n\nGitHub Actions 让 CI/CD 变得简单。\n\n## 基本工作流\n\n```yaml\nname: CI\non: [push]\njobs:\n  build:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: actions/checkout@v2\n      - name: Build\n        run: go build ./...\n      - name: Test\n        run: go test ./...\n```\n\n## 部署步骤\n\n- 构建镜像\n- 推送到 registry\n- 部署到服务器\n\n## 总结\n\n自动化让部署更可靠。','','published',3,4,4,'2026-06-24 09:32:16.379','2026-06-24 09:32:16.379'),(11,'设计模式学习笔记','记录常用设计模式的理解和应用。','# 设计模式学习笔记\n\n设计模式是解决常见问题的经验总结。\n\n## 创建型模式\n\n- 单例模式：确保只有一个实例\n- 工厂模式：封装对象创建过程\n\n## 结构型模式\n\n- 适配器模式：接口转换\n- 装饰器模式：动态添加功能\n\n## 行为型模式\n\n- 观察者模式：事件通知\n- 策略模式：算法选择\n\n## 总结\n\n设计模式让代码更灵活、可维护。','','published',0,4,5,'2026-06-24 09:32:16.384','2026-06-24 09:32:16.384'),(12,'Git 工作流总结','总结团队协作中常用的 Git 工作流。','# Git 工作流总结\n\n好的 Git 工作流提升团队效率。\n\n## Git Flow\n\n- main: 生产分支\n- develop: 开发分支\n- feature/*: 功能分支\n- release/*: 发布分支\n- hotfix/*: 紧急修复\n\n## GitHub Flow\n\n- main: 主分支\n- feature/*: 功能分支\n- 通过 PR 合并\n\n## 常用命令\n\n```bash\ngit checkout -b feature/xxx\ngit add .\ngit commit -m \"feat: xxx\"\ngit push origin feature/xxx\n```\n\n## 总结\n\n选择适合团队的工作流。','','published',5,4,5,'2026-06-24 09:32:16.389','2026-06-24 09:32:16.389'),(13,'快速排序算法详解','本文介绍快速排序算法的原理和实现','# 快速排序算法详解\n\n快速排序是一种高效的排序算法。\n\n## 原理\n\n1. 选择基准元素\n2. 分区操作\n3. 递归排序\n\n## 代码实现\n\n```python\ndef quick_sort(arr):\n    if len(arr) <= 1:\n        return arr\n    pivot = arr[0]\n    left = [x for x in arr[1:] if x < pivot]\n    right = [x for x in arr[1:] if x >= pivot]\n    return quick_sort(left) + [pivot] + quick_sort(right)\n```','','published',16,6,6,'2026-06-24 09:55:51.307','2026-06-24 09:55:51.307'),(14,'英语学习','欢迎加入我们的小组','阿萨法**************``[](url)[](url)[](url)- - [](url)![图片](/uploads/articles/2026/07/20260710204248-ff57577da5b1.png)','/uploads/articles/2026/07/20260710204238-804728825ab7.png','published',3,4,7,'2026-07-10 20:42:50.395','2026-07-10 20:42:50.395'),(15,'Markdown 格式测试文章','测试各种Markdown格式是否正常渲染','## 二级标题测试\n\n这是**加粗文字**测试。\n\n这是*斜体文字*测试。\n\n这是`行内代码`测试。\n\n### 代码块测试\n\n```javascript\nfunction hello() {\n  console.log(\"Hello, World!\");\n}\n```\n\n### 链接测试\n\n访问 [Google](https://www.google.com) 搜索引擎。\n\n### 列表测试\n\n- 列表项 1\n- 列表项 2\n- 列表项 3','','published',13,4,1,'2026-07-10 20:57:47.984','2026-07-10 20:57:47.984'),(16,'用 Markdown 打造可读的技术文章','整理标题、列表、代码块和链接等常用 Markdown 写作技巧。','# Markdown 写作技巧\n\n一篇好的技术文章需要清晰的结构和易读的表达。本文整理几个常用 Markdown 语法，帮助你把知识写得更清楚。\n\n## 推荐习惯\n\n- 使用层级标题组织内容\n- 用代码块展示完整示例\n- 为关键结论补充链接和说明','','published',6,4,1,'2026-07-11 10:45:56.195','2026-07-11 10:45:56.195'),(17,'GitHub Actions 自动部署实践','从代码提交到构建发布，搭建一条简单可靠的 CI/CD 流程。','# GitHub Actions 自动部署实践\n\n通过工作流文件可以把测试、构建和部署串联起来。每次提交代码后，系统会自动执行检查并发布最新版本。\n\n## 流程建议\n\n1. 先运行 lint 和测试\n2. 再构建前端资源\n3. 最后发布到服务器','','published',2,4,4,'2026-07-11 10:45:56.248','2026-07-11 10:45:56.248'),(18,'React Hooks 状态管理指南','理解 useState、useEffect 和自定义 Hook 的组合方式。','# React Hooks 状态管理指南\n\nHooks 让组件可以直接管理状态和副作用。合理拆分自定义 Hook，可以减少重复逻辑并提升可测试性。\n\n## 实践要点\n\n- 让每个 Hook 只负责一类逻辑\n- 避免在依赖数组中遗漏变量\n- 把可复用状态抽到独立 Hook 中','','published',0,4,2,'2026-07-11 10:45:56.256','2026-07-11 10:45:56.256'),(19,'Nginx 反向代理配置清单','记录部署 Web 服务时最常用的 Nginx 配置项和排查思路。','# Nginx 反向代理配置清单\n\nNginx 可以统一处理域名、静态资源和后端服务转发。配置时应优先保证代理头、超时和缓存策略清晰可控。\n\n## 排查顺序\n\n1. 检查端口是否监听\n2. 检查 upstream 是否可达\n3. 查看访问日志和错误日志','','published',16,4,4,'2026-07-11 10:45:56.263','2026-07-11 10:45:56.263'),(20,'测试文章：功能测试','这是一篇用于功能测试的文章','# 功能测试文章\n\n## 测试目的\n\n验证博客系统的文章创建功能。\n\n## 测试内容\n\n- 标题填写 ✓\n- 分类选择 ✓\n- 标签选择 ✓\n- 摘要填写 ✓\n- 内容编写 ✓\n\n## 结论\n\n所有功能正常运行。','','draft',0,4,1,'2026-07-13 10:25:39.139','2026-07-13 10:25:39.139');
/*!40000 ALTER TABLE `articles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `categories`
--

DROP TABLE IF EXISTS `categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `categories` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_categories_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `categories`
--

LOCK TABLES `categories` WRITE;
/*!40000 ALTER TABLE `categories` DISABLE KEYS */;
INSERT INTO `categories` VALUES (1,'技术分享','','2026-06-24 09:30:18.128','2026-06-24 09:30:18.128'),(2,'前端开发','','2026-06-24 09:32:16.336','2026-06-24 09:32:16.336'),(3,'后端架构','','2026-06-24 09:32:16.341','2026-06-24 09:32:16.341'),(4,'运维部署','','2026-06-24 09:32:16.346','2026-06-24 09:32:16.346'),(5,'学习笔记','','2026-06-24 09:32:16.350','2026-06-24 09:32:16.350'),(6,'算法','','2026-06-24 09:55:51.288','2026-06-24 09:55:51.288'),(7,'英语','','2026-07-10 20:42:50.366','2026-07-10 20:42:50.366'),(8,'数据库','','2026-07-13 08:51:33.213','2026-07-13 08:51:33.213'),(9,'人工智能','','2026-07-13 08:51:33.228','2026-07-13 08:51:33.228'),(10,'云原生','','2026-07-13 08:51:33.232','2026-07-13 08:51:33.232'),(11,'网络安全','','2026-07-13 08:51:33.234','2026-07-13 08:51:33.234'),(12,'生活随笔','','2026-07-13 08:51:33.237','2026-07-13 08:51:33.237'),(13,'测试分类','这是一个用于功能测试的分类','2026-07-13 10:26:13.682','2026-07-13 10:26:13.682');
/*!40000 ALTER TABLE `categories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `comments`
--

DROP TABLE IF EXISTS `comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `comments` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `article_id` bigint(20) unsigned NOT NULL,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `content` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` tinyint(4) NOT NULL DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `parent_id` bigint(20) unsigned DEFAULT NULL,
  `reply_to_id` bigint(20) unsigned DEFAULT NULL,
  `guest_name` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `guest_email` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `guest_website` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_comments_article` (`article_id`),
  KEY `fk_comments_user` (`user_id`),
  KEY `fk_comments_parent` (`parent_id`),
  KEY `fk_comments_reply_to` (`reply_to_id`),
  CONSTRAINT `fk_comments_article` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`) ON UPDATE CASCADE,
  CONSTRAINT `fk_comments_parent` FOREIGN KEY (`parent_id`) REFERENCES `comments` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_comments_reply_to` FOREIGN KEY (`reply_to_id`) REFERENCES `comments` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_comments_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `comments`
--

LOCK TABLES `comments` WRITE;
/*!40000 ALTER TABLE `comments` DISABLE KEYS */;
INSERT INTO `comments` VALUES (1,19,4,'wqr',1,'2026-07-13 09:25:04.348','2026-07-13 09:25:04.348',NULL,NULL,NULL,NULL,NULL),(2,19,4,'111',1,'2026-07-13 10:11:49.695','2026-07-13 10:11:49.695',NULL,NULL,NULL,NULL,NULL),(3,19,4,'这是一条测试评论',1,'2026-07-13 10:24:18.997','2026-07-13 10:24:18.997',NULL,NULL,NULL,NULL,NULL),(4,16,7,'111',1,'2026-07-13 15:44:28.857','2026-07-13 15:44:28.857',NULL,NULL,NULL,NULL,NULL),(5,16,7,'111',1,'2026-07-13 15:59:15.048','2026-07-13 15:59:15.048',NULL,NULL,NULL,NULL,NULL),(6,19,NULL,'111',1,'2026-07-13 16:08:06.132','2026-07-13 16:10:27.821',NULL,NULL,'aaa','dzzd98784@gmail.com',''),(7,10,4,'1',1,'2026-07-13 19:59:18.780','2026-07-13 19:59:18.780',NULL,NULL,'','','');
/*!40000 ALTER TABLE `comments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tags`
--

DROP TABLE IF EXISTS `tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tags` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tags_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tags`
--

LOCK TABLES `tags` WRITE;
/*!40000 ALTER TABLE `tags` DISABLE KEYS */;
INSERT INTO `tags` VALUES (1,'Markdown','2026-07-11 10:45:18.895','2026-07-11 10:45:18.895'),(2,'Git','2026-07-11 10:45:18.933','2026-07-11 10:45:18.933'),(3,'GitHub Actions','2026-07-11 10:45:18.937','2026-07-11 10:45:18.937'),(4,'Nginx','2026-07-11 10:45:18.941','2026-07-11 10:45:18.941'),(5,'RESTful API','2026-07-11 10:45:18.945','2026-07-11 10:45:18.945'),(6,'微服务','2026-07-11 10:45:18.947','2026-07-11 10:45:18.947'),(7,'React Hooks','2026-07-11 10:45:18.950','2026-07-11 10:45:18.950'),(8,'Docker','2026-07-11 10:45:18.953','2026-07-11 10:45:18.953'),(9,'Go','2026-07-11 10:45:18.956','2026-07-11 10:45:18.956'),(10,'英语学习','2026-07-11 10:45:18.959','2026-07-11 10:45:18.959'),(11,'测试标签','2026-07-13 10:26:33.227','2026-07-13 10:26:33.227');
/*!40000 ALTER TABLE `tags` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `nickname` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `role` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'user',
  `avatar` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `status` tinyint(4) NOT NULL DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'dong','$2a$10$jKgSi76j1d8GNVdtibZ0Te6SKDVeFVrnrby97USpDAf7xUBpxDcrS','dong','user','',1,'2026-06-15 16:58:34.772','2026-06-15 16:58:34.772'),(2,'dsfafag','$2a$10$WHHERn3z5LwHpkxtdWEu5eYouKId6Cg/lRqa6oJYOEI1PFWdSCh1.','dsegsgsgesg','user','',1,'2026-06-15 17:17:30.814','2026-06-15 17:17:30.814'),(3,'testuser','$2a$10$rtGVpxmkBzbTHYqii2Nf7eD0kxu4/e1tObgY66gmXDFkbUMsv2bGu','�����û�','user','',1,'2026-06-15 17:18:11.709','2026-07-10 14:54:15.882'),(4,'admin','$2a$10$5jDi3mDw6ZJc9puVNtQxa.Oeop4GmIbGMJovSgXAucR196dg3Y452','管理员','admin','/uploads/articles/2026/07/20260710194712-acb07ee6adb9.png',1,'2026-06-24 09:27:14.679','2026-07-10 19:47:17.332'),(6,'demo','$2a$10$uN3kBMOfutQm6flyLBv4je5xcA5ZrrH5/kq/4PXd5vBx5vFqJs9.a','演示用户','user','',1,'2026-06-24 09:44:32.123','2026-06-24 09:44:32.123'),(7,'ddz','$2a$10$e/ywSuvrC/.AMJrcELMi1OGYYdG9ig6ed5cHXow3Z9Ar..dnPkvFi','ddz','user','',1,'2026-07-10 11:45:15.573','2026-07-10 14:54:15.149');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'gin_blog'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-07-14 11:29:21
