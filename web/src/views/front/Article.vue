<template>
  <div v-if="article" class="article-page">
    <div class="article-layout">
      <article class="article-main">
        <div class="article-header card" :data-pet-hint="`正在阅读：${article.title}`">
          <img
            v-if="article.coverImage"
            :src="article.coverImage"
            class="article-cover"
            :alt="article.title"
          />
          <h1 class="article-title">{{ article.title }}</h1>
          <div class="article-meta">
            <span class="meta-item">作者 {{ article.author?.nickname || '-' }}</span>
            <span class="meta-item">{{ formatDate(article.createdAt) }}</span>
            <span class="meta-item">浏览 {{ article.viewCount || 0 }}</span>
            <span class="meta-category">{{ article.category?.name || '未分类' }}</span>
            <span v-for="tag in article.tags || []" :key="tag.id" class="meta-tag">{{ tag.name }}</span>
          </div>
        </div>

        <div class="article-content card" v-html="renderedContent"></div>

        <div class="comment-section card">
          <div class="section-head">
            <h3 class="section-title">评论（{{ comments.length }}）</h3>
            <span class="section-tip">支持二级评论</span>
          </div>

          <div class="comment-form root-comment-form">
            <div class="comment-form-fields" :class="{ 'is-authenticated': userStore.isLoggedIn }">
              <div class="comment-form-avatar" aria-hidden="true">
                {{ (userStore.isLoggedIn ? userStore.user?.nickname : guestName)?.[0] || '?' }}
              </div>
              <template v-if="userStore.isLoggedIn">
                <div class="comment-form-identity">
                  <span class="comment-form-identity-name">{{ userStore.user?.nickname || userStore.user?.username }}</span>
                  <span class="comment-form-identity-status">已登录，将以账户身份发表评论</span>
                </div>
              </template>
              <template v-else>
                <div class="comment-form-field">
                  <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M20 21a8 8 0 0 0-16 0M12 13a5 5 0 1 0 0-10 5 5 0 0 0 0 10Z" /></svg>
                  <input v-model="guestName" type="text" maxlength="50" placeholder="昵称" />
                </div>
                <div class="comment-form-field">
                  <svg viewBox="0 0 24 24" aria-hidden="true"><rect x="3" y="5" width="18" height="14" rx="2" /><path d="m3 7 9 6 9-6" /></svg>
                  <input v-model="guestEmail" type="email" maxlength="255" placeholder="邮箱" />
                </div>
                <div class="comment-form-field comment-website-field">
                  <svg viewBox="0 0 24 24" aria-hidden="true"><circle cx="12" cy="12" r="9" /><path d="M3 12h18M12 3c2.4 2.5 3.6 5.5 3.6 9S14.4 18.5 12 21c-2.4-2.5-3.6-5.5-3.6-9S9.6 5.5 12 3Z" /></svg>
                  <input v-model="guestWebsite" type="url" maxlength="255" placeholder="网站（可选）" />
                </div>
              </template>
            </div>
            <div v-if="replyTarget && !userStore.isLoggedIn" class="root-reply-banner">
              正在回复 {{ replyTarget.author?.nickname || '匿名用户' }}
              <button type="button" class="reply-cancel" @click="clearReply">取消</button>
            </div>
            <textarea
              v-model="rootCommentContent"
              class="input comment-input"
              placeholder="写下你的评论"
              rows="4"
            ></textarea>
            <button
              class="btn btn-primary"
              :disabled="submitting || !rootCommentContent.trim()"
              @click="submitComment()"
            >
              {{ submitting && !replyTarget ? '提交中...' : '发表评论' }}
            </button>
          </div>
          <div v-if="commentsLoading" class="loading-state">评论加载中...</div>
          <div v-else class="comment-list">
            <div v-for="comment in commentTree" :key="comment.id" class="comment-thread">
              <div class="comment-item">
                <div class="comment-avatar">{{ comment.author?.nickname?.[0] || '?' }}</div>
                <div class="comment-body">
                  <div class="comment-header">
                    <span class="comment-author">{{ comment.author?.nickname || '匿名用户' }}</span>
                    <span v-if="comment.replyToAuthor" class="reply-to">
                      回复 {{ comment.replyToAuthor.nickname || '匿名用户' }}
                    </span>
                    <span class="comment-time">{{ formatDate(comment.createdAt) }}</span>
                    <button
                      v-if="canReply(comment)"
                      type="button"
                      class="reply-link"
                      @click="startReply(comment)"
                    >
                      回复
                    </button>
                    <button
                      v-if="canDelete(comment)"
                      type="button"
                      class="delete-link"
                      @click="removeComment(comment)"
                    >
                      删除
                    </button>
                  </div>
                  <div class="comment-content">{{ comment.content }}</div>

                  <div
                    v-if="userStore.isLoggedIn && replyTarget?.id === comment.id"
                    class="comment-form inline-reply-form"
                  >
                    <div class="reply-banner">
                      正在回复 {{ comment.author?.nickname || '匿名用户' }}
                      <button type="button" class="reply-cancel" @click="clearReply">取消</button>
                    </div>
                    <textarea
                      v-model="replyContent"
                      class="input comment-input"
                      :placeholder="`回复 ${comment.author?.nickname || '匿名用户'}`"
                      rows="3"
                    ></textarea>
                    <button
                      class="btn btn-primary"
                      :disabled="submitting || !replyContent.trim()"
                      @click="submitComment(comment)"
                    >
                      {{ submitting ? '提交中...' : '发表回复' }}
                    </button>
                  </div>
                </div>
              </div>

              <div v-if="comment.children.length" class="reply-list">
                <div v-for="reply in comment.children" :key="reply.id" class="comment-item comment-item-reply">
                  <div class="comment-avatar">{{ reply.author?.nickname?.[0] || '?' }}</div>
                  <div class="comment-body">
                    <div class="comment-header">
                      <span class="comment-author">{{ reply.author?.nickname || '匿名用户' }}</span>
                      <span v-if="reply.replyToAuthor" class="reply-to">
                        回复 {{ reply.replyToAuthor.nickname || '匿名用户' }}
                      </span>
                      <span class="comment-time">{{ formatDate(reply.createdAt) }}</span>
                      <button
                        v-if="canReply(reply)"
                        type="button"
                        class="reply-link"
                        @click="startReply(reply)"
                      >
                        回复
                      </button>
                      <button
                        v-if="canDelete(reply)"
                        type="button"
                        class="delete-link"
                        @click="removeComment(reply)"
                      >
                        删除
                      </button>
                    </div>
                    <div class="comment-content">{{ reply.content }}</div>

                    <div
                      v-if="userStore.isLoggedIn && replyTarget?.id === reply.id"
                      class="comment-form inline-reply-form"
                    >
                      <div class="reply-banner">
                        正在回复 {{ reply.author?.nickname || '匿名用户' }}
                        <button type="button" class="reply-cancel" @click="clearReply">取消</button>
                      </div>
                      <textarea
                        v-model="replyContent"
                        class="input comment-input"
                        :placeholder="`回复 ${reply.author?.nickname || '匿名用户'}`"
                        rows="3"
                      ></textarea>
                      <button
                        class="btn btn-primary"
                        :disabled="submitting || !replyContent.trim()"
                        @click="submitComment(reply)"
                      >
                        {{ submitting ? '提交中...' : '发表回复' }}
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-if="!commentsLoading && comments.length === 0" class="empty-comments">
            暂无评论，来发表第一条评论吧。
          </div>
        </div>
      </article>
    </div>
  </div>

  <div v-else-if="pageLoading" class="loading-state">加载中...</div>
  <div v-else class="loading-state">文章不存在或暂时无法访问。</div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { marked } from 'marked'
import { getArticleFull } from '../../api/article'
import { createComment, deleteMyComment, getComments } from '../../api/comment'
import { useUserStore } from '../../stores/user'
import { message } from '../../utils/message'

const route = useRoute()
const userStore = useUserStore()

const article = ref(null)
const comments = ref([])
const rootCommentContent = ref('')
const guestName = ref('')
const guestEmail = ref('')
const guestWebsite = ref('')
const replyContent = ref('')
const submitting = ref(false)
const pageLoading = ref(false)
const commentsLoading = ref(false)
const replyTarget = ref(null)

const renderedContent = computed(() => {
  try {
    return marked(article.value?.content || '')
  } catch (error) {
    return article.value?.content || ''
  }
})

const commentTree = computed(() => {
  const roots = []
  const byId = new Map()

  for (const comment of comments.value) {
    byId.set(comment.id, { ...comment, children: [] })
  }

  function findRootParentId(comment) {
    if (comment.parentId) {
      return comment.parentId
    }
    if (!comment.replyToId) {
      return null
    }

    const replyTo = byId.get(comment.replyToId)
    if (!replyTo) {
      return null
    }
    return replyTo.parentId || replyTo.id
  }

  function resolveReplyToAuthor(comment) {
    if (comment.replyToAuthor || !comment.replyToId) {
      return comment.replyToAuthor || null
    }

    const replyTo = byId.get(comment.replyToId)
    return replyTo?.author || null
  }

  for (const comment of byId.values()) {
    comment.replyToAuthor = resolveReplyToAuthor(comment)

    const parentId = findRootParentId(comment)
    if (parentId) {
      const parent = byId.get(parentId)
      if (parent) {
        parent.children.push(comment)
        continue
      }
    }

    roots.push(comment)
  }

  for (const root of roots) {
    root.children.sort((a, b) => new Date(a.createdAt) - new Date(b.createdAt))
  }

  return roots
})

function formatDate(dateStr) {
  if (!dateStr) {
    return '-'
  }

  const d = new Date(dateStr)
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

function canReply(comment) {
  return !userStore.isLoggedIn || comment.userId !== userStore.user?.id
}

function canDelete(comment) {
  return userStore.isLoggedIn && comment.userId === userStore.user?.id
}

async function removeComment(comment) {
  if (!canDelete(comment)) {
    return
  }

  try {
    await deleteMyComment(comment.id)
    if (replyTarget.value?.id === comment.id) {
      clearReply()
    }
    message.success('评论已删除')
    await loadComments()
  } catch (error) {
    // The request layer already shows the backend message.
  }
}

function clearReply() {
  replyTarget.value = null
  replyContent.value = ''
}

function startReply(comment) {
  if (!canReply(comment)) {
    return
  }
  replyTarget.value = comment
  replyContent.value = ''
}

async function loadArticle() {
  const res = await getArticleFull(route.params.id)
  article.value = res.data
}

async function loadComments() {
  commentsLoading.value = true
  try {
    const res = await getComments(route.params.id)
    comments.value = res.data || []
  } finally {
    commentsLoading.value = false
  }
}

async function initializePage() {
  pageLoading.value = true
  article.value = null
  comments.value = []
  rootCommentContent.value = ''
  guestName.value = ''
  guestEmail.value = ''
  guestWebsite.value = ''
  replyContent.value = ''
  replyTarget.value = null

  try {
    await loadArticle()
    await loadComments()
  } catch (error) {
    article.value = null
    comments.value = []
  } finally {
    pageLoading.value = false
  }
}

async function submitComment(targetComment = null) {
  const replyComment = targetComment || replyTarget.value
  const isInlineReply = !!targetComment
  const content = isInlineReply ? replyContent.value.trim() : rootCommentContent.value.trim()

  if (!userStore.isLoggedIn) {
    if (!guestName.value.trim()) {
      message.warning('请填写昵称')
      return
    }
    if (!guestEmail.value.trim()) {
      message.warning('请填写邮箱')
      return
    }
    if (!/^\S+@\S+\.\S+$/.test(guestEmail.value.trim())) {
      message.warning('请输入正确的邮箱地址，例如 name@example.com')
      return
    }
  }
  if (!content) {
    message.warning('请输入评论内容')
    return
  }

  submitting.value = true
  try {
    await createComment({
      articleId: Number(route.params.id),
      replyToId: replyComment?.id,
      content,
      guestName: userStore.isLoggedIn ? undefined : guestName.value.trim(),
      guestEmail: userStore.isLoggedIn ? undefined : guestEmail.value.trim(),
      guestWebsite: userStore.isLoggedIn ? undefined : guestWebsite.value.trim()
    })

    if (isInlineReply) {
      replyContent.value = ''
      clearReply()
    } else {
      rootCommentContent.value = ''
      clearReply()
    }

    message.success('评论发布成功')
    await loadComments()
  } catch (error) {
    // The request layer already shows the backend message.
  } finally {
    submitting.value = false
  }
}

watch(
  () => route.params.id,
  () => {
    initializePage()
  },
  { immediate: true }
)
</script>

<style scoped>
.article-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr);
}

.article-main {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.article-cover {
  width: 100%;
  height: 300px;
  object-fit: cover;
  border-radius: var(--radius-sm);
  margin-bottom: 20px;
}

.article-title {
  font-size: 28px;
  font-weight: 700;
  line-height: 1.4;
  margin-bottom: 16px;
}

.article-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  align-items: center;
}

.meta-item {
  color: var(--text-muted);
  font-size: 13px;
}

.meta-category {
  padding: 3px 12px;
  background: var(--accent-dim);
  color: var(--accent);
  border-radius: 20px;
  font-size: 12px;
}

.meta-tag {
  padding: 3px 10px;
  border: 1px solid var(--border-light);
  color: var(--text-secondary);
  border-radius: 20px;
  font-size: 12px;
}

.article-content {
  line-height: 1.8;
  font-size: 16px;
  color: var(--text-primary);
}

.article-content :deep(h1),
.article-content :deep(h2),
.article-content :deep(h3) {
  margin-top: 24px;
  margin-bottom: 12px;
  color: var(--text-primary);
}

.article-content :deep(p) {
  margin-bottom: 16px;
}

.article-content :deep(code) {
  background: rgba(0, 0, 0, 0.3);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 14px;
}

.article-content :deep(pre) {
  background: rgba(0, 0, 0, 0.4);
  padding: 16px;
  border-radius: var(--radius-sm);
  overflow-x: auto;
  margin-bottom: 16px;
}

.article-content :deep(img) {
  max-width: 100%;
  border-radius: var(--radius-sm);
}

.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border);
}

.section-title {
  font-size: 18px;
  font-weight: 600;
}

.section-tip {
  display: none;
  color: var(--text-muted);
  font-size: 13px;
}

.comment-form {
  margin-bottom: 24px;
}

.root-comment-form {
  display: grid;
  grid-template-columns: 1fr auto;
  margin-bottom: 28px;
  overflow: hidden;
  border: 1px solid var(--accent);
  border-radius: 18px;
  background: rgba(0, 0, 0, 0.16);
  box-shadow: inset 0 0 0 1px rgba(45, 212, 191, 0.05);
}

.comment-form-fields {
  grid-column: 1 / -1;
  display: grid;
  grid-template-columns: 96px repeat(3, minmax(0, 1fr));
  min-height: 84px;
  align-items: center;
  border-bottom: 1px dashed var(--border-light);
}

.comment-form-fields.is-authenticated {
  grid-template-columns: 96px minmax(0, 1fr);
}

.comment-form-avatar {
  width: 52px;
  height: 52px;
  margin-left: 20px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--accent-dim);
  color: var(--accent);
  font-size: 20px;
  font-weight: 600;
}

.comment-form-field {
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 14px;
  height: 38px;
  padding: 0 24px;
  color: var(--text-secondary);
  font-size: 15px;
  border-left: 1px dashed var(--border-light);
}

.comment-form-identity {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 4px;
  padding: 0 24px;
  border-left: 1px dashed var(--border-light);
}

.comment-form-identity-name {
  overflow: hidden;
  color: var(--text-primary);
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.comment-form-identity-status {
  color: var(--text-muted);
  font-size: 13px;
}

.comment-form-field svg,
.root-comment-form .btn svg {
  width: 21px;
  height: 21px;
  flex: 0 0 auto;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.8;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.comment-form-field input {
  width: 100%;
  min-width: 0;
  padding: 0;
  border: 0;
  outline: 0;
  background: transparent;
  color: var(--text-secondary);
  font: inherit;
}

.comment-form-field input::placeholder {
  color: var(--text-secondary);
  opacity: 1;
}

.comment-form-field input:focus {
  color: var(--text-primary);
}

.comment-form-field span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.root-comment-form .comment-input {
  grid-column: 1 / -1;
  min-height: 190px;
  padding: 26px 28px;
  border: 0;
  border-radius: 0;
  background: transparent;
  box-shadow: none;
}

.root-reply-banner {
  grid-column: 1 / -1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 10px 20px;
  border-bottom: 1px dashed var(--border-light);
  background: rgba(45, 212, 191, 0.08);
  color: var(--text-secondary);
  font-size: 14px;
}

.root-comment-form .comment-input:focus {
  box-shadow: none;
}

.root-comment-form .btn {
  grid-column: 2;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  justify-self: end;
  margin: 0 20px 14px 0;
  padding: 11px 20px;
  border-radius: 13px;
}

.root-comment-form .btn::before {
  content: '↗';
  font-size: 19px;
  line-height: 1;
}

.inline-reply-form {
  margin-top: 12px;
  margin-bottom: 0;
  padding: 12px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  background: rgba(255, 255, 255, 0.03);
}

.reply-banner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
  padding: 10px 12px;
  border-radius: var(--radius-sm);
  background: rgba(45, 212, 191, 0.1);
  color: var(--text-secondary);
  font-size: 14px;
}

.reply-cancel,
.reply-link,
.delete-link {
  color: var(--accent);
  background: transparent;
  border: none;
  cursor: pointer;
  font-size: 13px;
}

.delete-link {
  color: #f87171;
}

.comment-input {
  width: 100%;
  min-height: 100px;
  resize: vertical;
  font-family: inherit;
}

.inline-reply-form .comment-input {
  min-height: 88px;
}

.comment-form .btn {
  margin-top: 12px;
}

.comment-form .btn:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.comment-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.comment-thread {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.comment-item {
  display: flex;
  gap: 12px;
  padding: 16px 0;
  border-bottom: 1px solid var(--border);
}

.comment-item-reply {
  margin-left: 48px;
  padding-top: 12px;
  padding-bottom: 12px;
}

.comment-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--accent-dim);
  color: var(--accent);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  flex-shrink: 0;
}

.comment-body {
  flex: 1;
}

.comment-header {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 8px;
}

.comment-author {
  font-weight: 500;
  font-size: 14px;
}

.reply-to,
.comment-time {
  color: var(--text-muted);
  font-size: 12px;
}

.comment-content {
  color: var(--text-secondary);
  font-size: 14px;
  line-height: 1.6;
}

.reply-list {
  display: flex;
  flex-direction: column;
  margin-left: 48px;
  padding-left: 16px;
  border-left: 2px solid var(--border);
  background: rgba(255, 255, 255, 0.02);
  border-radius: 0 0 var(--radius-sm) var(--radius-sm);
}

.empty-comments {
  text-align: center;
  padding: 24px 0;
  color: var(--text-muted);
  font-size: 14px;
}

.loading-state {
  text-align: center;
  padding: 40px 0;
  color: var(--text-muted);
}

@media (max-width: 768px) {
  .article-cover {
    height: 200px;
  }

  .comment-item-reply {
    margin-left: 0;
  }

  .reply-list {
    margin-left: 20px;
    padding-left: 12px;
  }

  .comment-form-fields {
    grid-template-columns: 64px minmax(0, 1fr);
  }

  .comment-form-fields.is-authenticated {
    grid-template-columns: 64px minmax(0, 1fr);
  }

  .comment-form-avatar {
    width: 40px;
    height: 40px;
    margin-left: 12px;
    font-size: 16px;
  }

  .comment-form-field {
    padding: 0 14px;
  }

  .comment-form-identity {
    padding: 0 14px;
  }

  .comment-form-field:nth-child(n + 3) {
    display: none;
  }

  .root-comment-form .comment-input {
    min-height: 150px;
    padding: 20px;
  }
}
</style>
