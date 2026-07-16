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
        <div class="editor-wrap" :class="{ fullscreen: isEditorFullscreen }">
          <div class="editor-toolbar">
            <div class="toolbar-group">
              <button class="toolbar-action" type="button" title="撤销" :disabled="!canUndo" @click="undoEditor">
                <span class="toolbar-icon">↶</span><span>撤销</span>
              </button>
              <button class="toolbar-action" type="button" title="重做" :disabled="!canRedo" @click="redoEditor">
                <span class="toolbar-icon">↷</span><span>重做</span>
              </button>
            </div>
            <span class="toolbar-divider" aria-hidden="true"></span>
            <div class="toolbar-group">
              <div class="toolbar-menu-wrap">
                <button class="toolbar-action" type="button" @click="toggleToolbarMenu('heading')">
                  <span class="toolbar-icon toolbar-icon-format">H</span><span>格式</span><span class="toolbar-caret">▾</span>
                </button>
                <div v-if="activeToolbarMenu === 'heading'" class="toolbar-menu toolbar-menu-list">
                  <button v-for="level in 6" :key="level" type="button" @click="applyHeading(level)">H{{ level }} 标题</button>
                </div>
              </div>
              <button class="toolbar-action" type="button" title="加粗" @click="applyInlineFormat('**', '**')">
                <span class="toolbar-icon toolbar-icon-bold">B</span><span>加粗</span>
              </button>
              <button class="toolbar-action" type="button" title="斜体" @click="applyInlineFormat('*', '*')">
                <span class="toolbar-icon toolbar-icon-italic">I</span><span>斜体</span>
              </button>
              <button class="toolbar-action" type="button" title="删除线" @click="applyInlineFormat('~~', '~~')">
                <span class="toolbar-icon"><s>S</s></span><span>删除线</span>
              </button>
              <div class="toolbar-menu-wrap">
                <button class="toolbar-action" type="button" @click="toggleToolbarMenu('color')">
                  <span class="toolbar-icon toolbar-icon-color">A</span><span>颜色</span><span class="toolbar-caret">▾</span>
                </button>
                <div v-if="activeToolbarMenu === 'color'" class="toolbar-menu color-menu">
                  <button v-for="color in colorOptions" :key="color.value" type="button" :title="color.label" :style="{ background: color.value }" @click="applyTextColor(color.value)"></button>
                </div>
              </div>
              <div class="toolbar-menu-wrap">
                <button class="toolbar-action" type="button" @click="toggleToolbarMenu('background')">
                  <span class="toolbar-icon toolbar-icon-background">A</span><span>背景</span><span class="toolbar-caret">▾</span>
                </button>
                <div v-if="activeToolbarMenu === 'background'" class="toolbar-menu color-menu">
                  <button v-for="color in backgroundOptions" :key="color.value" type="button" :title="color.label" :style="{ background: color.value }" @click="applyTextBackground(color.value)"></button>
                </div>
              </div>
            </div>
            <span class="toolbar-divider" aria-hidden="true"></span>
            <div class="toolbar-group">
              <div class="toolbar-menu-wrap">
                <button class="toolbar-action" type="button" @click="toggleToolbarMenu('list')">
                  <span class="toolbar-icon">☷</span><span>列表</span><span class="toolbar-caret">▾</span>
                </button>
                <div v-if="activeToolbarMenu === 'list'" class="toolbar-menu toolbar-menu-list">
                  <button type="button" @click="applyList('unordered')">无序列表</button>
                  <button type="button" @click="applyList('ordered')">有序列表</button>
                  <button type="button" @click="applyList('task')">任务列表</button>
                </div>
              </div>
              <div class="toolbar-menu-wrap">
                <button class="toolbar-action" type="button" @click="toggleToolbarMenu('align')">
                  <span class="toolbar-icon">≡</span><span>对齐</span><span class="toolbar-caret">▾</span>
                </button>
                <div v-if="activeToolbarMenu === 'align'" class="toolbar-menu toolbar-menu-list">
                  <button type="button" @click="applyAlignment('left')">左对齐</button>
                  <button type="button" @click="applyAlignment('center')">居中对齐</button>
                  <button type="button" @click="applyAlignment('right')">右对齐</button>
                </div>
              </div>
              <button class="toolbar-action" type="button" @click="insertSnippet('---')">
                <span class="toolbar-icon">—</span><span>水平线</span>
              </button>
              <button class="toolbar-action" type="button" @click="applyQuote">
                <span class="toolbar-icon">❝</span><span>块引用</span>
              </button>
            </div>
            <span class="toolbar-divider" aria-hidden="true"></span>
            <div class="toolbar-group">
              <div class="toolbar-menu-wrap">
                <button class="toolbar-action" type="button" @click="toggleToolbarMenu('code')">
                  <span class="toolbar-icon">&lt;/&gt;</span><span>代码</span><span class="toolbar-caret">▾</span>
                </button>
                <div v-if="activeToolbarMenu === 'code'" class="toolbar-menu toolbar-menu-list">
                  <button type="button" @click="applyInlineCode">行内代码</button>
                  <button type="button" @click="applyCodeBlock">代码块</button>
                </div>
              </div>
              <div class="toolbar-menu-wrap">
                <button class="toolbar-action" type="button" @click="toggleToolbarMenu('table')">
                  <span class="toolbar-icon">▦</span><span>表格</span><span class="toolbar-caret">▾</span>
                </button>
                <div v-if="activeToolbarMenu === 'table'" class="toolbar-menu table-settings-menu">
                  <label>
                    <span>列数</span>
                    <input v-model.number="tableSettings.columns" type="number" min="1" max="10" />
                  </label>
                  <label>
                    <span>内容行数</span>
                    <input v-model.number="tableSettings.rows" type="number" min="1" max="20" />
                  </label>
                  <label>
                    <span>对齐方式</span>
                    <select v-model="tableSettings.alignment">
                      <option value="left">左对齐</option>
                      <option value="center">居中</option>
                      <option value="right">右对齐</option>
                    </select>
                  </label>
                  <button class="table-insert-button" type="button" @click="insertCustomTable">插入表格</button>
                </div>
              </div>
              <div class="toolbar-menu-wrap">
                <button class="toolbar-action" type="button" @click="toggleToolbarMenu('image')">
                  <span class="toolbar-icon">▧</span><span>图像</span><span class="toolbar-caret">▾</span>
                </button>
                <div v-if="activeToolbarMenu === 'image'" class="toolbar-menu toolbar-menu-list">
                  <button type="button" @click="openContentImagePicker">上传本地图片</button>
                  <button type="button" @click="insertImageUrl">使用图片 URL</button>
                </div>
              </div>
              <button class="toolbar-action" type="button" @click="insertVideo">
                <span class="toolbar-icon">▶</span><span>视频</span>
              </button>
              <div class="toolbar-menu-wrap">
                <button class="toolbar-action" type="button" @click="toggleToolbarMenu('formula')">
                  <span class="toolbar-icon">∑</span><span>公式</span><span class="toolbar-caret">▾</span>
                </button>
                <div v-if="activeToolbarMenu === 'formula'" class="toolbar-menu toolbar-menu-list">
                  <button type="button" @click="insertFormula(false)">行内公式</button>
                  <button type="button" @click="insertFormula(true)">块级公式</button>
                </div>
              </div>
              <button class="toolbar-action" type="button" @click="insertLink">
                <span class="toolbar-icon">↗</span><span>链接</span>
              </button>
            </div>
            <span class="toolbar-divider" aria-hidden="true"></span>
            <div class="toolbar-group">
              <div class="toolbar-menu-wrap">
                <button class="toolbar-action" type="button" @click="toggleToolbarMenu('template')">
                  <span class="toolbar-icon">▤</span><span>模板</span><span class="toolbar-caret">▾</span>
                </button>
                <div v-if="activeToolbarMenu === 'template'" class="toolbar-menu toolbar-menu-list toolbar-menu-wide">
                  <button v-for="template in articleTemplates" :key="template.name" type="button" @click="insertArticleTemplate(template)">{{ template.name }}</button>
                </div>
              </div>
              <button class="toolbar-action" type="button" :class="{ active: editorPanel === 'catalog' }" @click="toggleEditorPanel('catalog')">
                <span class="toolbar-icon">☰</span><span>目录</span>
              </button>
              <button class="toolbar-action" type="button" :class="{ active: isEditorFullscreen }" @click="toggleEditorFullscreen">
                <span class="toolbar-icon">⛶</span><span>{{ isEditorFullscreen ? '退出全屏' : '全屏' }}</span>
              </button>
            </div>
          </div>
          <input
            ref="contentImageInputRef"
            class="content-image-input"
            type="file"
            accept="image/*"
            @change="handleContentImageUpload"
          />
          <div class="editor-workspace" :class="{ 'has-panel': editorPanel }">
            <div
              ref="editorRef"
              class="input editor-textarea preview-content"
              contenteditable="true"
              role="textbox"
              aria-multiline="true"
              data-placeholder="在这里直接编辑文章内容"
              @input="handleVisualEditorInput"
              @keydown="handleEditorShortcut"
              @paste="handleEditorPaste"
              @focus="saveEditorSelection"
              @keyup="saveEditorSelection"
              @mouseup="saveEditorSelection"
            ></div>
            <aside v-if="editorPanel === 'catalog'" class="editor-side-panel editor-catalog">
              <h4>文章目录</h4>
              <button v-for="heading in editorHeadings" :key="`${heading.index}-${heading.text}`" type="button" :style="{ paddingLeft: `${12 + (heading.level - 1) * 12}px` }" @click="goToHeading(heading)">{{ heading.text }}</button>
              <p v-if="editorHeadings.length === 0">正文中暂无标题</p>
            </aside>
          </div>
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
import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { createAdminArticle, getAdminArticleDetail, updateArticle, uploadFile } from '../../api/article'
import { getAdminCategories } from '../../api/category'
import { getAdminTags } from '../../api/tag'
import { message } from '../../utils/message'
import { renderMarkdown } from '../../utils/markdown'

const route = useRoute()
const router = useRouter()
// 根据当前响应式状态计算派生数据。
const isEdit = computed(() => !!route.params.id)
const editorRef = ref(null)
const contentImageInputRef = ref(null)
const activeToolbarMenu = ref('')
const editorPanel = ref('')
const isEditorFullscreen = ref(false)
const editorHistory = ref([''])
const editorHistoryIndex = ref(0)
const tableSettings = ref({ columns: 3, rows: 3, alignment: 'left' })
let contentImageSelection = null
let savedEditorRange = null

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
const colorOptions = [
  { label: '默认', value: '#e8eaed' },
  { label: '红色', value: '#ef4444' },
  { label: '橙色', value: '#f97316' },
  { label: '黄色', value: '#eab308' },
  { label: '绿色', value: '#22c55e' },
  { label: '青色', value: '#2dd4bf' },
  { label: '蓝色', value: '#3b82f6' },
  { label: '紫色', value: '#a855f7' }
]
const backgroundOptions = [
  { label: '深灰', value: '#374151' },
  { label: '红色', value: '#7f1d1d' },
  { label: '橙色', value: '#7c2d12' },
  { label: '黄色', value: '#713f12' },
  { label: '绿色', value: '#14532d' },
  { label: '青色', value: '#134e4a' },
  { label: '蓝色', value: '#1e3a8a' },
  { label: '紫色', value: '#581c87' }
]
const articleTemplates = [
  {
    name: '通用文章结构',
    content: '## 前言\n\n在这里介绍文章背景。\n\n## 正文\n\n在这里编写主要内容。\n\n## 总结\n\n在这里总结全文。'
  },
  {
    name: '技术教程',
    content: '## 目标\n\n说明本教程要完成的目标。\n\n## 环境准备\n\n- 环境要求\n- 依赖工具\n\n## 实现步骤\n\n### 第一步\n\n### 第二步\n\n## 常见问题\n\n## 总结'
  },
  {
    name: '更新日志',
    content: '## 新增\n\n- 新功能\n\n## 优化\n\n- 优化内容\n\n## 修复\n\n- 修复问题'
  }
]

const canUndo = computed(() => editorHistoryIndex.value > 0)
const canRedo = computed(() => editorHistoryIndex.value < editorHistory.value.length - 1)
const editorHeadings = computed(() => {
  const documentFragment = new DOMParser().parseFromString(form.value.content || '', 'text/html')
  return Array.from(documentFragment.querySelectorAll('h1, h2, h3, h4, h5, h6')).map((heading, index) => ({
    index,
    level: Number(heading.tagName.slice(1)),
    text: heading.textContent?.trim() || `标题 ${index + 1}`
  }))
})

// 将 Markdown 模板转换为 HTML 后插入可视化编辑器。
function insertSnippet(snippet) {
  insertHtmlAtSelection(renderMarkdown(snippet))
}

function commitEditorContent(content, syncEditor = true) {
  form.value.content = content
  pushEditorHistory(content)
  if (syncEditor) {
    nextTick(() => {
      if (editorRef.value && editorRef.value.innerHTML !== content) {
        editorRef.value.innerHTML = content
      }
    })
  }
}

function resetEditorHistory(content) {
  editorHistory.value = [content]
  editorHistoryIndex.value = 0
}

function pushEditorHistory(content) {
  if (editorHistory.value[editorHistoryIndex.value] === content) return
  const history = editorHistory.value.slice(0, editorHistoryIndex.value + 1)
  history.push(content)
  if (history.length > 100) history.shift()
  editorHistory.value = history
  editorHistoryIndex.value = history.length - 1
}

function handleVisualEditorInput(event) {
  const content = event.currentTarget.innerHTML
  form.value.content = content
  pushEditorHistory(content)
  saveEditorSelection()
}

function undoEditor() {
  if (!canUndo.value) return
  editorHistoryIndex.value--
  commitEditorContent(editorHistory.value[editorHistoryIndex.value])
}

function redoEditor() {
  if (!canRedo.value) return
  editorHistoryIndex.value++
  commitEditorContent(editorHistory.value[editorHistoryIndex.value])
}

function toggleToolbarMenu(menu) {
  activeToolbarMenu.value = activeToolbarMenu.value === menu ? '' : menu
}

function closeToolbarMenu() {
  activeToolbarMenu.value = ''
}

function applyInlineFormat(before) {
  closeToolbarMenu()
  const command = before === '**' ? 'bold' : before === '*' ? 'italic' : 'strikeThrough'
  executeEditorCommand(command)
}

function applyHeading(level) {
  closeToolbarMenu()
  executeEditorCommand('formatBlock', `h${level}`)
}

function applyList(type) {
  closeToolbarMenu()
  if (type === 'ordered') {
    executeEditorCommand('insertOrderedList')
  } else if (type === 'task') {
    insertHtmlAtSelection('<ul class="task-list"><li><input type="checkbox" /> 任务项</li></ul>')
  } else {
    executeEditorCommand('insertUnorderedList')
  }
}

function applyQuote() {
  closeToolbarMenu()
  executeEditorCommand('formatBlock', 'blockquote')
}

function applyTextColor(color) {
  closeToolbarMenu()
  executeEditorCommand('foreColor', color)
}

function applyTextBackground(color) {
  closeToolbarMenu()
  executeEditorCommand('hiliteColor', color)
}

function applyAlignment(alignment) {
  closeToolbarMenu()
  const command = alignment === 'center' ? 'justifyCenter' : alignment === 'right' ? 'justifyRight' : 'justifyLeft'
  executeEditorCommand(command)
}

function applyInlineCode() {
  closeToolbarMenu()
  wrapSelectionInElement('code', '代码')
}

function applyCodeBlock() {
  closeToolbarMenu()
  restoreEditorSelection()
  const selection = window.getSelection()
  const code = selection?.toString() || '在这里输入代码'
  const pre = document.createElement('pre')
  const codeElement = document.createElement('code')
  codeElement.textContent = code
  pre.appendChild(codeElement)
  insertNodeAtSelection(pre)
}

function insertCustomTable() {
  closeToolbarMenu()
  const columns = Math.min(10, Math.max(1, Number(tableSettings.value.columns) || 1))
  const rows = Math.min(20, Math.max(1, Number(tableSettings.value.rows) || 1))
  const alignment = tableSettings.value.alignment
  const header = `<thead><tr>${Array.from({ length: columns }, (_, index) => `<th style="text-align:${alignment}">表头 ${index + 1}</th>`).join('')}</tr></thead>`
  const body = `<tbody>${Array.from({ length: rows }, (_, rowIndex) => `<tr>${Array.from({ length: columns }, (_, columnIndex) => `<td style="text-align:${alignment}">内容 ${rowIndex + 1}-${columnIndex + 1}</td>`).join('')}</tr>`).join('')}</tbody>`
  insertHtmlAtSelection(`<table>${header}${body}</table><p><br></p>`)
}

function insertImageUrl() {
  closeToolbarMenu()
  const imageUrl = window.prompt('请输入图片 URL')?.trim()
  if (!imageUrl) return
  const description = window.prompt('请输入图片说明', '图片')?.trim() || '图片'
  insertHtmlAtSelection(`<img src="${escapeHtmlAttribute(imageUrl)}" alt="${escapeHtmlAttribute(description)}" />`)
}

function insertVideo() {
  closeToolbarMenu()
  const videoUrl = window.prompt('请输入视频 URL')?.trim()
  if (!videoUrl) return
  insertHtmlAtSelection(`<video controls preload="metadata" src="${escapeHtmlAttribute(videoUrl)}"></video><p><br></p>`)
}

function insertFormula(displayMode) {
  closeToolbarMenu()
  const formula = window.prompt('请输入 LaTeX 公式', 'E = mc^2')?.trim()
  if (!formula) return
  if (displayMode) {
    insertHtmlAtSelection(renderMarkdown(`$$\n${formula}\n$$`))
  } else {
    const rendered = new DOMParser().parseFromString(renderMarkdown(`$${formula}$`), 'text/html').body.firstElementChild?.innerHTML
    insertHtmlAtSelection(rendered || formula)
  }
}

function insertLink() {
  closeToolbarMenu()
  restoreEditorSelection()
  const selected = window.getSelection()?.toString() || ''
  const url = window.prompt('请输入链接地址', 'https://')?.trim()
  if (!url) return
  const label = selected || window.prompt('请输入链接文字', '链接文字')?.trim() || '链接文字'
  if (selected) {
    executeEditorCommand('createLink', url)
  } else {
    insertHtmlAtSelection(`<a href="${escapeHtmlAttribute(url)}">${escapeHtml(label)}</a>`)
  }
}

function insertArticleTemplate(template) {
  closeToolbarMenu()
  insertSnippet(template.content)
}

function saveEditorSelection() {
  const editor = editorRef.value
  const selection = window.getSelection()
  if (!editor || !selection?.rangeCount) return
  const range = selection.getRangeAt(0)
  if (editor.contains(range.commonAncestorContainer)) {
    savedEditorRange = range.cloneRange()
  }
}

function restoreEditorSelection(range = savedEditorRange) {
  const editor = editorRef.value
  if (!editor) return false
  editor.focus()
  const selection = window.getSelection()
  selection.removeAllRanges()
  if (range && editor.contains(range.commonAncestorContainer)) {
    selection.addRange(range)
  } else {
    const fallbackRange = document.createRange()
    fallbackRange.selectNodeContents(editor)
    fallbackRange.collapse(false)
    selection.addRange(fallbackRange)
  }
  return true
}

function executeEditorCommand(command, value = null) {
  if (!restoreEditorSelection()) return
  document.execCommand(command, false, value)
  syncVisualEditorContent()
}

function syncVisualEditorContent() {
  if (!editorRef.value) return
  form.value.content = editorRef.value.innerHTML
  pushEditorHistory(form.value.content)
  saveEditorSelection()
}

function insertHtmlAtSelection(html, range = savedEditorRange) {
  if (!restoreEditorSelection(range)) return
  const selection = window.getSelection()
  if (!selection?.rangeCount) return
  const activeRange = selection.getRangeAt(0)
  activeRange.deleteContents()
  const template = document.createElement('template')
  template.innerHTML = html
  const fragment = template.content
  const lastNode = fragment.lastChild
  activeRange.insertNode(fragment)
  if (lastNode) {
    activeRange.setStartAfter(lastNode)
    activeRange.collapse(true)
    selection.removeAllRanges()
    selection.addRange(activeRange)
  }
  syncVisualEditorContent()
}

function insertNodeAtSelection(node) {
  if (!restoreEditorSelection()) return
  const selection = window.getSelection()
  const range = selection.getRangeAt(0)
  range.deleteContents()
  range.insertNode(node)
  range.setStartAfter(node)
  range.collapse(true)
  selection.removeAllRanges()
  selection.addRange(range)
  syncVisualEditorContent()
}

function wrapSelectionInElement(tagName, fallbackText) {
  if (!restoreEditorSelection()) return
  const selection = window.getSelection()
  const range = selection.getRangeAt(0)
  const element = document.createElement(tagName)
  element.textContent = selection.toString() || fallbackText
  range.deleteContents()
  range.insertNode(element)
  range.selectNodeContents(element)
  selection.removeAllRanges()
  selection.addRange(range)
  syncVisualEditorContent()
}

function escapeHtml(value) {
  const element = document.createElement('div')
  element.textContent = value
  return element.innerHTML
}

function escapeHtmlAttribute(value) {
  return escapeHtml(value).replaceAll('`', '&#96;')
}

function toggleEditorPanel(panel) {
  closeToolbarMenu()
  editorPanel.value = editorPanel.value === panel ? '' : panel
}

function toggleEditorFullscreen() {
  closeToolbarMenu()
  isEditorFullscreen.value = !isEditorFullscreen.value
  document.body.style.overflow = isEditorFullscreen.value ? 'hidden' : ''
}

function goToHeading(heading) {
  const headingElement = editorRef.value?.querySelectorAll('h1, h2, h3, h4, h5, h6')[heading.index]
  if (!headingElement) return
  headingElement.scrollIntoView({ behavior: 'smooth', block: 'center' })
  const range = document.createRange()
  range.selectNodeContents(headingElement)
  range.collapse(false)
  savedEditorRange = range
  restoreEditorSelection(range)
}

function handleEditorShortcut(event) {
  if (event.key === 'Escape') {
    closeToolbarMenu()
    if (isEditorFullscreen.value) toggleEditorFullscreen()
    return
  }
  if (!(event.ctrlKey || event.metaKey)) return
  const key = event.key.toLowerCase()
  if (key === 'z') {
    event.preventDefault()
    event.shiftKey ? redoEditor() : undoEditor()
  } else if (key === 'y') {
    event.preventDefault()
    redoEditor()
  } else if (key === 'b') {
    event.preventDefault()
    applyInlineFormat('**', '**')
  } else if (key === 'i') {
    event.preventDefault()
    applyInlineFormat('*', '*')
  } else if (key === 'k') {
    event.preventDefault()
    insertLink()
  }
}

function openContentImagePicker() {
  closeToolbarMenu()
  saveEditorSelection()
  contentImageSelection = savedEditorRange?.cloneRange() || null
  contentImageInputRef.value?.click()
}

// 上传正文图片，并在打开文件选择器前记录的光标位置显示图片。
async function handleContentImageUpload(event) {
  const file = event.target.files?.[0]
  if (!file) return

  try {
    await uploadAndInsertContentImage(file, contentImageSelection)
  } catch (error) {
    message.error('正文图片上传失败')
  } finally {
    contentImageSelection = null
    event.target.value = ''
  }
}

async function uploadAndInsertContentImage(file, selection) {
  const res = await uploadFile(file)
  const imageUrl = res.data?.url || res.data || ''
  if (!imageUrl) throw new Error('empty image url')

  insertContentImage(imageUrl, file.name, selection)
  message.success('图片已上传并插入正文')
}

function insertContentImage(imageUrl, fileName = '图片', selection) {
  const alt = fileName.replace(/\.[^.]+$/, '').replace(/[[\]]/g, '') || '图片'
  insertHtmlAtSelection(`<img src="${escapeHtmlAttribute(imageUrl)}" alt="${escapeHtmlAttribute(alt)}" /><p><br></p>`, selection)
}

// 支持从系统剪贴板或网页直接粘贴图片。
async function handleEditorPaste(event) {
  const clipboard = event.clipboardData
  if (!clipboard) return

  saveEditorSelection()
  const selection = savedEditorRange?.cloneRange() || null
  const imageItem = Array.from(clipboard.items || []).find(item => item.kind === 'file' && item.type.startsWith('image/'))
  const imageFile = imageItem?.getAsFile() || Array.from(clipboard.files || []).find(file => file.type.startsWith('image/'))

  if (imageFile) {
    event.preventDefault()
    try {
      await uploadAndInsertContentImage(imageFile, selection)
    } catch (error) {
      message.error('粘贴图片上传失败')
    }
    return
  }

  const html = clipboard.getData('text/html')
  const imageSource = html ? new DOMParser().parseFromString(html, 'text/html').querySelector('img')?.src : ''
  if (!imageSource) return

  event.preventDefault()
  if (imageSource.startsWith('data:image/') || imageSource.startsWith('blob:')) {
    try {
      const blob = await fetch(imageSource).then(response => response.blob())
      const extension = blob.type.split('/')[1] || 'png'
      const file = new File([blob], `clipboard-${Date.now()}.${extension}`, { type: blob.type })
      await uploadAndInsertContentImage(file, selection)
    } catch (error) {
      message.error('粘贴图片上传失败')
    }
  } else {
    insertContentImage(imageSource, '图片', selection)
    message.success('图片已插入正文')
  }
}

// 处理用户操作或浏览器事件。
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

// 处理用户操作或浏览器事件。
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
    commitEditorContent(renderMarkdown(content))

    if (!form.value.title.trim()) {
      const heading = content.match(/^#\s+(.+)$/m)?.[1]?.trim()
      form.value.title = heading || file.name.replace(/\.(md|markdown|txt)$/i, '')
    }

    message.success('文章已导入')
  } catch (error) {
    message.error('文章导入失败')
  } finally {
    event.target.value = ''
  }
}

// 处理用户操作或浏览器事件。
async function handleSave() {
  syncVisualEditorContent()
  const hasVisualContent = Boolean(
    editorRef.value?.textContent?.trim() ||
    editorRef.value?.querySelector('img, video, table, hr, .katex')
  )
  if (!form.value.title.trim() || !form.value.categoryId || !hasVisualContent) {
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
  const [categoryRes, tagRes] = await Promise.all([getAdminCategories(), getAdminTags()])
  categories.value = categoryRes.data || []
  tags.value = tagRes.data || []

  if (!isEdit.value) return

  const res = await getAdminArticleDetail(route.params.id)
  const data = res.data
  const visualContent = renderMarkdown(data.content || '')
  form.value = {
    title: data.title,
    categoryId: data.categoryId,
    summary: data.summary,
    coverImage: data.coverImage,
    content: visualContent,
    status: data.status,
    tagIds: (data.tags || []).map(tag => tag.id)
  }
  resetEditorHistory(visualContent)
  nextTick(() => {
    if (editorRef.value) editorRef.value.innerHTML = visualContent
  })
})

onUnmounted(() => {
  document.body.style.overflow = ''
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

.editor-toolbar {
  position: relative;
  z-index: 5;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px;
  padding: 9px 12px;
  background: #f7f8fa;
  border: 1px solid var(--border);
  border-bottom: none;
  border-radius: var(--radius-sm) var(--radius-sm) 0 0;
  color: #505363;
  box-shadow: 0 1px 0 rgba(0, 0, 0, .08);
}

.toolbar-group {
  display: flex;
  align-items: stretch;
  gap: 2px;
}

.toolbar-divider {
  width: 1px;
  height: 42px;
  margin: 4px 2px;
  background: rgba(80, 83, 99, .18);
}

.editor-toolbar .toolbar-action {
  min-width: 50px;
  height: 52px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  padding: 3px 7px;
  background: transparent;
  border: 1px solid transparent;
  border-radius: 6px;
  color: #505363;
  cursor: pointer;
  font-size: 12px;
  line-height: 1.15;
  white-space: nowrap;
  transition: all 0.2s;
}

.editor-toolbar .toolbar-action:hover,
.editor-toolbar .toolbar-action.active {
  background: #e7e9ee;
  border-color: #d8dbe3;
  color: #202331;
}

.editor-toolbar .toolbar-action:disabled {
  cursor: not-allowed;
  opacity: .35;
}

.toolbar-icon {
  width: 28px;
  height: 24px;
  min-height: 24px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 23px;
  line-height: 1;
}

.toolbar-icon-bold {
  font-family: Georgia, serif;
  font-weight: 800;
}

.toolbar-icon-italic {
  font-family: Georgia, serif;
  font-style: italic;
}

.toolbar-icon-format,
.toolbar-icon-color,
.toolbar-icon-background {
  position: relative;
  padding: 0;
  border: 0;
  border-radius: 0;
  background: transparent;
  color: #505363;
  font-family: Georgia, serif;
}

.toolbar-icon-color::after,
.toolbar-icon-background::after {
  content: '';
  position: absolute;
  right: 4px;
  bottom: 0;
  left: 4px;
  height: 3px;
  border-radius: 999px;
}

.toolbar-icon-color::after {
  background: #ef4444;
}

.toolbar-icon-background::after {
  background: #3b82f6;
}

.toolbar-caret {
  position: absolute;
  top: 8px;
  right: 4px;
  font-size: 9px;
}

.toolbar-menu-wrap {
  position: relative;
}

.toolbar-menu-wrap > .toolbar-action {
  position: relative;
  padding-right: 12px;
}

.toolbar-menu {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  z-index: 30;
  min-width: 128px;
  padding: 6px;
  border: 1px solid var(--border-light);
  border-radius: 9px;
  background: #1b2433;
  box-shadow: 0 12px 28px rgba(0, 0, 0, .38);
}

.toolbar-menu-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.editor-toolbar .toolbar-menu-list button {
  width: 100%;
  min-height: 34px;
  padding: 7px 10px;
  border: 0;
  border-radius: 6px;
  background: transparent;
  color: var(--text-secondary);
  cursor: pointer;
  text-align: left;
  white-space: nowrap;
}

.editor-toolbar .toolbar-menu-list button:hover {
  background: var(--accent-dim);
  color: var(--accent);
}

.toolbar-menu-wide {
  min-width: 170px;
}

.table-settings-menu {
  width: 220px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 12px;
}

.table-settings-menu label {
  display: grid;
  grid-template-columns: 72px minmax(0, 1fr);
  align-items: center;
  gap: 8px;
  color: var(--text-secondary);
  font-size: 13px;
}

.table-settings-menu input,
.table-settings-menu select {
  width: 100%;
  height: 32px;
  padding: 0 8px;
  border: 1px solid var(--border-light);
  border-radius: 6px;
  outline: none;
  background: var(--bg-primary);
  color: var(--text-primary);
}

.table-settings-menu input:focus,
.table-settings-menu select:focus {
  border-color: var(--accent);
}

.editor-toolbar .table-settings-menu .table-insert-button {
  width: 100%;
  height: 34px;
  border: 0;
  border-radius: 6px;
  background: var(--accent);
  color: var(--bg-primary);
  cursor: pointer;
  font-weight: 600;
}

.toolbar-group:last-child .toolbar-menu {
  right: 0;
  left: auto;
}

.color-menu {
  width: 152px;
  display: grid;
  grid-template-columns: repeat(4, 28px);
  gap: 6px;
}

.editor-toolbar .color-menu button {
  width: 28px;
  min-width: 28px;
  height: 28px;
  border: 2px solid rgba(255, 255, 255, .18);
  border-radius: 6px;
  cursor: pointer;
}

.editor-toolbar .color-menu button:hover {
  border-color: #fff;
}

.content-image-input {
  display: none;
}

.editor-workspace {
  display: grid;
  grid-template-columns: minmax(0, 1fr);
  min-height: 430px;
}

.editor-workspace.has-panel {
  grid-template-columns: minmax(0, 1fr) minmax(280px, 38%);
}

.editor-textarea {
  border-radius: 0 0 var(--radius-sm) var(--radius-sm);
  overflow: auto;
  padding: 20px;
  font-family: inherit;
  font-size: 15px;
  line-height: 1.8;
  min-height: 400px;
  height: 100%;
  outline: none;
  resize: vertical;
  cursor: text;
}

.editor-textarea:empty::before {
  content: attr(data-placeholder);
  color: var(--text-muted);
  pointer-events: none;
}

.editor-workspace.has-panel .editor-textarea {
  border-radius: 0 0 0 var(--radius-sm);
}

.editor-side-panel {
  min-width: 0;
  max-height: 620px;
  overflow: auto;
  padding: 18px;
  border: 1px solid var(--border);
  border-left: 0;
  border-radius: 0 0 var(--radius-sm) 0;
  background: var(--bg-secondary);
}

.editor-side-panel > h4 {
  margin-bottom: 14px;
  color: var(--text-primary);
  font-size: 15px;
}

.editor-catalog {
  display: flex;
  flex-direction: column;
  align-items: stretch;
}

.editor-catalog button {
  padding-block: 7px;
  border: 0;
  background: transparent;
  color: var(--text-secondary);
  cursor: pointer;
  text-align: left;
}

.editor-catalog button:hover {
  color: var(--accent);
  background: var(--accent-dim);
}

.editor-catalog p {
  color: var(--text-muted);
  font-size: 13px;
}

.editor-wrap.fullscreen {
  position: fixed;
  inset: 10px;
  z-index: 1200;
  display: flex;
  flex-direction: column;
  border-radius: 12px;
  background: var(--bg-primary);
  box-shadow: 0 24px 80px rgba(0, 0, 0, .65);
}

.editor-wrap.fullscreen .editor-toolbar {
  flex: 0 0 auto;
  border-radius: 12px 12px 0 0;
}

.editor-wrap.fullscreen .editor-workspace {
  flex: 1;
  min-height: 0;
}

.editor-wrap.fullscreen .editor-textarea,
.editor-wrap.fullscreen .editor-side-panel {
  max-height: none;
  resize: none;
}

.preview-content {
  line-height: 1.8;
  color: var(--text-primary);
}

.preview-content :deep(h1),
.preview-content :deep(h2),
.preview-content :deep(h3),
.preview-content :deep(h4),
.preview-content :deep(h5),
.preview-content :deep(h6) {
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

.preview-content :deep(pre code) {
  padding: 0;
  background: transparent;
}

.preview-content :deep(blockquote) {
  margin: 16px 0;
  padding: 10px 16px;
  border-left: 4px solid var(--accent);
  background: var(--accent-dim);
  color: var(--text-secondary);
}

.preview-content :deep(ul),
.preview-content :deep(ol) {
  margin: 12px 0 16px;
  padding-left: 28px;
}

.preview-content :deep(li) {
  margin: 6px 0;
}

.preview-content :deep(input[type='checkbox']) {
  margin-right: 7px;
  accent-color: var(--accent);
}

.preview-content :deep(img) {
  max-width: 100%;
  border-radius: var(--radius-sm);
}

.preview-content :deep(video) {
  display: block;
  width: min(100%, 860px);
  margin: 16px auto;
  border-radius: var(--radius-sm);
  background: #000;
}

.preview-content :deep(.katex-display) {
  max-width: 100%;
  overflow-x: auto;
  overflow-y: hidden;
  padding: 8px 0;
}

.preview-content :deep(table) {
  width: 100%;
  margin: 16px 0;
  border-collapse: collapse;
}

.preview-content :deep(th),
.preview-content :deep(td) {
  padding: 10px 12px;
  border: 1px solid var(--border-light);
  text-align: left;
}

.preview-content :deep(th) {
  background: var(--accent-dim);
}

.preview-content :deep(hr) {
  margin: 24px 0;
  border: 0;
  border-top: 1px solid var(--border-light);
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

  .editor-workspace.has-panel {
    grid-template-columns: minmax(0, 1fr);
  }

  .editor-side-panel {
    max-height: 360px;
    border-left: 1px solid var(--border);
    border-radius: 0 0 var(--radius-sm) var(--radius-sm);
  }

  .toolbar-divider {
    display: none;
  }
}
</style>
