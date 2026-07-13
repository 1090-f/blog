package utils

const (
	defaultPage     = 1
	defaultPageSize = 10
	maxPageSize     = 100
)

func NormalizePage(page, pageSize int) (normalizedPage int, normalizedPageSize int, offset int, limit int) {
	if page < 1 {
		page = defaultPage
	}
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	return page, pageSize, (page - 1) * pageSize, pageSize
}
