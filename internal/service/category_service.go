package service

import (
	"errors"
	"strings"

	"blog/internal/dto"
	"blog/internal/model"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryExists   = errors.New("category already exists")
	ErrInvalidCategory  = errors.New("invalid category params")
	ErrCategoryInUse    = errors.New("category is in use by existing articles")
)

type CategoryService struct {
	categoryDAO CategoryStore
	articleDAO  CategoryArticleStore
}

type CategoryStore interface {
	List() ([]model.Category, error)
	Create(category *model.Category) error
	FindByID(id uint) (*model.Category, error)
	Update(category *model.Category) error
	Delete(id uint) error
}

type CategoryArticleStore interface {
	CountByCategoryID(categoryID uint) (int64, error)
}

func NewCategoryService(categoryDAO CategoryStore, articleDAO CategoryArticleStore) *CategoryService {
	return &CategoryService{
		categoryDAO: categoryDAO,
		articleDAO:  articleDAO,
	}
}

func (s *CategoryService) List() ([]model.Category, error) {
	return s.categoryDAO.List()
}

func (s *CategoryService) Create(req dto.CreateCategoryRequest) (*model.Category, error) {
	name := strings.TrimSpace(req.Name)
	description := strings.TrimSpace(req.Description)
	if name == "" || len(name) > 50 || len(description) > 255 {
		return nil, ErrInvalidCategory
	}

	category := &model.Category{
		Name:        name,
		Description: description,
	}

	if err := s.categoryDAO.Create(category); err != nil {
		var mysqlErr *mysqlDriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, ErrCategoryExists
		}
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) Update(id uint, req dto.UpdateCategoryRequest) (*model.Category, error) {
	name := strings.TrimSpace(req.Name)
	description := strings.TrimSpace(req.Description)
	if id == 0 || name == "" || len(name) > 50 || len(description) > 255 {
		return nil, ErrInvalidCategory
	}

	category, err := s.categoryDAO.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	category.Name = name
	category.Description = description

	if err := s.categoryDAO.Update(category); err != nil {
		var mysqlErr *mysqlDriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, ErrCategoryExists
		}
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) Delete(id uint) error {
	if id == 0 {
		return ErrCategoryNotFound
	}

	if _, err := s.categoryDAO.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCategoryNotFound
		}
		return err
	}

	count, err := s.articleDAO.CountByCategoryID(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrCategoryInUse
	}

	return s.categoryDAO.Delete(id)
}
