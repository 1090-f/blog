<template>
  <div class="article-form-page">
    <div class="form-card card">
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
          <select v-model="form.categoryId" class="input">
            <option :value="null">选择分类</option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
          </select>
        </div>
        <div class="form-group flex-1">
          <label class="form-label">状态</label>
          <div class="radio-group">
            <label class="radio-item" :class="{ active: form.status === 'draft' }">
              <input v-model="form.status" type="radio" value="draft" />
              <span>草稿</span>
            </label>
            <label class="radio-item" :class="{ active: form.status === 'published' }">
              <input v-model="form.status" type="radio" value="published" />
              <span>发布</span>
            </label>
          </div>
        </div>
      </div>

      <div class="form-group">
        <label class="form-label">标签</label>
        <div v-if="tags.length" class="tag-picker">
          <label v-for="tag in tags" :key="tag.id" class="tag-option" :class="{ active: form.tagIds.includes(tag.id) }">
            <input v-model="form.tagIds" type="checkbox" :value="tag.id" />
            <span>{{ tag.name }}</span>
          </label>
        </div>
        <p v-else class="form-tip">暂无可选标签，请先在标签管理中创建。</p>
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
        <label class="form-label">封面图 URL</label>
        <div class="input-with-upload">
          <input
            v-model="form.coverImage"
            class="input"
            type="text"
            placeholder="https://example.com/cover.jpg"
          />
          <label class="upload-btn">
            上传图片
            <input type="file" accept="image/*" hidden @change="handleCoverUpload" />
          </label>
        </div>
        <p class="upload-tip">上传格式支持 jpg、jpeg、png、webp，大小受服务端配置限制。</p>
        <div v-if="form.coverImage" class="cover-preview">
          <img :src="form.coverImage" alt="封面预览" />
        </div>
      </div>

      <div class="form-group">
        <div class="content-label-row">
          <label class="form-label">内容</label>
          <label v-if="!isEdit" class="import-btn">
            导入文章
            <input type="file" accept=".md,.markdown,.txt" hidden @change="handleArticleImport" />
          </label>
        </div>
        <p v-if="!isEdit" class="import-tip">支持导入 Markdown 或文本文件，文件内容会覆盖当前编辑器内容。</p>
        <div class="editor-tabs">
          <button :class="['tab-btn', { active: editorTab === 'edit' }]" @click="editorTab = 'edit'">编辑</button>
          <button :class="['tab-btn', { active: editorTab === 'preview' }]" @click="editorTab = 'preview'">预览</button>
        </div>
        <div v-show="editorTab === 'edit'" class="editor-wrap">
          <div class="editor-toolbar">
            <button type="button" title="加粗" @click="insertMarkdown('**', '**')">B</button>
            <button type="button" title="斜体" @click="insertMarkdown('*', '*')"><i>I</i></button>
            <button type="button" title="行内代码" @click="insertMarkdown('`', '`')">&lt;/&gt;</button>
            <button type="button" title="代码块" @click="insertMarkdown('```\n', '\n```')">{ }</button>
            <button type="button" title="链接" @click="insertMarkdown('[', '](url)')">Link</button>
            <button type="button" title="图片" @click="insertMarkdown('![alt](', ')')">Img</button>
            <button type="button" title="二级标题" @click="insertMarkdown('## ', '')">H2</button>
            <button type="button" title="列表" @click="insertMarkdown('- ', '')">List</button>
          </div>
          <textarea
            ref="editorRef"
            v-model="form.content"
            class="input editor-textarea"
            rows="20"
            placeholder="支持 Markdown 格式"
          ></textarea>
        </div>
        <div v-show="editorTab === 'preview'" class="preview-wrap">
          <div class="preview-content" v-html="previewHtml"></div>
        </div>
      </div>

      <div class="form-actions">
        <button class="btn btn-primary" :disabled="saving" @click="handleSave">
          {{ saving ? '保存中...' : (isEdit ? '更新文章' : '创建文章') }}
        </button>
        <button class="btn btn-outline" @click="$router.back()">取消</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { marked } from 'marked'
import { createAdminArticle, getArticleDetail, updateArticle, uploadFile } from '../../api/article'
import { getCategories } from '../../api/category'
import { getTags } from '../../api/tag'
import { message } from '../../utils/message'

const route = useRoute()
const router = useRouter()
const isEdit = computed(() => !!route.params.id)
const editorRef = ref(null)
const editorTab = ref('edit')

const form = ref({
  title: '',
  categoryId: null,
  summary: '',
  coverImage: '',
  content: '',
  status: 'draft',
  tagIds: []
})
const categories = ref([])
const tags = ref([])
const saving = ref(false)
const allowedArticleExtensions = ['md', 'markdown', 'txt']

const previewHtml = computed(() => {
  try {
    return marked(form.value.content || '')
  } catch (error) {
    return form.value.content
  }
})

function insertMarkdown(before, after) {
  const textarea = editorRef.value
  if (!textarea) return

  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const text = form.value.content
  const selected = text.substring(start, end)
  form.value.content = text.substring(0, start) + before + selected + after + text.substring(end)
}

async function handleCoverUpload(event) {
  const file = event.target.files?.[0]
  if (!file) return

  try {
    const res = await uploadFile(file)
    form.value.coverImage = res.data?.url || res.data || ''
    message.success('上传成功')
  } catch (error) {
    message.error('上传失败')
  } finally {
    event.target.value = ''
  }
}

async function handleArticleImport(event) {
  const file = event.target.files?.[0]
  if (!file) return

  const extension = file.name.split('.').pop()?.toLowerCase()
  if (!allowedArticleExtensions.includes(extension)) {
    message.warning('仅支持导入 Markdown（.md、.markdown）或纯文本（.txt）文件')
    event.target.value = ''
    return
  }

  try {
    const content = await file.text()
    form.value.content = content

    if (!form.value.title.trim()) {
      const heading = content.match(/^#\s+(.+)$/m)?.[1]?.trim()
      form.value.title = heading || file.name.replace(/\.(md|markdown|txt)$/i, '')
    }

    editorTab.value = 'edit'
    message.success('文章已导入')
  } catch (error) {
    message.error('文章导入失败')
  } finally {
    event.target.value = ''
  }
}

async function handleSave() {
  if (!form.value.title.trim() || !form.value.categoryId || !form.value.content.trim()) {
    message.warning('请填写标题、分类和内容')
    return
  }

  saving.value = true
  try {
    if (isEdit.value) {
      await updateArticle(route.params.id, form.value)
      message.success('文章已更新')
    } else {
      await createAdminArticle(form.value)
      message.success('文章已创建')
    }
    router.push('/admin/articles')
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  const [categoryRes, tagRes] = await Promise.all([getCategories(), getTags()])
  categories.value = categoryRes.data || []
  tags.value = tagRes.data || []

  if (!isEdit.value) return

  const res = await getArticleDetail(route.params.id)
  const data = res.data
  form.value = {
    title: data.title,
    categoryId: data.categoryId,
    summary: data.summary,
    coverImage: data.coverImage,
    content: data.content,
    status: data.status,
    tagIds: (data.tags || []).map(tag => tag.id)
  }
})
</script>

<style scoped>
.article-form-page {
  width: 100%;
}

.form-card {
  width: 100%;
  max-width: none;
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

.content-label-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.content-label-row .form-label {
  margin-bottom: 8px;
}

.import-btn {
  margin-bottom: 8px;
  color: var(--accent);
  cursor: pointer;
  font-size: 13px;
}

.import-tip {
  margin: -4px 0 12px;
  color: var(--text-muted);
  font-size: 12px;
}

.form-row {
  display: flex;
  gap: 16px;
}

.tag-picker {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag-option {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 7px 12px;
  border: 1px solid var(--border);
  border-radius: 20px;
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 13px;
}

.tag-option input {
  accent-color: var(--accent);
}

.tag-option.active {
  border-color: var(--accent);
  background: var(--accent-dim);
  color: var(--accent);
}

.form-tip {
  color: var(--text-muted);
  font-size: 13px;
}

.flex-1 {
  flex: 1;
}

.radio-group {
  display: flex;
  gap: 12px;
}

.radio-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.2s;
}

.radio-item input {
  display: none;
}

.radio-item.active {
  background: var(--accent-dim);
  border-color: var(--accent);
  color: var(--accent);
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

.upload-tip {
  margin-top: 10px;
  color: var(--text-muted);
  font-size: 12px;
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

.editor-tabs {
  display: flex;
  gap: 4px;
  margin-bottom: 12px;
}

.tab-btn {
  padding: 8px 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 13px;
  transition: all 0.2s;
}

.tab-btn.active {
  background: var(--accent-dim);
  border-color: var(--accent);
  color: var(--accent);
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

.editor-textarea {
  border-radius: 0 0 var(--radius-sm) var(--radius-sm);
  font-family: 'Menlo', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.6;
  min-height: 400px;
  resize: vertical;
}

.preview-wrap {
  min-height: 400px;
  padding: 20px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
}

.preview-content {
  line-height: 1.8;
  color: var(--text-primary);
}

.preview-content :deep(h1),
.preview-content :deep(h2),
.preview-content :deep(h3) {
  margin-top: 20px;
  margin-bottom: 12px;
}

.preview-content :deep(p) {
  margin-bottom: 12px;
}

.preview-content :deep(code) {
  background: rgba(0, 0, 0, 0.3);
  padding: 2px 6px;
  border-radius: 4px;
}

.preview-content :deep(pre) {
  background: rgba(0, 0, 0, 0.4);
  padding: 16px;
  border-radius: var(--radius-sm);
  overflow-x: auto;
}

.form-actions {
  display: flex;
  gap: 12px;
  padding-top: 20px;
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

@media (max-width: 768px) {
  .form-row,
  .input-with-upload,
  .form-actions {
    flex-direction: column;
  }
}
</style>
