<template>
  <div class="front-layout">
    <header class="front-header">
      <div class="header-inner">
        <router-link to="/" class="logo">
          <span class="logo-icon">🏠</span>
          <span class="logo-text">Gin Blog</span>
        </router-link>

        <nav class="nav-links">
          <router-link to="/" class="nav-link home-nav-link">🏠 主页</router-link>
          <div
            class="article-nav-group"
            :class="{ 'is-article-menu-open': articleMenuOpen }"
            @mouseenter="articleMenuOpen = true"
            @mouseleave="articleMenuOpen = false"
            @focusin="articleMenuOpen = true"
            @focusout="handleArticleMenuFocusOut"
          >
            <div class="article-nav-trigger-wrap">
              <router-link to="/archive" class="nav-link article-nav-trigger" @click="articleMenuOpen = false">
                📝 文章
                <span class="article-caret article-caret-down">⌄</span>
                <span class="article-caret article-caret-up">⌃</span>
              </router-link>
            </div>
            <div class="article-subnav">
              <router-link to="/archive" class="subnav-link" @click="articleMenuOpen = false">
                <span class="subnav-icon">▣</span>
                <span>归档</span>
              </router-link>
              <router-link to="/categories" class="subnav-link" @click="articleMenuOpen = false">
                <span class="subnav-icon">▰</span>
                <span>分类</span>
              </router-link>
              <router-link to="/tags" class="subnav-link" @click="articleMenuOpen = false">
                <span class="subnav-icon">#</span>
                <span>标签</span>
              </router-link>
            </div>
          </div>

          <template v-if="userStore.isLoggedIn">
            <div
              class="user-menu"
              :class="{ 'is-user-menu-open': userMenuOpen }"
              @mouseenter="userMenuOpen = true"
              @mouseleave="userMenuOpen = false"
              @focusin="userMenuOpen = true"
              @focusout="handleUserMenuFocusOut"
            >
              <button class="user-menu-trigger" type="button" @click="userMenuOpen = !userMenuOpen">
                <span v-if="userStore.user?.avatar" class="user-avatar user-avatar-image">
                  <img :src="userStore.user.avatar" :alt="userStore.user?.nickname || 'user avatar'" />
                </span>
                <span v-else class="user-avatar">{{ userStore.user?.nickname?.[0] || 'U' }}</span>
                <span class="user-name">{{ userStore.user?.nickname || userStore.user?.username }}</span>
                <span class="user-caret user-caret-down"></span>
                <span class="user-caret user-caret-up"></span>
              </button>

              <div class="user-dropdown">
                <button type="button" class="user-dropdown-item subnav-link" @click="goProfile">
                  <span class="subnav-icon">◎</span>
                  个人中心
                </button>
                <button type="button" class="user-dropdown-item subnav-link" @click="handleLogout">
                  <span class="subnav-icon">↪</span>
                  退出登录
                </button>
              </div>
            </div>
          </template>

          <template v-else>
            <router-link to="/login" class="nav-link">🔑 登录</router-link>
          </template>
        </nav>

        <button class="mobile-menu-btn" type="button" @click="mobileOpen = !mobileOpen">☰</button>
      </div>
    </header>

    <div v-if="mobileOpen" class="mobile-menu">
      <router-link to="/" @click="mobileOpen = false">🏠 主页</router-link>
      <router-link to="/archive" @click="mobileOpen = false">📝 文章归档</router-link>

      <template v-if="userStore.isLoggedIn">
        <router-link to="/profile" @click="mobileOpen = false">👤 个人中心</router-link>
        <a @click="handleLogout">🚪 退出登录</a>
      </template>

      <template v-else>
        <router-link to="/login" @click="mobileOpen = false">🔑 登录</router-link>
      </template>
    </div>

    <div class="front-scroll">
      <main class="front-main">
        <div class="front-shell">
          <aside class="front-sidebar front-sidebar-left">
            <div class="card author-card">
              <button class="author-profile-trigger" type="button" aria-label="查看关于我" @click="goAbout">
                <img src="/author-profile.png?v=2" :alt="`${authorName} avatar`" class="author-avatar-image" />
              </button>
              <h3 class="author-name">{{ authorName }}</h3>
              <span class="author-accent" aria-hidden="true"></span>
              <p class="author-description">七月初七，淮水竹亭</p>
              <nav class="author-links" aria-label="个人链接">
                <a href="tencent://AddContact/?fromId=45&fromSubId=1&subcmd=all&uin=1438318243&website=www.oicqzone.com" aria-label="添加 QQ 好友" title="QQ：1438318243">
                  <svg class="author-social-icon qq-icon" viewBox="0 0 24 24" aria-hidden="true">
                    <path d="M12 2.2c-2.57 0-4.32 2.34-4.32 5.47 0 1.13.22 2.1.58 2.93-.93.79-1.64 1.99-1.95 3.46-.18.86.19 1.37.8 1.29.54-.07 1.01-.55 1.35-1.18.32 1.04.94 1.88 1.76 2.38-.26.47-.45.98-.45 1.43 0 .52.31.82.77.54.38-.24.92-.4 1.66-.4s1.28.16 1.66.4c.46.28.77-.02.77-.54 0-.45-.19-.96-.45-1.43.82-.5 1.44-1.34 1.76-2.38.34.63.81 1.11 1.35 1.18.61.08.98-.43.8-1.29-.31-1.47-1.02-2.67-1.95-3.46.36-.83.58-1.8.58-2.93C16.32 4.54 14.57 2.2 12 2.2Z" fill="currentColor" />
                  </svg>
                </a>
                <a href="https://github.com/" target="_blank" rel="noreferrer" aria-label="GitHub">
                  <svg class="author-social-icon" viewBox="0 0 19 19" aria-hidden="true"><use href="/icons.svg#github-icon" /></svg>
                </a>
              </nav>
            </div>

            <div class="card sidebar-meta-card">
              <h4 class="sidebar-card-title">分类</h4>
              <div class="sidebar-tags">
                <button v-for="category in categories" :key="category.id" type="button" @click="goToCategory(category.id)">
                  <span>{{ category.name }}</span>
                  <span class="sidebar-tag-count">{{ category.articleCount || 0 }}</span>
                </button>
              </div>
            </div>

            <div class="card sidebar-meta-card">
              <h4 class="sidebar-card-title">标签</h4>
              <div class="sidebar-tags">
                <button v-for="tag in tags" :key="tag.id" type="button" @click="goToTag(tag.id)">{{ tag.name }}</button>
              </div>
            </div>
          </aside>

          <section class="front-view">
            <ArchiveToolbar
              v-if="showArchiveToolbar"
              :categories="categories"
              :site-stats="siteStats"
            />
            <router-view />
          </section>

          <aside class="front-sidebar front-sidebar-right">
            <WeatherCard />
            <div class="card site-stats-card">
              <h4 class="site-stats-title">站点统计</h4>
              <div class="site-stats-list">
                <div class="site-stat-item"><span class="site-stat-icon">▤</span><span>文章</span><strong>{{ formatNumber(siteStats.articleCount) }}</strong></div>
                <div class="site-stat-item"><span class="site-stat-icon">▱</span><span>分类</span><strong>{{ formatNumber(siteStats.categoryCount) }}</strong></div>
                <div class="site-stat-item"><span class="site-stat-icon">◇</span><span>标签</span><strong>{{ formatNumber(siteStats.tagCount) }}</strong></div>
                <div class="site-stat-item"><span class="site-stat-icon">≡</span><span>总字数</span><strong>{{ formatNumber(siteStats.totalWords) }}</strong></div>
                <div class="site-stat-item"><span class="site-stat-icon">◷</span><span>运行时长</span><strong>{{ formatRuntime(siteStats.firstPublishedAt) }}</strong></div>
                <div class="site-stat-item"><span class="site-stat-icon">〽</span><span>最后活动</span><strong>{{ formatLastActivity(siteStats.lastActivityAt) }}</strong></div>
              </div>
            </div>
            <CalendarCard />
          </aside>
        </div>
      </main>

      <footer class="front-footer">
        <div class="footer-inner">
          <div class="footer-brand">
            <span class="logo-icon">🏠</span>
            <p>Gin Blog - 记录技术与生活</p>
          </div>
          <div class="footer-links">
            <router-link to="/">首页</router-link>
            <router-link to="/archive">文章</router-link>
          </div>
          <div class="footer-copy">
            © {{ new Date().getFullYear() }} Gin Blog. All rights reserved.
          </div>
        </div>
      </footer>
    </div>

    <BlogPet />
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import BlogPet from '../components/BlogPet.vue'
import ArchiveToolbar from '../components/ArchiveToolbar.vue'
import WeatherCard from '../components/WeatherCard.vue'
import CalendarCard from '../components/CalendarCard.vue'
import { getArticles } from '../api/article'
import { getCategories } from '../api/category'
import { getTags } from '../api/tag'
import { getSiteStats } from '../api/site'

const userStore = useUserStore()
const router = useRouter()
const route = useRoute()
const mobileOpen = ref(false)
const articleMenuOpen = ref(false)
const userMenuOpen = ref(false)
const categories = ref([])
const tags = ref([])
const siteStats = ref({
  articleCount: 0,
  categoryCount: 0,
  tagCount: 0,
  totalWords: 0,
  firstPublishedAt: null,
  lastActivityAt: null
})
const authorName = computed(() => userStore.user?.nickname || '管理员')
const showArchiveToolbar = computed(() => ['Home', 'Archive', 'Article', 'Categories', 'Tags'].includes(route.name))

function goProfile() {
  mobileOpen.value = false
  userMenuOpen.value = false
  router.push('/profile')
}

function handleLogout() {
  mobileOpen.value = false
  userMenuOpen.value = false
  userStore.logout()
  router.push('/')
}

function goAbout() {
  router.push('/about')
}

function handleArticleMenuFocusOut(event) {
  if (!event.currentTarget.contains(event.relatedTarget)) {
    articleMenuOpen.value = false
  }
}

function handleUserMenuFocusOut(event) {
  if (!event.currentTarget.contains(event.relatedTarget)) {
    userMenuOpen.value = false
  }
}

function goToCategory(categoryId) {
  router.push({ path: '/archive', query: { categoryId } })
}

function goToTag(tagId) {
  router.push({ path: '/archive', query: { tagId } })
}

function formatNumber(value) {
  return new Intl.NumberFormat('zh-CN').format(value || 0)
}

function formatRuntime(dateStr) {
  if (!dateStr) return '—'
  const days = Math.max(1, Math.floor((Date.now() - new Date(dateStr).getTime()) / 86400000))
  return `${days} 天`
}

function formatLastActivity(dateStr) {
  if (!dateStr) return '—'
  const days = Math.max(0, Math.floor((Date.now() - new Date(dateStr).getTime()) / 86400000))
  return days === 0 ? '今天' : `${days} 天前`
}

async function loadCategoryCounts(rawCategories) {
  return Promise.all(rawCategories.map(async category => {
    try {
      const response = await getArticles({ categoryId: category.id, page: 1, pageSize: 1 })
      return { ...category, articleCount: response.data?.total || 0 }
    } catch {
      return { ...category, articleCount: 0 }
    }
  }))
}

onMounted(async () => {
  try {
    const [categoryRes, tagRes, statsRes] = await Promise.all([getCategories(), getTags(), getSiteStats()])
    categories.value = await loadCategoryCounts(categoryRes.data || [])
    tags.value = tagRes.data || []
    siteStats.value = { ...siteStats.value, ...(statsRes.data || {}) }
  } catch (error) {
    categories.value = []
    tags.value = []
  }
})
</script>

<style scoped>
.front-layout {
  position: relative;
  width: 100%;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: visible;
  background-color: var(--bg-primary);
  background-image:
    linear-gradient(to bottom, rgba(15, 20, 25, 0.16) 0%, rgba(15, 20, 25, 0.34) 62%, var(--bg-primary) 100%),
    url('/blog-background-top.png?v=1');
  background-position: center top;
  background-size: cover;
  background-repeat: no-repeat;
  background-attachment: fixed;
}

.front-header {
  flex: 0 0 60px;
  position: sticky;
  top: 0;
  z-index: 100;
  background: transparent;
  backdrop-filter: none;
  border-bottom: 1px solid transparent;
}

.front-scroll {
  position: relative;
  z-index: 1;
  flex: 1;
  min-height: calc(100vh - 60px);
  display: flex;
  flex-direction: column;
  overflow: visible;
}

.header-inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
}

.logo-icon {
  font-size: 24px;
}

.logo-text {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 8px;
}

.article-nav-group {
  position: relative;
  display: block;
  width: 106px;
}

.article-nav-group::after {
  content: '';
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  height: 8px;
}

.article-nav-trigger {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 7px;
  width: 106px;
  height: 38px;
  padding: 0 10px;
  font-size: 14px;
  text-align: center;
  border-radius: 9px;
  transition: all 0.2s ease;
}

.article-subnav {
  display: none;
  position: absolute;
  top: calc(100% + 8px);
  left: 0;
  width: 195px;
  padding: 7px 7px;
  background: rgba(22, 23, 24, 0.98);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 14px;
  box-shadow: 0 16px 36px rgba(0, 0, 0, 0.38);
  z-index: 120;
}

.subnav-link {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 7px 12px;
  color: #c8c8c8;
  font-size: 14px;
  font-weight: 500;
  border-radius: 12px;
  transition: color 0.2s ease, background 0.2s ease;
}

.subnav-link:hover,
.subnav-link.router-link-active {
  color: #fff;
  background: rgba(255, 255, 255, 0.08);
}

.subnav-icon {
  width: 18px;
  color: #d4d4d4;
  font-size: 16px;
  line-height: 1;
  text-align: center;
}

.article-caret {
  display: inline-block;
  width: 7px;
  height: 7px;
  margin-left: 1px;
  color: var(--text-secondary);
  font-size: 0;
  line-height: 1;
  border-right: 1.5px solid currentColor;
  border-bottom: 1.5px solid currentColor;
  transform: rotate(45deg) translateY(-2px);
}

.article-caret-up {
  display: none;
  transform: rotate(225deg) translateY(-2px);
}

.article-nav-group.is-article-menu-open .article-subnav {
  display: block;
}

.article-nav-group.is-article-menu-open .article-caret-down {
  display: none;
}

.article-nav-group.is-article-menu-open .article-caret-up {
  display: inline-block;
}

.article-nav-group.is-article-menu-open .article-nav-trigger {
  color: var(--accent);
  background: var(--accent-dim);
}

.article-nav-group.is-article-menu-open .article-caret {
  color: var(--accent);
}

.nav-link {
  padding: 8px 16px;
  color: var(--text-secondary);
  font-size: 14px;
  border-radius: var(--radius-sm);
  transition: all 0.2s ease;
}

.nav-link:hover {
  color: var(--text-primary);
  background: var(--bg-card);
}

.nav-link.router-link-active {
  color: var(--accent);
  background: var(--accent-dim);
}

.home-nav-link.router-link-active:not(.router-link-exact-active) {
  color: var(--text-secondary);
  background: transparent;
}

.user-menu {
  position: relative;
  display: flex;
  align-items: center;
  gap: 8px;
  height: 38px;
  padding: 4px 10px;
  background: transparent;
  border-radius: 9px;
  transition: all 0.2s ease;
}

.user-menu::after {
  content: '';
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  height: 8px;
}

.user-menu-trigger {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: transparent;
  border: none;
  padding: 0;
  cursor: pointer;
}

.user-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--accent);
  color: var(--bg-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  overflow: hidden;
}

.user-avatar-image {
  background: transparent;
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.user-name {
  color: var(--text-secondary);
  font-size: 14px;
  transition: color 0.2s ease;
}

.user-menu:hover .user-name,
.user-menu:focus-within .user-name {
  color: var(--accent);
}

.user-dropdown {
  display: none;
  position: absolute;
  top: 100%;
  left: 0;
  margin-top: 8px;
  width: 195px;
  padding: 7px;
  background: rgba(22, 23, 24, 0.98);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 14px;
  box-shadow: 0 16px 36px rgba(0, 0, 0, 0.38);
  overflow: hidden;
  z-index: 20;
}

.user-menu.is-user-menu-open .user-dropdown {
  display: block;
}

.user-menu.is-user-menu-open {
  background: var(--accent-dim);
}

.user-caret {
  display: inline-block;
  width: 7px;
  height: 7px;
  margin-left: 2px;
  color: var(--text-secondary);
  border-right: 1.5px solid currentColor;
  border-bottom: 1.5px solid currentColor;
  transform: rotate(45deg) translateY(-2px);
}

.user-caret-up {
  display: none;
  transform: rotate(225deg) translateY(-2px);
}

.user-menu.is-user-menu-open .user-caret-down {
  display: none;
}

.user-menu.is-user-menu-open .user-caret-up {
  display: inline-block;
}

.user-menu.is-user-menu-open .user-caret {
  color: var(--accent);
}

.user-dropdown-item {
  width: 100%;
  text-align: left;
  background: transparent;
  border: none;
  cursor: pointer;
}

.user-dropdown-item:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
}

.mobile-menu-btn {
  display: none;
  background: none;
  border: none;
  color: var(--text-primary);
  font-size: 24px;
  cursor: pointer;
}

.mobile-menu {
  display: none;
  position: absolute;
  top: 60px;
  left: 0;
  right: 0;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border);
  padding: 16px;
  z-index: 99;
}

.mobile-menu a {
  display: block;
  padding: 12px 16px;
  color: var(--text-secondary);
  border-radius: var(--radius-sm);
  cursor: pointer;
}

.mobile-menu a:hover {
  background: var(--bg-card);
  color: var(--accent);
}

.front-main {
  flex: 1;
  max-width: 1600px;
  margin: 0 auto;
  padding: 24px 20px;
  width: 100%;
}

.front-shell {
  display: grid;
  grid-template-columns: 240px minmax(0, 1fr) 280px;
  gap: 24px;
  align-items: start;
}

.front-view {
  min-width: 0;
}

.front-sidebar {
  display: flex;
  flex-direction: column;
  gap: 20px;
  margin-top: 324px;
}

.author-card {
  padding: 14px 14px 20px;
  text-align: center;
}

.author-profile-trigger {
  display: block;
  width: 100%;
  padding: 0;
  overflow: hidden;
  border: 0;
  border-radius: 18px;
  background: #d9d9d9;
  cursor: pointer;
}

.author-avatar-image {
  display: block;
  width: 100%;
  aspect-ratio: 1;
  object-fit: cover;
  object-position: center 13%;
  transition: transform 0.25s ease;
}

.author-profile-trigger:hover .author-avatar-image {
  transform: scale(1.03);
}

.author-name {
  margin: 18px 0 7px;
  font-size: 20px;
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
  margin-bottom: 18px;
  color: var(--text-secondary);
  font-size: 14px;
}

.author-links {
  display: grid;
  grid-template-columns: repeat(2, 48px);
  justify-content: center;
  gap: 10px;
}

.author-links a {
  display: grid;
  place-items: center;
  aspect-ratio: 1;
  border-radius: 12px;
  background: rgba(45, 212, 191, 0.16);
  color: #c8d7d5;
  font-size: 20px;
}

.author-links a:hover {
  background: var(--accent);
  color: var(--bg-primary);
}

.author-social-icon {
  width: 22px;
  height: 22px;
  color: currentColor;
}

.site-stats-card {
  padding: 22px;
}

.site-stats-title,
.sidebar-card-title {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 20px;
  font-size: 19px;
}

.site-stats-title::before {
  content: '';
  width: 4px;
  height: 21px;
  border-radius: 999px;
  background: var(--accent);
}

.site-stats-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.site-stat-item {
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--text-secondary);
  font-size: 15px;
}

.site-stat-item strong {
  margin-left: auto;
  color: var(--text-primary);
  font-size: 18px;
  white-space: nowrap;
}

.site-stat-icon {
  width: 20px;
  color: var(--accent);
  font-size: 21px;
  font-weight: 700;
  line-height: 1;
  text-align: center;
}

.sidebar-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.sidebar-tags button {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  border: 0;
  border-radius: 999px;
  padding: 6px 10px;
  background: var(--accent-dim);
  color: var(--accent);
  font-size: 12px;
  cursor: pointer;
}

.sidebar-tag-count {
  min-width: 1.25em;
  padding: 1px 4px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.1);
  font-size: 11px;
  font-weight: 600;
  line-height: 1.2;
  text-align: center;
}

.sidebar-tags button:hover {
  background: var(--accent);
  color: var(--bg-primary);
}

.front-footer {
  background: transparent;
  border-top: 1px solid transparent;
  padding: 40px 20px;
  margin-top: auto;
}

.footer-inner {
  max-width: 1600px;
  margin: 0 auto;
  text-align: center;
}

.footer-brand {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-bottom: 20px;
}

.footer-brand p {
  color: var(--text-secondary);
  font-size: 14px;
}

.footer-links {
  display: flex;
  justify-content: center;
  gap: 24px;
  margin-bottom: 20px;
}

.footer-links a {
  color: var(--text-muted);
  font-size: 14px;
}

.footer-links a:hover {
  color: var(--accent);
}

.footer-copy {
  color: var(--text-muted);
  font-size: 13px;
}

@media (max-width: 768px) {
  .nav-links {
    display: none;
  }

  .mobile-menu-btn {
    display: block;
  }

  .mobile-menu {
    display: block;
  }

  .article-nav-group {
    display: contents;
    width: auto;
  }

  .article-subnav {
    display: none;
  }
}

@media (max-width: 1200px) {
  .front-shell {
    grid-template-columns: 220px minmax(0, 1fr) 240px;
  }
}

@media (max-width: 1024px) {
  .front-shell {
    grid-template-columns: minmax(0, 1fr);
  }

  .front-sidebar {
    display: none;
  }
}
</style>
