package service

import (
	"errors"
	"strings"

	"blog/internal/dto"
	"blog/internal/model"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// 标签服务相关的错误哨兵值。
var (
	ErrTagNotFound = errors.New("tag not found")      // 标签不存在
	ErrTagExists   = errors.New("tag already exists")  // 标签名已存在
	ErrInvalidTag  = errors.New("invalid tag params")  // 标签参数无效
	ErrTagInUse    = errors.New("tag is in use by existing articles") // 标签下有文章，无法删除
)

// TagStore 标签服务所需的持久化操作抽象。
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

// TagService 标签业务逻辑层，处理标签 CRUD 操作。
type TagService struct {
	tagStore TagStore
}

// NewTagService 创建并初始化标签业务实例。
func NewTagService(tagStore TagStore) *TagService {
	return &TagService{tagStore: tagStore}
}

// List 查询所有标签列表，按数据库默认排序返回。
func (s *TagService) List() ([]model.Tag, error) {
	return s.tagStore.List()
}

// Create 创建新标签，校验名称非空且长度合法后持久化，名称重复时返回 ErrTagExists。
func (s *TagService) Create(req dto.CreateTagRequest) (*model.Tag, error) {
	name := strings.TrimSpace(req.Name)
	// 名称不能为空，且不能超过 50 个字符
	if name == "" || len(name) > 50 {
		return nil, ErrInvalidTag
	}

	tag := &model.Tag{Name: name}
	if err := s.tagStore.Create(tag); err != nil {
		// MySQL 错误码 1062 = 唯一约束冲突，说明标签名已存在
		var mysqlErr *mysqlDriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, ErrTagExists
		}
		return nil, err
	}
	return tag, nil
}

// Update 更新已有标签的名称，校验参数后持久化，名称重复时返回 ErrTagExists。
func (s *TagService) Update(id uint, req dto.UpdateTagRequest) (*model.Tag, error) {
	name := strings.TrimSpace(req.Name)
	// ID 不能为空，名称不能为空且不能超过 50 个字符
	if id == 0 || name == "" || len(name) > 50 {
		return nil, ErrInvalidTag
	}

	// 先确认标签存在
	tag, err := s.tagStore.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTagNotFound
		}
		return nil, err
	}
	tag.Name = name
	// 持久化更新后的标签名称
	if err := s.tagStore.Update(tag); err != nil {
		// MySQL 错误码 1062 = 唯一约束冲突，说明标签名已存在
		var mysqlErr *mysqlDriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, ErrTagExists
		}
		return nil, err
	}
	return tag, nil
}

// Delete 删除指定标签，若标签下仍有文章则拒绝删除。
func (s *TagService) Delete(id uint) error {
	if id == 0 {
		return ErrTagNotFound
	}
	// 先确认标签存在
	if _, err := s.tagStore.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTagNotFound
		}
		return err
	}

	// 检查是否有关联文章
	count, err := s.tagStore.CountByTagID(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrTagInUse
	}
	return s.tagStore.Delete(id)
}

// NormalizeTagIDs 去除标签 ID 列表中的零值和重复项，返回去重排序后的结果。
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
