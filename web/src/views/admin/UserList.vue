<template>
  <div class="user-page">
    <div class="action-bar">
      <h3 class="page-subtitle">用户管理</h3>
      <div class="filter-row">
        <input
          v-model="keyword"
          class="input search-input"
          type="text"
          placeholder="搜索用户名或昵称"
          @keyup.enter="handleSearch"
        />
        <select v-model="role" class="input filter-select" @change="handleSearch">
          <option value="">全部角色</option>
          <option value="admin">管理员</option>
          <option value="user">普通用户</option>
        </select>
        <select v-model="status" class="input filter-select" @change="handleSearch">
          <option value="">全部状态</option>
          <option value="1">启用</option>
          <option value="0">禁用</option>
        </select>
        <button class="btn btn-primary" @click="handleSearch">查询</button>
      </div>
    </div>

    <div class="card">
      <div v-if="loading" class="loading-state">加载中...</div>
      <div v-else-if="users.length === 0" class="empty-state">暂无用户</div>
      <div v-else class="table-wrap">
        <table class="data-table">
          <thead>
            <tr>
              <th>用户名</th>
              <th>昵称</th>
              <th>角色</th>
              <th>状态</th>
              <th>创建时间</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in users" :key="user.id">
              <td>{{ user.username }}</td>
              <td>{{ user.nickname }}</td>
              <td>
                <span class="tag">{{ user.role === 'admin' ? '管理员' : '普通用户' }}</span>
              </td>
              <td>
                <span :class="['status-tag', user.status === 1 ? 'status-enabled' : 'status-disabled']">
                  {{ user.status === 1 ? '启用' : '禁用' }}
                </span>
              </td>
              <td>{{ formatDate(user.createdAt) }}</td>
              <td>
                <button
                  class="btn btn-outline btn-sm"
                  :disabled="savingId === user.id"
                  @click="toggleStatus(user)"
                >
                  {{ user.status === 1 ? '禁用' : '启用' }}
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
      <button class="btn btn-outline" :disabled="page >= Math.ceil(total / pageSize)" @click="changePage(page + 1)">下一页</button>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { getAdminUsers, updateUserStatus } from '../../api/user'
import { message } from '../../utils/message'

const users = ref([])
const loading = ref(false)
const savingId = ref(null)
const keyword = ref('')
const role = ref('')
const status = ref('')
const page = ref(1)
const pageSize = 10
const total = ref(0)

function formatDate(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

function handleSearch() {
  page.value = 1
  fetchUsers()
}

function changePage(nextPage) {
  page.value = nextPage
  fetchUsers()
}

async function fetchUsers() {
  loading.value = true
  try {
    const params = { page: page.value, pageSize }
    if (keyword.value.trim()) params.keyword = keyword.value.trim()
    if (role.value) params.role = role.value
    if (status.value !== '') params.status = Number(status.value)

    const res = await getAdminUsers(params)
    users.value = res.data?.list || []
    total.value = res.data?.total || 0
  } finally {
    loading.value = false
  }
}

async function toggleStatus(user) {
  savingId.value = user.id
  try {
    const nextStatus = user.status === 1 ? 0 : 1
    await updateUserStatus(user.id, { status: nextStatus })
    message.success(nextStatus === 1 ? '用户已启用' : '用户已禁用')
    await fetchUsers()
  } finally {
    savingId.value = null
  }
}

onMounted(fetchUsers)
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

.filter-select {
  width: 140px;
}

.table-wrap {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid var(--border);
  font-size: 14px;
}

.data-table th {
  color: var(--text-muted);
  font-weight: 500;
}

.tag {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 4px;
  background: var(--accent-dim);
  color: var(--accent);
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

.btn-sm {
  padding: 6px 12px;
  font-size: 13px;
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
