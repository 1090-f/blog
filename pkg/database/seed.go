package database

import (
	"errors"
	"fmt"
	"time"

	"blog/internal/model"
	"gorm.io/gorm"
)

// sampleArticleDefinition 描述首次启动时写入的示例文章及其关联数据。
type sampleArticleDefinition struct {
	Title               string
	Summary             string
	Content             string
	CategoryName        string
	CategoryDescription string
	TagNames            []string
	DaysAgo             int
}

// EnsureSampleContent 在文章表为空时创建可直接展示的示例内容。
// 分类和标签按名称复用，全部写入在同一个事务中完成，避免产生半成品数据。
func EnsureSampleContent(db *gorm.DB) (bool, error) {
	seeded := false
	err := db.Transaction(func(tx *gorm.DB) error {
		var articleCount int64
		if err := tx.Model(&model.Article{}).Count(&articleCount).Error; err != nil {
			return fmt.Errorf("count articles: %w", err)
		}
		if articleCount > 0 {
			return nil
		}

		author, err := findSampleAuthor(tx)
		if err != nil {
			return err
		}

		now := time.Now()
		for _, definition := range defaultSampleArticles() {
			category := model.Category{
				Name:        definition.CategoryName,
				Description: definition.CategoryDescription,
			}
			if err := tx.Where("name = ?", category.Name).FirstOrCreate(&category).Error; err != nil {
				return fmt.Errorf("ensure sample category %q: %w", category.Name, err)
			}

			tags := make([]model.Tag, 0, len(definition.TagNames))
			for _, name := range definition.TagNames {
				tag := model.Tag{Name: name}
				if err := tx.Where("name = ?", name).FirstOrCreate(&tag).Error; err != nil {
					return fmt.Errorf("ensure sample tag %q: %w", name, err)
				}
				tags = append(tags, tag)
			}

			publishedAt := now.AddDate(0, 0, -definition.DaysAgo)
			article := model.Article{
				Title:      definition.Title,
				Summary:    definition.Summary,
				Content:    definition.Content,
				Status:     "published",
				ViewCount:  0,
				UserID:     author.ID,
				CategoryID: category.ID,
				Tags:       tags,
				CreatedAt:  publishedAt,
				UpdatedAt:  publishedAt,
			}
			if err := tx.Omit("User", "Category").Create(&article).Error; err != nil {
				return fmt.Errorf("create sample article %q: %w", article.Title, err)
			}
		}

		seeded = true
		return nil
	})
	return seeded, err
}

// findSampleAuthor 优先使用启用状态的管理员，确保示例文章满足作者外键约束。
func findSampleAuthor(tx *gorm.DB) (*model.User, error) {
	var author model.User
	err := tx.Where("role = ? AND status = ?", "admin", 1).Order("id ASC").First(&author).Error
	if err == nil {
		return &author, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("find sample author: %w", err)
	}

	err = tx.Where("status = ?", 1).Order("id ASC").First(&author).Error
	if err == nil {
		return &author, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("cannot seed sample content without an active user; enable admin bootstrap first")
	}
	return nil, fmt.Errorf("find fallback sample author: %w", err)
}

func defaultSampleArticles() []sampleArticleDefinition {
	return []sampleArticleDefinition{
		{
			Title:               "欢迎来到我的博客",
			Summary:             "这是一篇自动生成的欢迎文章，你可以从管理后台编辑或删除它。",
			CategoryName:        "站点公告",
			CategoryDescription: "博客更新、功能说明与站点动态",
			TagNames:            []string{"欢迎", "博客"},
			DaysAgo:             2,
			Content: `# 欢迎来到我的博客

你好，欢迎来到这里！这篇文章由系统在首次启动时自动创建，用来确认文章、分类和标签功能已经正常工作。

## 从哪里开始

- 在管理后台新建或编辑文章
- 使用 Markdown 编写正文
- 为文章选择分类和标签
- 上传封面与正文图片

当你发布自己的第一篇文章后，可以放心删除这篇示例内容。愿这个小站记录下值得珍藏的想法与经历。`,
		},
		{
			Title:               "用 Markdown 写下第一篇文章",
			Summary:             "快速了解标题、列表、引用、代码块和链接等常用 Markdown 写法。",
			CategoryName:        "写作指南",
			CategoryDescription: "关于内容创作与 Markdown 排版的实用技巧",
			TagNames:            []string{"Markdown", "写作"},
			DaysAgo:             1,
			Content: `# 用 Markdown 写下第一篇文章

Markdown 让你把注意力放在内容本身，同时保持清晰、稳定的排版。

## 常用语法

### 列表

1. 先确定文章主题
2. 写出内容提纲
3. 补充细节并发布

### 引用

> 写作不是等待灵感，而是给思考留下可以回看的痕迹。

### 代码

~~~go
package main

import "fmt"

func main() {
	fmt.Println("Hello, blog!")
}
~~~

你可以进入管理后台编辑本文，观察这些语法在页面中的实际效果。`,
		},
		{
			Title:               "从 Docker 部署开始记录",
			Summary:             "这个博客已经运行在容器中，数据卷、健康检查与更新流程都值得被记录。",
			CategoryName:        "技术实践",
			CategoryDescription: "开发、部署与运维过程中的经验记录",
			TagNames:            []string{"Docker", "Go", "部署"},
			DaysAgo:             0,
			Content: `# 从 Docker 部署开始记录

现在，这个博客已经由 Go、Vue、MySQL 和 Docker Compose 共同运行。

## 日常维护建议

- 更新前备份 MySQL 数据和上传文件
- 使用健康检查确认服务状态
- 不要使用 ` + "`docker compose down -v`" + `，避免删除持久化数据卷
- 定期更新依赖和服务器安全补丁
- 管理后台只向可信网络或 SSH 隧道开放

部署不是终点。把遇到的问题、解决方案和新的想法持续记录下来，博客才会真正成为自己的知识库。`,
		},
	}
}
