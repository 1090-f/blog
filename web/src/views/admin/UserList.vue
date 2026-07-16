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
        <button class="btn btn-primary" @click="openCreateDialog">创建管理员</button>
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
                <select
                  class="input role-select"
                  :value="user.role"
                  :disabled="savingId === user.id || isProtectedAdmin(user)"
                  @change="handleRoleChange(user, $event)"
                >
                  <option value="admin">管理员</option>
                  <option value="user">普通用户</option>
                </select>
                <span v-if="isCurrentUser(user)" class="current-label">当前账号</span>
                <span v-else-if="user.role === 'admin'" class="current-label">受保护</span>
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
                  :disabled="savingId === user.id || isProtectedAdmin(user)"
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

    <div v-if="createDialogVisible" class="modal-overlay" @click.self="closeCreateDialog">
      <div class="modal-card card">
        <h3>创建管理员</h3>
        <p class="modal-tip">新账号将直接获得后台管理权限，请仅为可信人员创建。</p>
        <div class="form-group">
          <label class="form-label">用户名</label>
          <input v-model="createForm.username" class="input" maxlength="50" placeholder="3-50 位字符" />
        </div>
        <div class="form-group">
          <label class="form-label">昵称</label>
          <input v-model="createForm.nickname" class="input" maxlength="50" placeholder="2-50 位字符" />
        </div>
        <div class="form-group">
          <label class="form-label">密码</label>
          <input v-model="createForm.password" class="input" type="password" maxlength="32" placeholder="至少 6 位" />
        </div>
        <div class="form-group">
          <label class="form-label">确认密码</label>
          <input v-model="createForm.confirmPassword" class="input" type="password" maxlength="32" placeholder="再次输入密码" />
        </div>
        <div class="modal-actions">
          <button class="btn btn-outline" :disabled="creating" @click="closeCreateDialog">取消</button>
          <button class="btn btn-primary" :disabled="creating" @click="handleCreateAdmin">
            {{ creating ? '创建中...' : '确认创建' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { createAdminUser, getAdminUsers, updateUserRole, updateUserStatus } from '../../api/user'
import { useUserStore } from '../../stores/user'
import { message } from '../../utils/message'

const userStore = useUserStore()

const users = ref([])
const loading = ref(false)
const savingId = ref(null)
const keyword = ref('')
const role = ref('')
const status = ref('')
const page = ref(1)
const pageSize = 10
const total = ref(0)
const createDialogVisible = ref(false)
const creating = ref(false)
const createForm = ref({ username: '', nickname: '', password: '', confirmPassword: '' })

function isCurrentUser(user) {
  return user.id === userStore.user?.id
}

function isProtectedAdmin(user) {
  return user.role === 'admin'
}

function openCreateDialog() {
  createForm.value = { username: '', nickname: '', password: '', confirmPassword: '' }
  createDialogVisible.value = true
}

function closeCreateDialog() {
  if (!creating.value) createDialogVisible.value = false
}

async function handleCreateAdmin() {
  const form = createForm.value
  if (!form.username.trim() || !form.nickname.trim() || !form.password) {
    message.warning('请填写完整的管理员信息')
    return
  }
  if (form.password !== form.confirmPassword) {
    message.warning('两次输入的密码不一致')
    return
  }
  creating.value = true
  try {
    await createAdminUser({ username: form.username.trim(), nickname: form.nickname.trim(), password: form.password })
    message.success('管理员创建成功')
    createDialogVisible.value = false
    page.value = 1
    await fetchUsers()
  } finally {
    creating.value = false
  }
}

// 将原始数据格式化为界面展示内容。
function formatDate(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

// 处理用户操作或浏览器事件。
function handleSearch() {
  page.value = 1
  fetchUsers()
}

// 更新当前筛选条件或页面状态。
function changePage(nextPage) {
  page.value = nextPage
  fetchUsers()
}

// 加载当前页面所需的数据。
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

// 切换对应的界面状态。
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

async function handleRoleChange(user, event) {
  const nextRole = event.target.value
  if (nextRole === user.role) return
  savingId.value = user.id
  try {
    await updateUserRole(user.id, { role: nextRole })
    message.success(nextRole === 'admin' ? '已设为管理员' : '已设为普通用户')
    await fetchUsers()
  } catch {
    event.target.value = user.role
    // 请求层会统一展示错误消息，此处只恢复下拉框原值。
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

.role-select {
  width: 112px;
  padding: 5px 8px;
  font-size: 13px;
}

.current-label {
  display: inline-block;
  margin-left: 8px;
  color: var(--text-muted);
  font-size: 12px;
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

.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  background: rgba(0, 0, 0, 0.55);
}

.modal-card {
  width: min(440px, 100%);
  padding: 24px;
}

.modal-card h3 {
  margin-bottom: 8px;
}

.modal-tip {
  margin-bottom: 20px;
  color: var(--text-muted);
  font-size: 13px;
  line-height: 1.6;
}

.form-group {
  margin-bottom: 16px;
}

.form-label {
  display: block;
  margin-bottom: 6px;
  color: var(--text-secondary);
  font-size: 14px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}

@media (max-width: 960px) {
  .action-bar {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
