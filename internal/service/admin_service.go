package service

// AdminDashboardStats 管理端仪表盘聚合统计数据。
type AdminDashboardStats struct {
	ArticleCount   int64 `json:"articleCount"`
	PublishedCount int64 `json:"publishedCount"`
	DraftCount     int64 `json:"draftCount"`
	CategoryCount  int64 `json:"categoryCount"`
	UserCount      int64 `json:"userCount"`
	TotalViews     int64 `json:"totalViews"`
	CommentCount   int64 `json:"commentCount"`
}

// AdminStatsArticleStore 仪表盘所需的文章统计操作抽象。
type AdminStatsArticleStore interface {
	CountAll() (int64, error)
	CountByStatus(status string) (int64, error)
	SumViewCount() (int64, error)
}

// AdminStatsCategoryStore 仪表盘所需的分类统计操作抽象。
type AdminStatsCategoryStore interface {
	CountAll() (int64, error)
}

// AdminStatsUserStore 仪表盘所需的用户统计操作抽象。
type AdminStatsUserStore interface {
	CountAll() (int64, error)
}

// AdminStatsCommentStore 仪表盘所需的评论统计操作抽象。
type AdminStatsCommentStore interface {
	CountAll() (int64, error)
}

// AdminService 管理端仪表盘统计业务逻辑层。
type AdminService struct {
	articleStore  AdminStatsArticleStore
	categoryStore AdminStatsCategoryStore
	userStore     AdminStatsUserStore
	commentStore  AdminStatsCommentStore
}

// 创建并初始化管理端统计实例。
func NewAdminService(articleStore AdminStatsArticleStore, categoryStore AdminStatsCategoryStore, userStore AdminStatsUserStore, commentStores ...AdminStatsCommentStore) *AdminService {
	var commentStore AdminStatsCommentStore
	if len(commentStores) > 0 {
		commentStore = commentStores[0]
	}
	return &AdminService{
		articleStore:  articleStore,
		categoryStore: categoryStore,
		userStore:     userStore,
		commentStore:  commentStore,
	}
}

// 获取管理端仪表盘统计数据。
func (s *AdminService) Dashboard() (*AdminDashboardStats, error) {
	articleCount, err := s.articleStore.CountAll()
	if err != nil {
		return nil, err
	}
	publishedCount, err := s.articleStore.CountByStatus(ArticleStatusPublished)
	if err != nil {
		return nil, err
	}
	draftCount, err := s.articleStore.CountByStatus(ArticleStatusDraft)
	if err != nil {
		return nil, err
	}
	categoryCount, err := s.categoryStore.CountAll()
	if err != nil {
		return nil, err
	}
	userCount, err := s.userStore.CountAll()
	if err != nil {
		return nil, err
	}
	totalViews, err := s.articleStore.SumViewCount()
	if err != nil {
		return nil, err
	}
	var commentCount int64
	if s.commentStore != nil {
		commentCount, err = s.commentStore.CountAll()
		if err != nil {
			return nil, err
		}
	}

	return &AdminDashboardStats{
		ArticleCount:   articleCount,
		PublishedCount: publishedCount,
		DraftCount:     draftCount,
		CategoryCount:  categoryCount,
		UserCount:      userCount,
		TotalViews:     totalViews,
		CommentCount:   commentCount,
	}, nil
}
