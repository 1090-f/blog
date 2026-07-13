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
	ErrTagNotFound = errors.New("tag not found")
	ErrTagExists   = errors.New("tag already exists")
	ErrInvalidTag  = errors.New("invalid tag params")
	ErrTagInUse    = errors.New("tag is in use by existing articles")
)

type TagStore interface {
	List() ([]model.Tag, error)
	ListByIDs(ids []uint) ([]model.Tag, error)
	Create(tag *model.Tag) error
	FindByID(id uint) (*model.Tag, error)
	Update(tag *model.Tag) error
	Delete(id uint) error
	CountByTagID(id uint) (int64, error)
	ReplaceArticleTags(articleID uint, tagIDs []uint) error
}

type TagService struct {
	tagStore TagStore
}

func NewTagService(tagStore TagStore) *TagService {
	return &TagService{tagStore: tagStore}
}

func (s *TagService) List() ([]model.Tag, error) {
	return s.tagStore.List()
}

func (s *TagService) Create(req dto.CreateTagRequest) (*model.Tag, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" || len(name) > 50 {
		return nil, ErrInvalidTag
	}

	tag := &model.Tag{Name: name}
	if err := s.tagStore.Create(tag); err != nil {
		var mysqlErr *mysqlDriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, ErrTagExists
		}
		return nil, err
	}
	return tag, nil
}

func (s *TagService) Update(id uint, req dto.UpdateTagRequest) (*model.Tag, error) {
	name := strings.TrimSpace(req.Name)
	if id == 0 || name == "" || len(name) > 50 {
		return nil, ErrInvalidTag
	}

	tag, err := s.tagStore.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTagNotFound
		}
		return nil, err
	}
	tag.Name = name
	if err := s.tagStore.Update(tag); err != nil {
		var mysqlErr *mysqlDriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, ErrTagExists
		}
		return nil, err
	}
	return tag, nil
}

func (s *TagService) Delete(id uint) error {
	if id == 0 {
		return ErrTagNotFound
	}
	if _, err := s.tagStore.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTagNotFound
		}
		return err
	}

	count, err := s.tagStore.CountByTagID(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrTagInUse
	}
	return s.tagStore.Delete(id)
}

func NormalizeTagIDs(ids []uint) []uint {
	seen := make(map[uint]struct{}, len(ids))
	result := make([]uint, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}
