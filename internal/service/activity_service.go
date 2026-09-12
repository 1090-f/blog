package service

import (
	"blog/internal/dao"
	"fmt"
	"time"
)

// DailyActivity 单日活动明细（文章数、评论数、总数）。
type DailyActivity struct {
	Date     string `json:"date"`
	Articles int64  `json:"articles"`
	Comments int64  `json:"comments"`
	Total    int64  `json:"total"`
}

// ActivityResponse 月度活动日历响应，按日列出活动明细。
type ActivityResponse struct {
	Year  int             `json:"year"`
	Month int             `json:"month"`
	Days  []DailyActivity `json:"days"`
}

// ActivityStore 活动统计服务所需的按日聚合操作抽象。
type ActivityStore interface {
	PublishedArticlesByDay(start, end time.Time) ([]dao.DailyActivityCount, error)
	ApprovedCommentsByDay(start, end time.Time) ([]dao.DailyActivityCount, error)
}

// ActivityService 活动日历统计业务逻辑层。
type ActivityService struct {
	store ActivityStore
}

// NewActivityService 创建并初始化活动统计实例。
func NewActivityService(store ActivityStore) *ActivityService {
	return &ActivityService{store: store}
}

// Get 获取指定月份的活动统计数据，返回该月每天的文章发布数和评论数。
// year 或 month 为 0 时自动取当前时间对应值。
func (s *ActivityService) Get(year, month int) (*ActivityResponse, error) {
	// 若未指定年或月，使用当前时间
	if year == 0 || month == 0 {
		now := time.Now()
		if year == 0 {
			year = now.Year()
		}
		if month == 0 {
			month = int(now.Month())
		}
	}
	// 校验月份合法性
	if month < 1 || month > 12 {
		return nil, fmt.Errorf("invalid activity month: %d", month)
	}

	// 计算查询时间范围：当月1日 ~ 下月1日（左闭右开）
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, 0)

	// 查询该月每天的已发布文章数和已审核评论数
	articleRows, err := s.store.PublishedArticlesByDay(start, end)
	if err != nil {
		return nil, err
	}
	commentRows, err := s.store.ApprovedCommentsByDay(start, end)
	if err != nil {
		return nil, err
	}

	// 将查询结果转为以日期为 key 的 map，便于后续按日遍历时快速查找
	articlesByDay := make(map[string]int64, len(articleRows))
	commentsByDay := make(map[string]int64, len(commentRows))
	for _, row := range articleRows {
		articlesByDay[row.Day.Format("2006-01-02")] = row.Count
	}
	for _, row := range commentRows {
		commentsByDay[row.Day.Format("2006-01-02")] = row.Count
	}

	// 逐日遍历，确保即使某天没有数据也会出现在结果中（值为 0）
	days := make([]DailyActivity, 0, end.Day())
	for day := start; day.Before(end); day = day.AddDate(0, 0, 1) {
		date := day.Format("2006-01-02")
		articles := articlesByDay[date]
		comments := commentsByDay[date]
		days = append(days, DailyActivity{
			Date:     date,
			Articles: articles,
			Comments: comments,
			Total:    articles + comments,
		})
	}

	return &ActivityResponse{Year: year, Month: month, Days: days}, nil
}
