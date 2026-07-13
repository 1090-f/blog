<template>
  <div class="category-page">
    <div class="action-bar">
      <h3 class="page-subtitle">分类列表</h3>
      <button class="btn btn-primary" @click="openDialog()">新增分类</button>
    </div>

    <div class="card">
      <div v-if="loading" class="loading-state">加载中...</div>
      <div v-else-if="categories.length > 0" class="category-list">
        <div v-for="cat in categories" :key="cat.id" class="category-item">
          <div class="category-info">
            <h4 class="category-name">{{ cat.name }}</h4>
            <p class="category-desc">{{ cat.description || '暂无描述' }}</p>
          </div>
          <div class="category-actions">
            <button class="btn btn-outline btn-sm" @click="openDialog(cat)">编辑</button>
            <button class="btn btn-danger btn-sm" @click="handleDelete(cat.id)">删除</button>
          </div>
        </div>
      </div>
      <div v-else class="empty-state">暂无分类</div>
    </div>

    <div v-if="dialogVisible" class="modal-overlay" @click.self="dialogVisible = false">
      <div class="modal-card card">
        <h3 class="modal-title">{{ editingId ? '编辑分类' : '新增分类' }}</h3>
        <div class="form-group">
          <label class="form-label">名称</label>
          <input v-model="form.name" class="input" type="text" placeholder="请输入分类名称" maxlength="50" />
        </div>
        <div class="form-group">
          <label class="form-label">描述</label>
          <textarea
            v-model="form.description"
            class="input"
            rows="3"
            placeholder="请输入分类描述，可留空"
            maxlength="255"
          ></textarea>
        </div>
        <div class="modal-actions">
          <button class="btn btn-outline" @click="dialogVisible = false">取消</button>
          <button class="btn btn-primary" :disabled="saving" @click="handleSave">
            {{ saving ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { createCategory, deleteCategory, getCategories, updateCategory } from '../../api/category'
import { message } from '../../utils/message'

const categories = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const editingId = ref(null)
const form = ref({ name: '', description: '' })
const saving = ref(false)

function openDialog(row) {
  if (row) {
    editingId.value = row.id
    form.value = { name: row.name, description: row.description || '' }
  } else {
    editingId.value = null
    form.value = { name: '', description: '' }
  }
  dialogVisible.value = true
}

async function fetchData() {
  loading.value = true
  try {
    const res = await getCategories()
    categories.value = res.data || []
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  if (!form.value.name.trim()) {
    message.warning('请输入分类名称')
    return
  }

  saving.value = true
  try {
    if (editingId.value) {
      await updateCategory(editingId.value, form.value)
      message.success('分类已更新')
    } else {
      await createCategory(form.value)
      message.success('分类已创建')
    }
    dialogVisible.value = false
    await fetchData()
  } finally {
    saving.value = false
  }
}

async function handleDelete(id) {
  if (!window.confirm('确认删除这个分类吗？')) return
  try {
    await deleteCategory(id)
    message.success('删除成功')
    await fetchData()
  } catch (error) {
    message.error('删除失败')
  }
}

onMounted(fetchData)
</script>

<style scoped>
.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-subtitle {
  font-size: 16px;
  font-weight: 600;
}

.category-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 16px 0;
  border-bottom: 1px solid var(--border);
}

.category-item:last-child {
  border-bottom: none;
}

.category-name {
  font-size: 15px;
  font-weight: 600;
  margin-bottom: 4px;
}

.category-desc {
  color: var(--text-muted);
  font-size: 13px;
}

.category-actions {
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

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-card {
  width: 100%;
  max-width: 480px;
  padding: 32px;
}

.modal-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 24px;
}

.form-group {
  margin-bottom: 20px;
}

.form-label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  color: var(--text-secondary);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 20px;
  border-top: 1px solid var(--border);
}
</style>
