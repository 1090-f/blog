package utils

// 分页查询的默认值和上限，防止前端传入非法参数。
const (
	defaultPage     = 1   // 默认页码
	defaultPageSize = 10  // 默认每页条数
	maxPageSize     = 100 // 每页条数上限
)

// NormalizePage 将前端传入的页码和页大小规范化为安全值，
// 同时计算数据库查询需要的 offset（偏移量）和 limit（取多少条）。
func NormalizePage(page, pageSize int) (normalizedPage int, normalizedPageSize int, offset int, limit int) {
	// 页码小于 1 时使用默认值
	if page < 1 {
		page = defaultPage
	}
	// 页小于 1 时使用默认值
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	// 页大小超过上限时截断
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	// 计算 offset：第 N 页的起始位置 = (N - 1) * pageSize
	return page, pageSize, (page - 1) * pageSize, pageSize
}
