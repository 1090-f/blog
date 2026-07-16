<template>
  <div class="dashboard">
    <div class="stats-grid">
      <div class="stat-card card">
        <div class="stat-icon icon-articles">AR</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.articleCount }}</div>
          <div class="stat-label">文章总数</div>
        </div>
      </div>
      <div class="stat-card card">
        <div class="stat-icon icon-published">PB</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.publishedCount }}</div>
          <div class="stat-label">已发布文章</div>
        </div>
      </div>
      <div class="stat-card card">
        <div class="stat-icon icon-comments">CO</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.commentCount }}</div>
          <div class="stat-label">评论总数</div>
        </div>
      </div>
      <div class="stat-card card">
        <div class="stat-icon icon-views">VW</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.totalViews }}</div>
          <div class="stat-label">总浏览量</div>
        </div>
      </div>
    </div>

    <div class="summary-grid">
      <div class="card summary-card">
        <h3 class="card-title">内容概览</h3>
        <div class="summary-list">
          <div class="summary-item">
            <span>草稿文章</span>
            <strong>{{ stats.draftCount }}</strong>
          </div>
          <div class="summary-item">
            <span>分类数量</span>
            <strong>{{ stats.categoryCount }}</strong>
          </div>
          <div class="summary-item">
            <span>用户数量</span>
            <strong>{{ stats.userCount }}</strong>
          </div>
        </div>
      </div>

      <div class="card summary-card">
        <h3 class="card-title">最近文章</h3>
        <div v-if="loadingRecent" class="empty-state">加载中...</div>
        <div v-else-if="recentArticles.length === 0" class="empty-state">暂无文章</div>
        <div v-else class="recent-list">
          <div v-for="article in recentArticles" :key="article.id" class="recent-item">
            <div class="recent-main">
              <div class="recent-title">{{ article.title }}</div>
              <div class="recent-meta">
                <span>{{ article.category?.name || '未分类' }}</span>
                <span>{{ article.status === 'published' ? '已发布' : '草稿' }}</span>
                <span>{{ formatDate(article.createdAt) }}</span>
              </div>
            </div>
            <div class="recent-views">{{ article.viewCount || 0 }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { getDashboardStats } from '../../api/admin'
import { getAdminArticles } from '../../api/article'

const stats = ref({
  articleCount: 0,
  publishedCount: 0,
  draftCount: 0,
  categoryCount: 0,
  userCount: 0,
  totalViews: 0,
  commentCount: 0
})
const recentArticles = ref([])
const loadingRecent = ref(false)

// 将原始数据格式化为界面展示内容。
function formatDate(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

// 加载当前页面所需的数据。
async function fetchDashboard() {
  const res = await getDashboardStats()
  stats.value = { ...stats.value, ...(res.data || {}) }
}

// 加载当前页面所需的数据。
async function fetchRecentArticles() {
  loadingRecent.value = true
  try {
    const res = await getAdminArticles({ page: 1, pageSize: 5 })
    recentArticles.value = res.data?.list || []
  } finally {
    loadingRecent.value = false
  }
}

onMounted(async () => {
  try {
    await Promise.all([fetchDashboard(), fetchRecentArticles()])
  } catch (error) {
    recentArticles.value = recentArticles.value || []
  }
})
</script>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 20px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 700;
  color: #071019;
  flex-shrink: 0;
}

.icon-articles {
  background: linear-gradient(135deg, #2dd4bf, #14b8a6);
}

.icon-published {
  background: linear-gradient(135deg, #60a5fa, #3b82f6);
}

.icon-comments {
  background: linear-gradient(135deg, #f59e0b, #f97316);
}

.icon-views {
  background: linear-gradient(135deg, #facc15, #eab308);
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
}

.stat-label {
  font-size: 13px;
  color: var(--text-muted);
  margin-top: 4px;
}

.summary-grid {
  display: grid;
  grid-template-columns: 1fr 1.4fr;
  gap: 20px;
  margin-top: 24px;
}

.summary-card {
  padding: 24px;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 18px;
}

.summary-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.summary-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 0;
  border-bottom: 1px solid var(--border);
  color: var(--text-secondary);
}

.summary-item:last-child {
  border-bottom: none;
}

.summary-item strong {
  color: var(--text-primary);
  font-size: 18px;
}

.recent-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.recent-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding-bottom: 14px;
  border-bottom: 1px solid var(--border);
}

.recent-item:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.recent-main {
  min-width: 0;
}

.recent-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.recent-meta {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  margin-top: 6px;
  color: var(--text-muted);
  font-size: 12px;
}

.recent-views {
  color: var(--text-secondary);
  font-weight: 600;
}

.empty-state {
  text-align: center;
  padding: 32px;
  color: var(--text-muted);
}

@media (max-width: 960px) {
  .stats-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .summary-grid {
    grid-template-columns: 1fr;
  }
}
</style>
