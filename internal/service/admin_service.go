package service

type AdminDashboardStats struct {
	ArticleCount   int64 `json:"articleCount"`
	PublishedCount int64 `json:"publishedCount"`
	DraftCount     int64 `json:"draftCount"`
	CategoryCount  int64 `json:"categoryCount"`
	UserCount      int64 `json:"userCount"`
	TotalViews     int64 `json:"totalViews"`
	CommentCount   int64 `json:"commentCount"`
}

type AdminStatsArticleStore interface {
	CountAll() (int64, error)
	CountByStatus(status string) (int64, error)
	SumViewCount() (int64, error)
}

type AdminStatsCategoryStore interface {
	CountAll() (int64, error)
}

type AdminStatsUserStore interface {
	CountAll() (int64, error)
}

type AdminStatsCommentStore interface {
	CountAll() (int64, error)
}

type AdminService struct {
	articleStore  AdminStatsArticleStore
	categoryStore AdminStatsCategoryStore
	userStore     AdminStatsUserStore
	commentStore  AdminStatsCommentStore
}

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
