<template>
  <div class="home">
    <div v-if="!aboutOpen" class="banner">
      <div class="banner-content">
        <p class="banner-subtitle">记录技术与生活，沉淀经验，也分享正在发生的思考。</p>
        <router-link to="/archive" class="btn btn-primary">浏览文章</router-link>
      </div>
    </div>

    <div class="home-grid">
      <aside class="sidebar sidebar-left">
        <div class="card author-card" data-pet-hint="个人资料">
          <div class="author-avatar-frame">
            <img
              src="/author-profile.png?v=2"
              :alt="`${authorName} avatar`"
              class="author-avatar-image"
            />
            <button
              class="author-profile-trigger"
              type="button"
              :aria-label="aboutOpen ? '返回最新文章' : '查看关于我'"
              @click="toggleAbout"
            >
              <svg viewBox="0 0 48 48" aria-hidden="true">
                <rect x="5" y="8" width="38" height="32" rx="5" />
                <circle cx="17" cy="20" r="4" />
                <path d="M11 33c1.8-5.3 10.2-5.3 12 0M29 18h8M29 25h8M29 32h5" />
              </svg>
            </button>
          </div>
          <h3 class="author-name">{{ authorName }}</h3>
          <span class="author-accent" aria-hidden="true"></span>
          <p class="author-description">{{ authorDescription }}</p>

          <nav class="author-links" aria-label="个人链接">
            <a
              v-for="link in authorLinks"
              :key="link.label"
              :href="link.href"
              class="author-link"
              :class="`author-link-${link.icon}`"
              :aria-label="link.label"
              :title="link.label"
              :target="link.external ? '_blank' : undefined"
              :rel="link.external ? 'noreferrer' : undefined"
            >
              <svg v-if="link.icon === 'bilibili'" viewBox="0 0 24 24" aria-hidden="true">
                <rect x="3" y="7" width="18" height="13" rx="4" />
                <path d="M8 4 6 2M16 4l2-2M8 13h.01M16 13h.01M8 17c2.2 1.2 5.8 1.2 8 0" />
              </svg>
              <svg v-else-if="link.icon === 'github'" viewBox="0 0 24 24" aria-hidden="true">
                <path d="M12 3a9 9 0 0 0-2.8 17.55c.45.08.62-.2.62-.44v-1.55c-2.53.55-3.06-1.07-3.06-1.07-.4-1.02-.98-1.29-.98-1.29-.8-.55.06-.54.06-.54.89.07 1.36.92 1.36.92.79 1.36 2.07.97 2.58.74.08-.57.31-.97.56-1.19-2.02-.23-4.15-1.01-4.15-4.5 0-.99.35-1.8.92-2.43-.09-.22-.4-1.15.09-2.4 0 0 .75-.24 2.47.93A8.6 8.6 0 0 1 12 7.4c.76 0 1.52.1 2.23.33 1.72-1.17 2.47-.93 2.47-.93.49 1.25.18 2.18.09 2.4.57.63.92 1.44.92 2.43 0 3.5-2.13 4.26-4.16 4.49.32.28.6.82.6 1.65v2.45c0 .24.16.52.62.44A9 9 0 0 0 12 3Z" />
              </svg>
              <svg v-else-if="link.icon === 'mail'" viewBox="0 0 24 24" aria-hidden="true">
                <rect x="3" y="5" width="18" height="14" rx="3" />
                <path d="m5 8 7 5 7-5" />
              </svg>
              <svg v-else viewBox="0 0 24 24" aria-hidden="true">
                <path d="M5 5h.01M5 12a7 7 0 0 1 7 7M5 5a14 14 0 0 1 14 14M5 19h.01" />
              </svg>
            </a>
          </nav>
          <p class="author-desc">一个面向写作、归档与分享的轻量博客系统。</p>
        </div>

        <div class="card">
          <h4 class="card-title">站点说明</h4>
          <p class="card-content">这里展示最新文章、分类入口和基础统计，适合作为博客首页和内容导航页。</p>
        </div>

        <div class="card">
          <h4 class="card-title">分类</h4>
          <div class="tag-cloud">
            <span v-for="cat in categories" :key="`tag-${cat.id}`" class="tag" @click="goToCategory(cat.id)">
              {{ cat.name }}
            </span>
            <span v-if="categories.length === 0" class="text-muted">暂无分类</span>
          </div>
        </div>

        <div class="card">
          <h4 class="card-title">标签</h4>
          <div class="tag-cloud">
            <span v-for="tag in tags" :key="`tag-${tag.id}`" class="tag" @click="goToTag(tag.id)">
              {{ tag.name }}
            </span>
            <span v-if="tags.length === 0" class="text-muted">暂无标签</span>
          </div>
        </div>
      </aside>

      <main class="main-content">
        <div v-if="aboutOpen" class="card about-inline-card">
          <About />
        </div>
        <div v-else class="card">
          <div class="section-tabs" role="tablist" aria-label="文章排序方式">
            <button
              type="button"
              class="section-tab"
              :class="{ active: articleListType === 'latest' }"
              :aria-selected="articleListType === 'latest'"
              role="tab"
              @click="articleListType = 'latest'"
            >
              最新文章
            </button>
            <button
              type="button"
              class="section-tab"
              :class="{ active: articleListType === 'popular' }"
              :aria-selected="articleListType === 'popular'"
              role="tab"
              @click="articleListType = 'popular'"
            >
              最热文章
            </button>
          </div>
          <div v-if="articles.length > 0" class="article-list">
            <div v-for="article in articles" :key="article.id" class="article-item" :data-pet-hint="`文章：${article.title}`" @click="goToArticle(article.id)">
              <span class="article-date">{{ formatDate(article.createdAt) }}</span>
              <span class="article-category">{{ article.category?.name || '未分类' }}</span>
              <span class="article-title">{{ article.title }}</span>
            </div>
          </div>
          <div v-else class="empty-state">
            <p>暂无文章</p>
          </div>
        </div>
      </main>

      <aside class="sidebar sidebar-right">
        <div class="card site-stats-card">
          <h4 class="card-title site-stats-title">站点统计</h4>
          <div class="stats-list">
            <div class="stat-item">
              <span class="stat-icon">▤</span>
              <span class="stat-label">文章</span>
              <span class="stat-value">{{ formatNumber(siteStats.articleCount) }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-icon">▱</span>
              <span class="stat-label">分类</span>
              <span class="stat-value">{{ formatNumber(siteStats.categoryCount) }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-icon">◇</span>
              <span class="stat-label">标签</span>
              <span class="stat-value">{{ formatNumber(siteStats.tagCount) }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-icon">≡</span>
              <span class="stat-label">总字数</span>
              <span class="stat-value">{{ formatNumber(siteStats.totalWords) }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-icon">◷</span>
              <span class="stat-label">运行时长</span>
              <span class="stat-value">{{ formatRuntime(siteStats.firstPublishedAt) }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-icon">〽</span>
              <span class="stat-label">最后活动</span>
              <span class="stat-value">{{ formatLastActivity(siteStats.lastActivityAt) }}</span>
            </div>
          </div>
        </div>
      </aside>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getLatestArticles, getPopularArticles } from '../../api/article'
import { getCategories } from '../../api/category'
import { getTags } from '../../api/tag'
import { getSiteStats } from '../../api/site'
import { useUserStore } from '../../stores/user'
import About from './About.vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const articleListType = ref('latest')
const latestArticles = ref([])
const popularArticles = ref([])
const categories = ref([])
const tags = ref([])
const articles = computed(() => articleListType.value === 'popular' ? popularArticles.value : latestArticles.value)
// 根据当前响应式状态计算派生数据。
const aboutOpen = computed(() => route.query.panel === 'about')
const siteStats = ref({
  articleCount: 0,
  categoryCount: 0,
  tagCount: 0,
  totalWords: 0,
  firstPublishedAt: null,
  lastActivityAt: null
})

const authorName = userStore.user?.nickname || 'Jason Shane'
const authorDescription = '\u4e03\u6708\u521d\u4e03\uff0c\u6dee\u6c34\u7af9\u4ead'
const authorLinks = [
  { label: 'Bilibili', icon: 'bilibili', href: 'https://www.bilibili.com/', external: true },
  { label: 'GitHub', icon: 'github', href: 'https://github.com/', external: true },
  { label: '\u53d1\u9001\u90ae\u4ef6', icon: 'mail', href: 'mailto:hello@ginblog.dev', external: false },
  { label: 'RSS \u8ba2\u9605', icon: 'rss', href: '/rss.xml', external: false }
]
/*
const authorDescription = '七月初七，淮水竹亭'
const authorLinks = [
  { label: 'Bilibili', icon: 'bilibili', href: 'https://www.bilibili.com/', external: true },
  { label: 'GitHub', icon: 'github', href: 'https://github.com/', external: true },
  { label: '发送邮件', icon: 'mail', href: 'mailto:hello@ginblog.dev', external: false },
  { label: 'RSS 订阅', icon: 'rss', href: '/rss.xml', external: false }
]
*/

// 将原始数据格式化为界面展示内容。
function formatDate(dateStr) {
  const d = new Date(dateStr)
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${month}-${day}`
}

// 将原始数据格式化为界面展示内容。
function formatNumber(value) {
  return new Intl.NumberFormat('zh-CN').format(value || 0)
}

// 将原始数据格式化为界面展示内容。
function formatRuntime(dateStr) {
  if (!dateStr) return '—'
  const days = Math.max(1, Math.floor((Date.now() - new Date(dateStr).getTime()) / 86400000))
  return `${days} 天`
}

// 将原始数据格式化为界面展示内容。
function formatLastActivity(dateStr) {
  if (!dateStr) return '—'
  const days = Math.max(0, Math.floor((Date.now() - new Date(dateStr).getTime()) / 86400000))
  if (days === 0) return '今天'
  return `${days} 天前`
}

// 跳转到对应页面。
function goToArticle(id) {
  router.push(`/article/${id}`)
}

// 跳转到对应页面。
function goToCategory(categoryId) {
  router.push({ path: '/archive', query: { categoryId } })
}

// 跳转到对应页面。
function goToTag(tagId) {
  router.push({ path: '/archive', query: { tagId } })
}

// 切换对应的界面状态。
function toggleAbout() {
  const query = { ...route.query }
  if (aboutOpen.value) {
    delete query.panel
  } else {
    query.panel = 'about'
  }
  router.push({ name: 'Home', query })
}

onMounted(async () => {
  try {
    const [latestRes, popularRes, categoriesRes, tagsRes, statsRes] = await Promise.all([
      getLatestArticles(10),
      getPopularArticles(10),
      getCategories(),
      getTags(),
      getSiteStats()
    ])
    latestArticles.value = latestRes.data || []
    popularArticles.value = popularRes.data || []
    categories.value = categoriesRes.data || []
    tags.value = tagsRes.data || []
    siteStats.value = { ...siteStats.value, ...(statsRes.data || {}) }
  } catch (error) {
    latestArticles.value = []
    popularArticles.value = []
    categories.value = []
    tags.value = []
  }
})
</script>

<style scoped>
.banner {
  position: relative;
  height: 300px;
  background: linear-gradient(135deg, rgba(45, 212, 191, 0.2), rgba(15, 20, 25, 0.8));
  border-radius: var(--radius-lg);
  overflow: hidden;
  margin-bottom: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.banner::before {
  display: none;
}

.banner-content {
  position: relative;
  text-align: center;
  z-index: 1;
}

.banner-subtitle {
  color: var(--text-secondary);
  font-size: 18px;
  margin-bottom: 24px;
}

.home-grid {
  display: block;
}

.home-grid > .sidebar {
  display: none;
}

.sidebar {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 16px;
}

.card-content {
  color: var(--text-secondary);
  font-size: 14px;
  line-height: 1.7;
}

.about-inline-card {
  padding: 0;
  overflow: hidden;
}

.about-inline-card :deep(.about-page) {
  max-width: none;
  padding: 30px;
}

.about-inline-card :deep(.about-section-lead) {
  font-size: 17px;
}

.about-inline-card :deep(.about-section) {
  margin-top: 48px;
}

.about-inline-card :deep(.project-card) {
  padding: 24px;
}

@media (max-width: 768px) {
  .about-inline-card :deep(.about-page) {
    padding: 22px 18px;
  }
}

.author-card {
  text-align: center;
  padding: 14px 14px 20px;
  overflow: hidden;
}

.author-avatar-frame {
  position: relative;
  width: 100%;
  aspect-ratio: 1;
  margin-bottom: 18px;
  overflow: hidden;
  border-radius: 18px;
  background: #d9d9d9;
}

.author-avatar-image {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
  object-position: center 13%;
  transition: transform 0.35s ease, filter 0.35s ease;
}

.author-avatar-frame:hover .author-avatar-image,
.author-avatar-frame:focus-within .author-avatar-image {
  transform: scale(1.03);
  filter: brightness(0.68);
}

.author-profile-trigger {
  position: absolute;
  inset: 0;
  display: grid;
  place-items: center;
  padding: 0;
  border: 0;
  background: rgba(15, 20, 25, 0.24);
  color: #fff;
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.25s ease, background 0.25s ease;
}

.author-avatar-frame:hover .author-profile-trigger,
.author-profile-trigger:focus-visible {
  opacity: 1;
  background: rgba(15, 20, 25, 0.4);
}

.author-profile-trigger svg {
  width: 54px;
  height: 54px;
  fill: none;
  stroke: currentColor;
  stroke-width: 2.6;
  stroke-linecap: round;
  stroke-linejoin: round;
  filter: drop-shadow(0 2px 5px rgba(0, 0, 0, 0.28));
}

.author-avatar {
  width: 80px;
  height: 80px;
  margin: 0 auto 16px;
  background: linear-gradient(135deg, var(--accent-dim), rgba(45, 212, 191, 0.3));
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  font-weight: 700;
}

.author-name {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 7px;
}

.author-accent {
  display: block;
  width: 38px;
  height: 6px;
  margin: 0 auto 12px;
  border-radius: 999px;
  background: var(--accent);
}

.author-description {
  color: var(--text-secondary);
  font-size: 14px;
  margin-bottom: 18px;
}

.author-card > .author-desc {
  display: none;
}

.author-links {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 10px;
}

.author-link {
  display: grid;
  place-items: center;
  width: 100%;
  aspect-ratio: 1;
  border-radius: 12px;
  background: rgba(45, 212, 191, 0.16);
  color: #c8d7d5;
  transition: color 0.2s ease, background 0.2s ease, transform 0.2s ease;
}

.author-link:hover,
.author-link:focus-visible {
  color: var(--bg-primary);
  background: var(--accent);
  transform: translateY(-2px);
}

.author-link svg {
  width: 27px;
  height: 27px;
  fill: none;
  stroke: currentColor;
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.author-link-github svg {
  fill: currentColor;
  stroke: none;
}

.author-desc {
  color: var(--text-muted);
  font-size: 13px;
}

.section-tabs {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 20px;
  border-bottom: 1px solid var(--border);
}

.section-tab {
  position: relative;
  padding: 0 0 12px;
  border: 0;
  background: transparent;
  color: var(--text-muted);
  font-size: 18px;
  font-weight: 600;
  cursor: pointer;
  transition: color 0.2s ease;
}

.section-tab::after {
  content: '';
  position: absolute;
  right: 0;
  bottom: -1px;
  left: 0;
  height: 2px;
  border-radius: 999px;
  background: var(--accent);
  opacity: 0;
  transform: scaleX(0.45);
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.section-tab:hover,
.section-tab.active {
  color: var(--text-primary);
}

.section-tab.active::after {
  opacity: 1;
  transform: scaleX(1);
}

.article-list {
  display: flex;
  flex-direction: column;
}

.article-item {
  display: grid;
  grid-template-columns: 80px 100px 1fr;
  gap: 12px;
  padding: 14px 0;
  border-bottom: 1px solid var(--border);
  cursor: pointer;
  transition: background-color 0.2s ease, padding-left 0.2s ease, padding-right 0.2s ease;
}

.article-item:hover {
  padding-left: 12px;
  padding-right: 12px;
  background: var(--accent-dim);
  border-radius: var(--radius-sm);
}

.article-item:hover .article-title {
  color: var(--accent);
}

.article-item:last-child {
  border-bottom: none;
}

.article-date {
  color: var(--text-muted);
  font-size: 13px;
}

.article-category {
  color: var(--accent);
  font-size: 13px;
}

.article-title {
  color: var(--text-primary);
  font-size: 14px;
}

.empty-state {
  text-align: center;
  padding: 40px 0;
  color: var(--text-muted);
}

.empty-state .btn {
  margin-top: 16px;
}

.stats-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.site-stats-card {
  padding: 22px;
}

.site-stats-title {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 22px;
  font-size: 20px;
  color: var(--text-primary);
}

.site-stats-title::before {
  content: '';
  width: 4px;
  height: 22px;
  border-radius: 999px;
  background: var(--accent);
}

.stat-icon {
  width: 20px;
  color: var(--accent);
  font-size: 22px;
  font-weight: 700;
  line-height: 1;
  text-align: center;
}

.stat-label {
  flex: 1;
  color: var(--text-secondary);
  font-size: 15px;
}

.stat-value {
  color: var(--text-primary);
  font-size: 18px;
  font-weight: 700;
  white-space: nowrap;
}

.tag-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.text-muted {
  color: var(--text-muted);
  font-size: 13px;
}

@media (max-width: 1024px) {
  .home-grid {
    grid-template-columns: 1fr;
  }

  .sidebar-right {
    display: none;
  }
}

@media (max-width: 768px) {
  .banner {
    height: 200px;
  }

  .banner-subtitle {
    font-size: 14px;
  }

  .sidebar-left {
    display: none;
  }

  .article-item {
    grid-template-columns: 70px 84px 1fr;
  }
}
</style>
