<template>
  <div class="write-page">
    <div class="page-banner">
      <h1 class="page-title">发布文章</h1>
      <p class="page-desc">把你的想法整理成内容，发布给更多人看到。</p>
    </div>

    <div class="write-form card">
      <div class="form-group">
        <label class="form-label">标题</label>
        <input
          v-model="form.title"
          class="input"
          type="text"
          placeholder="请输入文章标题"
          maxlength="150"
        />
      </div>

      <div class="form-row">
        <div class="form-group flex-1">
          <label class="form-label">分类</label>
          <select v-model="form.categoryId" class="input category-select" @change="handleCategoryChange">
            <option :value="null">选择分类</option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
            <option value="custom">+ 自定义分类</option>
          </select>
          <input
            v-if="showCustomCategory"
            v-model="form.customCategory"
            class="input custom-category-input"
            type="text"
            placeholder="输入新分类名称"
            maxlength="50"
          />
        </div>
      </div>

      <div class="form-group">
        <label class="form-label">摘要</label>
        <textarea
          v-model="form.summary"
          class="input"
          rows="3"
          placeholder="请输入文章摘要，可留空"
          maxlength="255"
        ></textarea>
      </div>

      <div class="form-group">
        <label class="form-label">封面图</label>
        <div class="input-with-upload">
          <input
            v-model="form.coverImage"
            class="input"
            type="text"
            placeholder="https://example.com/image.jpg"
          />
          <label class="upload-btn">
            上传
            <input type="file" accept="image/*" hidden @change="handleCoverUpload" />
          </label>
        </div>
        <div v-if="form.coverImage" class="cover-preview">
          <img :src="form.coverImage" alt="封面预览" />
        </div>
      </div>

      <div class="form-group">
        <label class="form-label">内容</label>
        <div class="editor-toolbar">
          <button type="button" title="加粗" @click="insertMarkdown('**', '**')">B</button>
          <button type="button" title="斜体" @click="insertMarkdown('*', '*')"><i>I</i></button>
          <button type="button" title="行内代码" @click="insertMarkdown('`', '`')">&lt;/&gt;</button>
          <button type="button" title="代码块" @click="insertMarkdown('```\n', '\n```')">{ }</button>
          <button type="button" title="链接" @click="insertMarkdown('[', '](url)')">Link</button>
          <button type="button" title="图片" @click="insertImage">Img</button>
          <button type="button" title="二级标题" @click="insertMarkdown('## ', '')">H2</button>
          <button type="button" title="列表" @click="insertMarkdown('- ', '')">List</button>
        </div>
        <textarea
          ref="editorRef"
          v-model="form.content"
          class="input content-textarea"
          rows="15"
          placeholder="支持 Markdown 格式"
        ></textarea>
      </div>

      <div class="form-actions">
        <button class="btn btn-primary" :disabled="saving" @click="handleSave">
          {{ saving ? '发布中...' : '发布文章' }}
        </button>
        <button class="btn btn-outline" @click="$router.back()">取消</button>
      </div>
    </div>

    <div v-if="showImageUpload" class="modal-overlay" @click.self="showImageUpload = false">
      <div class="modal-card card">
        <h3 class="modal-title">上传图片</h3>
        <div class="upload-area">
          <input ref="imageInput" type="file" accept="image/*" hidden @change="handleImageUpload" />
          <button class="btn btn-outline" @click="openImagePicker">选择图片</button>
          <p class="upload-hint">支持 jpg、jpeg、png、webp，大小受服务端限制。</p>
        </div>
        <div class="modal-actions">
          <button class="btn btn-outline" @click="showImageUpload = false">关闭</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { nextTick, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { createArticle, uploadFile } from '../../api/article'
import { createUserCategory, getCategories } from '../../api/category'
import { message } from '../../utils/message'

const router = useRouter()
const editorRef = ref(null)
const imageInput = ref(null)
const showImageUpload = ref(false)

const form = ref({
  title: '',
  categoryId: null,
  customCategory: '',
  summary: '',
  coverImage: '',
  content: '',
  status: 'published'
})
const categories = ref([])
const saving = ref(false)
const showCustomCategory = ref(false)

function handleCategoryChange() {
  showCustomCategory.value = form.value.categoryId === 'custom'
  if (!showCustomCategory.value) {
    form.value.customCategory = ''
  }
}

async function createCustomCategory(name) {
  const res = await createUserCategory({ name, description: '' })
  return res.data?.id
}

function insertMarkdown(before, after) {
  const textarea = editorRef.value
  if (!textarea) return

  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const text = form.value.content
  const selected = text.substring(start, end)
  form.value.content = text.substring(0, start) + before + selected + after + text.substring(end)

  nextTick(() => {
    textarea.focus()
    const cursorStart = start + before.length
    const cursorEnd = selected ? cursorStart + selected.length : cursorStart
    textarea.setSelectionRange(cursorStart, cursorEnd)
  })
}

function insertImage() {
  showImageUpload.value = true
}

function openImagePicker() {
  imageInput.value?.click()
}

async function handleCoverUpload(event) {
  const file = event.target.files?.[0]
  if (!file) return

  try {
    const res = await uploadFile(file)
    form.value.coverImage = res.data?.url || res.data
    message.success('上传成功')
  } catch (error) {
    message.error('上传失败')
  } finally {
    event.target.value = ''
  }
}

async function handleImageUpload(event) {
  const file = event.target.files?.[0]
  if (!file) return

  try {
    const res = await uploadFile(file)
    const url = res.data?.url || res.data
    insertMarkdown(`![图片](${url})`, '')
    showImageUpload.value = false
    message.success('图片已插入')
  } catch (error) {
    message.error('上传失败')
  } finally {
    event.target.value = ''
  }
}

async function handleSave() {
  if (!form.value.title.trim() || !form.value.content.trim()) {
    message.warning('请填写标题和内容')
    return
  }

  let categoryId = form.value.categoryId

  if (categoryId === 'custom') {
    if (!form.value.customCategory.trim()) {
      message.warning('请输入自定义分类名称')
      return
    }

    try {
      categoryId = await createCustomCategory(form.value.customCategory.trim())
      const res = await getCategories()
      categories.value = res.data || []
    } catch (error) {
      message.error('创建分类失败')
      return
    }
  }

  if (!categoryId) {
    message.warning('请选择分类')
    return
  }

  saving.value = true
  try {
    await createArticle({
      ...form.value,
      categoryId
    })
    message.success('文章发布成功')
      router.push('/archive')
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  try {
    const res = await getCategories()
    categories.value = res.data || []
  } catch (error) {
    categories.value = []
  }
})
</script>

<style scoped>
.write-page {
  max-width: 800px;
  margin: 0 auto;
}

.page-banner {
  text-align: center;
  padding: 40px 0 32px;
}

.page-title {
  font-size: 32px;
  font-weight: 700;
  margin-bottom: 8px;
}

.page-desc {
  color: var(--text-muted);
  font-size: 15px;
}

.write-form {
  padding: 32px;
}

.form-group {
  margin-bottom: 24px;
}

.form-label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
}

.form-row {
  display: flex;
  gap: 16px;
}

.flex-1 {
  flex: 1;
}

.input-with-upload {
  display: flex;
  gap: 8px;
}

.input-with-upload .input {
  flex: 1;
}

.upload-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
  white-space: nowrap;
}

.upload-btn:hover {
  border-color: var(--accent);
  color: var(--accent);
}

.cover-preview {
  margin-top: 12px;
}

.cover-preview img {
  max-width: 200px;
  max-height: 120px;
  border-radius: var(--radius-sm);
  object-fit: cover;
}

.editor-toolbar {
  display: flex;
  gap: 4px;
  padding: 8px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-bottom: none;
  border-radius: var(--radius-sm) var(--radius-sm) 0 0;
}

.editor-toolbar button {
  min-width: 40px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: 1px solid transparent;
  border-radius: 4px;
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 13px;
  transition: all 0.2s;
}

.editor-toolbar button:hover {
  background: var(--bg-card);
  border-color: var(--border);
  color: var(--text-primary);
}

.content-textarea {
  border-radius: 0 0 var(--radius-sm) var(--radius-sm);
  font-family: 'Menlo', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.6;
  min-height: 300px;
  resize: vertical;
}

.form-actions {
  display: flex;
  gap: 12px;
  padding-top: 16px;
  border-top: 1px solid var(--border);
}

select.input {
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' fill='%239aa0a6' viewBox='0 0 16 16'%3E%3Cpath d='M8 11L3 6h10l-5 5z'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 12px center;
  padding-right: 36px;
}

select.input option {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.category-select {
  min-width: 150px;
}

.custom-category-input {
  margin-top: 12px;
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
  max-width: 400px;
  padding: 32px;
}

.modal-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 24px;
}

.upload-area {
  text-align: center;
  padding: 20px;
  border: 2px dashed var(--border);
  border-radius: var(--radius-sm);
  margin-bottom: 20px;
}

.upload-hint {
  color: var(--text-muted);
  font-size: 13px;
  margin-top: 12px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

@media (max-width: 768px) {
  .input-with-upload,
  .form-actions {
    flex-direction: column;
  }
}
</style>
