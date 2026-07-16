package service

import (
	"errors"
	"strings"

	"blog/internal/dto"
	"blog/internal/model"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// 分类服务相关的错误值。
var (
	ErrCategoryNotFound = errors.New("category not found")                      // 分类不存在
	ErrCategoryExists   = errors.New("category already exists")                 // 分类名已存在
	ErrInvalidCategory  = errors.New("invalid category params")                 // 分类参数无效
	ErrCategoryInUse    = errors.New("category is in use by existing articles") // 分类下有文章，无法删除
)

// CategoryService 分类业务逻辑层，处理分类 CRUD 操作。
type CategoryService struct {
	categoryDAO CategoryStore
	articleDAO  CategoryArticleStore
}

// CategoryStore 分类服务所需的持久化操作抽象。
type CategoryStore interface {
	List() ([]model.Category, error)
	Create(category *model.Category) error
	FindByID(id uint) (*model.Category, error)
	Update(category *model.Category) error
	Delete(id uint) error
}

// CategoryArticleStore 分类服务所需的文章数量统计操作抽象。
type CategoryArticleStore interface {
	CountByCategoryID(categoryID uint) (int64, error)
}

// NewCategoryService 创建并初始化分类业务实例。
func NewCategoryService(categoryDAO CategoryStore, articleDAO CategoryArticleStore) *CategoryService {
	return &CategoryService{
		categoryDAO: categoryDAO,
		articleDAO:  articleDAO,
	}
}

// List 查询所有分类列表，按数据库默认排序返回。
func (s *CategoryService) List() ([]model.Category, error) {
	return s.categoryDAO.List()
}

// Create 创建新分类，校验名称和描述长度后持久化，名称重复时返回 ErrCategoryExists。
func (s *CategoryService) Create(req dto.CreateCategoryRequest) (*model.Category, error) {
	name := strings.TrimSpace(req.Name)
	description := strings.TrimSpace(req.Description)
	// 名称不能为空，名称和描述有长度上限
	if name == "" || len(name) > 50 || len(description) > 255 {
		return nil, ErrInvalidCategory
	}

	category := &model.Category{
		Name:        name,
		Description: description,
	}

	if err := s.categoryDAO.Create(category); err != nil {
		// 判断是否是唯一约束错误
		var mysqlErr *mysqlDriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, ErrCategoryExists
		}
		return nil, err
	}

	return category, nil
}

// Update 更新已有分类的名称和描述，校验参数后持久化，名称重复时返回 ErrCategoryExists。
func (s *CategoryService) Update(id uint, req dto.UpdateCategoryRequest) (*model.Category, error) {
	name := strings.TrimSpace(req.Name)
	description := strings.TrimSpace(req.Description)
	// ID 不能为空，名称和描述有长度上限
	if id == 0 || name == "" || len(name) > 50 || len(description) > 255 {
		return nil, ErrInvalidCategory
	}

	// 先确认分类存在
	category, err := s.categoryDAO.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	category.Name = name
	category.Description = description

	// 持久化更新后的分类字段
	if err := s.categoryDAO.Update(category); err != nil {
		// MySQL 错误码 1062 = 唯一约束冲突，说明分类名已存在
		var mysqlErr *mysqlDriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, ErrCategoryExists
		}
		return nil, err
	}

	return category, nil
}

// Delete 删除指定分类，若分类下仍有文章则拒绝删除。
func (s *CategoryService) Delete(id uint) error {
	// 校验id
	if id == 0 {
		return ErrCategoryNotFound
	}
	// 检查分类是否存在
	if _, err := s.categoryDAO.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCategoryNotFound
		}
		return err
	}
	// 检查是否有关联文章
	count, err := s.articleDAO.CountByCategoryID(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrCategoryInUse
	}

	return s.categoryDAO.Delete(id)
}
