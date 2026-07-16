<template>
  <div class="comment-page">
    <div class="action-bar">
      <h3 class="page-subtitle">评论管理</h3>
      <div class="filter-row">
        <input
          v-model="keyword"
          class="input search-input"
          type="text"
          placeholder="搜索评论、用户或文章"
          @keyup.enter="handleSearch"
        />
        <input
          v-model="articleId"
          class="input article-id-input"
          type="number"
          min="1"
          placeholder="文章 ID"
          @keyup.enter="handleSearch"
        />
        <select v-model="status" class="input filter-select" @change="handleSearch">
          <option value="">全部状态</option>
          <option value="1">已显示</option>
          <option value="0">已隐藏</option>
        </select>
        <button class="btn btn-primary" @click="handleSearch">查询</button>
      </div>
    </div>

    <div class="card">
      <div v-if="loading" class="loading-state">加载中...</div>
      <div v-else-if="comments.length === 0" class="empty-state">暂无评论</div>
      <div v-else class="table-wrap">
        <table class="data-table">
          <thead>
            <tr>
              <th>评论内容</th>
              <th>评论者</th>
              <th>所属文章</th>
              <th>状态</th>
              <th>发布时间</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="comment in comments" :key="comment.id">
              <td class="content-cell">{{ comment.content }}</td>
              <td>
                <div>{{ comment.nickname || comment.username || '-' }}</div>
                <div class="muted-text">@{{ comment.username || '-' }}</div>
              </td>
              <td>
                <div class="article-title-cell">{{ comment.articleTitle || '未找到文章' }}</div>
                <div class="muted-text">#{{ comment.articleId }}</div>
              </td>
              <td>
                <span :class="['status-tag', comment.status === 1 ? 'status-enabled' : 'status-disabled']">
                  {{ comment.status === 1 ? '已显示' : '已隐藏' }}
                </span>
              </td>
              <td>{{ formatDate(comment.createdAt) }}</td>
              <td class="actions-cell">
                <button
                  class="btn btn-outline btn-sm"
                  :disabled="savingId === comment.id"
                  @click="toggleStatus(comment)"
                >
                  {{ comment.status === 1 ? '隐藏' : '显示' }}
                </button>
                <button
                  class="btn btn-danger btn-sm"
                  :disabled="deletingId === comment.id"
                  @click="removeComment(comment)"
                >
                  删除
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-if="total > pageSize" class="pagination-wrap">
      <button class="btn btn-outline" :disabled="page <= 1" @click="changePage(page - 1)">上一页</button>
      <span class="page-info">{{ page }} / {{ Math.ceil(total / pageSize) }}</span>
      <button
        class="btn btn-outline"
        :disabled="page >= Math.ceil(total / pageSize)"
        @click="changePage(page + 1)"
      >
        下一页
      </button>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { deleteAdminComment, getAdminComments, updateCommentStatus } from '../../api/comment'
import { message } from '../../utils/message'

const comments = ref([])
const loading = ref(false)
const savingId = ref(null)
const deletingId = ref(null)
const keyword = ref('')
const articleId = ref('')
const status = ref('')
const page = ref(1)
const pageSize = 10
const total = ref(0)

// 将原始数据格式化为界面展示内容。
function formatDate(dateStr) {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

// 处理用户操作或浏览器事件。
function handleSearch() {
  page.value = 1
  fetchComments()
}

// 更新当前筛选条件或页面状态。
function changePage(nextPage) {
  page.value = nextPage
  fetchComments()
}

// 加载当前页面所需的数据。
async function fetchComments() {
  loading.value = true
  try {
    const params = { page: page.value, pageSize }
    if (keyword.value.trim()) params.keyword = keyword.value.trim()
    if (articleId.value) params.articleId = Number(articleId.value)
    if (status.value !== '') params.status = Number(status.value)

    const res = await getAdminComments(params)
    comments.value = res.data?.list || []
    total.value = res.data?.total || 0
  } finally {
    loading.value = false
  }
}

// 切换对应的界面状态。
async function toggleStatus(comment) {
  savingId.value = comment.id
  try {
    const nextStatus = comment.status === 1 ? 0 : 1
    await updateCommentStatus(comment.id, { status: nextStatus })
    message.success(nextStatus === 1 ? '评论已显示' : '评论已隐藏')
    await fetchComments()
  } finally {
    savingId.value = null
  }
}

// 处理当前模块的相关逻辑。
async function removeComment(comment) {
  if (!window.confirm('确定删除这条评论吗？删除后不可恢复。')) return

  deletingId.value = comment.id
  try {
    await deleteAdminComment(comment.id)
    message.success('评论已删除')
    if (comments.value.length === 1 && page.value > 1) page.value -= 1
    await fetchComments()
  } finally {
    deletingId.value = null
  }
}

onMounted(fetchComments)
</script>

<style scoped>
.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  gap: 16px;
}

.page-subtitle {
  font-size: 16px;
  font-weight: 600;
}

.filter-row {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.search-input {
  min-width: 240px;
}

.article-id-input {
  width: 110px;
}

.filter-select {
  width: 120px;
}

.table-wrap {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  min-width: 920px;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid var(--border);
  font-size: 14px;
  vertical-align: top;
}

.data-table th {
  color: var(--text-muted);
  font-weight: 500;
}

.content-cell {
  max-width: 280px;
  white-space: normal;
  line-height: 1.6;
}

.article-title-cell {
  max-width: 220px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.muted-text {
  margin-top: 4px;
  color: var(--text-muted);
  font-size: 12px;
}

.status-tag {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 4px;
  font-size: 12px;
}

.status-enabled {
  background: rgba(34, 197, 94, 0.15);
  color: #22c55e;
}

.status-disabled {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}

.actions-cell {
  white-space: nowrap;
}

.actions-cell .btn + .btn {
  margin-left: 8px;
}

.btn-sm {
  padding: 6px 12px;
  font-size: 13px;
}

.btn-danger {
  border: 1px solid rgba(239, 68, 68, 0.4);
  background: transparent;
  color: #f87171;
}

.btn-danger:hover {
  background: rgba(239, 68, 68, 0.12);
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
