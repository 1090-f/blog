<template>
  <div class="article-list-page">
    <div class="action-bar">
      <div class="filter-group">
        <input
          v-model="keyword"
          class="input search-input"
          type="text"
          placeholder="搜索标题关键字"
          @keyup.enter="handleSearch"
        />
        <select v-model="categoryId" class="input filter-select" @change="handleSearch">
          <option value="">全部分类</option>
          <option v-for="category in categories" :key="category.id" :value="String(category.id)">
            {{ category.name }}
          </option>
        </select>
        <select v-model="tagId" class="input filter-select" @change="handleSearch">
          <option value="">全部标签</option>
          <option v-for="tag in tags" :key="tag.id" :value="String(tag.id)">
            {{ tag.name }}
          </option>
        </select>
        <select v-model="statusFilter" class="input filter-select" @change="handleSearch">
          <option value="">全部状态</option>
          <option value="published">已发布</option>
          <option value="draft">草稿</option>
        </select>
        <button class="btn btn-outline" @click="resetFilters">重置</button>
      </div>
      <button class="btn btn-primary" @click="$router.push('/admin/article/new')">新建文章</button>
    </div>

    <div class="card">
      <div v-if="loading" class="loading-state">加载中...</div>
      <div v-else-if="articles.length > 0" class="article-table">
        <div v-for="article in articles" :key="article.id" class="article-row">
          <div class="article-main">
            <h4 class="article-title">{{ article.title }}</h4>
            <p v-if="article.summary" class="article-summary">{{ article.summary }}</p>
            <div class="article-meta">
              <span class="tag">{{ article.category?.name || '未分类' }}</span>
              <span v-for="tag in article.tags || []" :key="tag.id" class="article-tag">{{ tag.name }}</span>
              <span :class="['status-tag', article.status === 'published' ? 'status-published' : 'status-draft']">
                {{ article.status === 'published' ? '已发布' : '草稿' }}
              </span>
              <span class="meta-text">浏览 {{ article.viewCount || 0 }}</span>
              <span class="meta-text">{{ formatDate(article.createdAt) }}</span>
            </div>
          </div>
          <div class="article-actions">
            <button class="btn btn-outline btn-sm" @click="$router.push(`/admin/article/edit/${article.id}`)">编辑</button>
            <button class="btn btn-danger btn-sm" @click="handleDelete(article.id)">删除</button>
          </div>
        </div>
      </div>
      <div v-else class="empty-state">暂无文章</div>
    </div>

    <div v-if="total > pageSize" class="pagination-wrap">
      <button class="btn btn-outline" :disabled="page <= 1" @click="changePage(page - 1)">上一页</button>
      <span class="page-info">{{ page }} / {{ Math.ceil(total / pageSize) }}</span>
      <button class="btn btn-outline" :disabled="page >= Math.ceil(total / pageSize)" @click="changePage(page + 1)">下一页</button>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { deleteArticle, getAdminArticles } from '../../api/article'
import { getAdminCategories } from '../../api/category'
import { getAdminTags } from '../../api/tag'
import { message } from '../../utils/message'

const articles = ref([])
const categories = ref([])
const tags = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = 15
const total = ref(0)
const keyword = ref('')
const categoryId = ref('')
const tagId = ref('')
const statusFilter = ref('')

// 将原始数据格式化为界面展示内容。
function formatDate(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

// 处理用户操作或浏览器事件。
function handleSearch() {
  page.value = 1
  fetchData()
}

// 重置当前筛选条件。
function resetFilters() {
  keyword.value = ''
  categoryId.value = ''
  tagId.value = ''
  statusFilter.value = ''
  handleSearch()
}

// 更新当前筛选条件或页面状态。
function changePage(nextPage) {
  page.value = nextPage
  fetchData()
}

// 加载当前页面所需的数据。
async function fetchCategories() {
  const res = await getAdminCategories()
  categories.value = res.data || []
}

// 加载当前页面所需的数据。
async function fetchTags() {
  const res = await getAdminTags()
  tags.value = res.data || []
}

// 加载当前页面所需的数据。
async function fetchData() {
  loading.value = true
  try {
    const params = { page: page.value, pageSize }
    if (keyword.value.trim()) params.keyword = keyword.value.trim()
    if (categoryId.value) params.categoryId = Number(categoryId.value)
    if (tagId.value) params.tagId = Number(tagId.value)
    if (statusFilter.value) params.status = statusFilter.value

    const res = await getAdminArticles(params)
    articles.value = res.data?.list || []
    total.value = res.data?.total || 0
  } finally {
    loading.value = false
  }
}

// 处理用户操作或浏览器事件。
async function handleDelete(id) {
  if (!window.confirm('确认删除这篇文章吗？')) return
  try {
    await deleteArticle(id)
    message.success('删除成功')
    await fetchData()
  } catch (error) {
    message.error('删除失败')
  }
}

onMounted(async () => {
  await Promise.all([fetchCategories(), fetchTags(), fetchData()])
})
</script>

<style scoped>
.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.filter-group {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.search-input {
  min-width: 240px;
}

.filter-select {
  width: 160px;
}

.article-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 18px 0;
  border-bottom: 1px solid var(--border);
}

.article-row:last-child {
  border-bottom: none;
}

.article-main {
  flex: 1;
  min-width: 0;
}

.article-title {
  font-size: 15px;
  font-weight: 600;
  margin-bottom: 8px;
  color: var(--text-primary);
}

.article-summary {
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.6;
  margin-bottom: 10px;
}

.article-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.tag {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 4px;
  background: var(--accent-dim);
  color: var(--accent);
  font-size: 12px;
}

.article-tag {
  display: inline-block;
  padding: 2px 8px;
  border: 1px solid var(--border);
  border-radius: 4px;
  color: var(--text-secondary);
  font-size: 12px;
}

.status-tag {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 4px;
  font-size: 12px;
}

.status-published {
  background: rgba(45, 212, 191, 0.15);
  color: var(--accent);
}

.status-draft {
  background: rgba(156, 163, 175, 0.15);
  color: var(--text-muted);
}

.meta-text {
  color: var(--text-muted);
  font-size: 13px;
}

.article-actions {
  display: flex;
  gap: 8px;
}

.btn-sm {
  padding: 6px 12px;
  font-size: 13px;
}

.btn-danger {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.btn-danger:hover {
  background: #ef4444;
  color: white;
}

.loading-state,
.empty-state {
  text-align: center;
  padding: 40px;
  color: var(--text-muted);
}

.pagination-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-top: 20px;
}

.page-info {
  color: var(--text-secondary);
  font-size: 14px;
}

@media (max-width: 960px) {
  .action-bar {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
