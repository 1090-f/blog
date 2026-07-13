package service

import (
	"blog/internal/dao"
	"fmt"
	"time"
)

type DailyActivity struct {
	Date     string `json:"date"`
	Articles int64  `json:"articles"`
	Comments int64  `json:"comments"`
	Total    int64  `json:"total"`
}

type ActivityResponse struct {
	Year  int             `json:"year"`
	Month int             `json:"month"`
	Days  []DailyActivity `json:"days"`
}

type ActivityStore interface {
	PublishedArticlesByDay(start, end time.Time) ([]dao.DailyActivityCount, error)
	ApprovedCommentsByDay(start, end time.Time) ([]dao.DailyActivityCount, error)
}

type ActivityService struct {
	store ActivityStore
}

func NewActivityService(store ActivityStore) *ActivityService {
	return &ActivityService{store: store}
}

func (s *ActivityService) Get(year, month int) (*ActivityResponse, error) {
	if year == 0 || month == 0 {
		now := time.Now()
		if year == 0 {
			year = now.Year()
		}
		if month == 0 {
			month = int(now.Month())
		}
	}
	if month < 1 || month > 12 {
		return nil, fmt.Errorf("invalid activity month: %d", month)
	}

	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, 0)
	articleRows, err := s.store.PublishedArticlesByDay(start, end)
	if err != nil {
		return nil, err
	}
	commentRows, err := s.store.ApprovedCommentsByDay(start, end)
	if err != nil {
		return nil, err
	}

	articlesByDay := make(map[string]int64, len(articleRows))
	commentsByDay := make(map[string]int64, len(commentRows))
	for _, row := range articleRows {
		articlesByDay[row.Day.Format("2006-01-02")] = row.Count
	}
	for _, row := range commentRows {
		commentsByDay[row.Day.Format("2006-01-02")] = row.Count
	}

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
