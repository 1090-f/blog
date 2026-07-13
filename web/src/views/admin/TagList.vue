<template>
  <div class="tag-page">
    <div class="action-bar">
      <h3 class="page-subtitle">标签管理</h3>
      <button class="btn btn-primary" @click="openDialog()">新增标签</button>
    </div>

    <div class="card">
      <div v-if="loading" class="loading-state">加载中...</div>
      <div v-else-if="tags.length === 0" class="empty-state">暂无标签</div>
      <div v-else class="tag-list">
        <div v-for="tag in tags" :key="tag.id" class="tag-item">
          <div class="tag-name">{{ tag.name }}</div>
          <div class="tag-actions">
            <button class="btn btn-outline btn-sm" @click="openDialog(tag)">编辑</button>
            <button class="btn btn-danger btn-sm" @click="handleDelete(tag.id)">删除</button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="dialogVisible" class="modal-overlay" @click.self="dialogVisible = false">
      <div class="modal-card card">
        <h3 class="modal-title">{{ editingId ? '编辑标签' : '新增标签' }}</h3>
        <div class="form-group">
          <label class="form-label">名称</label>
          <input v-model="form.name" class="input" type="text" placeholder="请输入标签名称" maxlength="50" />
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
import { createTag, deleteTag, getTags, updateTag } from '../../api/tag'
import { message } from '../../utils/message'

const tags = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const editingId = ref(null)
const form = ref({ name: '' })
const saving = ref(false)

function openDialog(tag) {
  editingId.value = tag?.id || null
  form.value = { name: tag?.name || '' }
  dialogVisible.value = true
}

async function fetchData() {
  loading.value = true
  try {
    const res = await getTags()
    tags.value = res.data || []
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  if (!form.value.name.trim()) {
    message.warning('请输入标签名称')
    return
  }

  saving.value = true
  try {
    if (editingId.value) {
      await updateTag(editingId.value, form.value)
      message.success('标签已更新')
    } else {
      await createTag(form.value)
      message.success('标签已创建')
    }
    dialogVisible.value = false
    await fetchData()
  } finally {
    saving.value = false
  }
}

async function handleDelete(id) {
  if (!window.confirm('确认删除这个标签吗？')) return
  try {
    await deleteTag(id)
    message.success('删除成功')
    await fetchData()
  } catch (error) {
    // The request layer already shows the backend message.
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

.tag-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 16px 0;
  border-bottom: 1px solid var(--border);
}

.tag-item:last-child {
  border-bottom: none;
}

.tag-name {
  display: inline-flex;
  padding: 5px 12px;
  border-radius: 20px;
  background: var(--accent-dim);
  color: var(--accent);
  font-size: 14px;
}

.tag-actions {
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
  max-width: 420px;
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
