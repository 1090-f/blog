<template>
  <div class="admin-layout">
    <aside class="admin-sidebar">
      <div class="sidebar-header">
        <router-link to="/" class="logo">
          <span class="logo-icon">GB</span>
          <span class="logo-text">Gin Blog</span>
        </router-link>
      </div>

      <nav class="sidebar-nav">
        <router-link to="/admin" class="nav-item" :class="{ active: route.path === '/admin' }">
          <span class="nav-icon">DB</span>
          <span>仪表盘</span>
        </router-link>
        <router-link to="/admin/articles" class="nav-item" :class="{ active: route.path.includes('/admin/articles') || route.path.includes('/admin/article/') }">
          <span class="nav-icon">AR</span>
          <span>文章管理</span>
        </router-link>
        <router-link to="/admin/categories" class="nav-item" :class="{ active: route.path.includes('/admin/categories') }">
          <span class="nav-icon">CT</span>
          <span>分类管理</span>
        </router-link>
        <router-link to="/admin/tags" class="nav-item" :class="{ active: route.path.includes('/admin/tags') }">
          <span class="nav-icon">TG</span>
          <span>标签管理</span>
        </router-link>
        <router-link to="/admin/comments" class="nav-item" :class="{ active: route.path.includes('/admin/comments') }">
          <span class="nav-icon">CO</span>
          <span>评论管理</span>
        </router-link>
        <router-link to="/admin/users" class="nav-item" :class="{ active: route.path.includes('/admin/users') }">
          <span class="nav-icon">US</span>
          <span>用户管理</span>
        </router-link>
      </nav>
    </aside>

    <div class="admin-main">
      <header class="admin-header">
        <h2 class="page-title">{{ pageTitle }}</h2>
        <div class="header-right">
          <span class="admin-user">{{ userStore.user?.nickname || userStore.user?.username || 'Admin' }}</span>
          <button class="btn btn-outline btn-sm" @click="handleLogout">退出登录</button>
        </div>
      </header>

      <div class="admin-content">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../stores/user'

const userStore = useUserStore()
const router = useRouter()
const route = useRoute()

const pageTitles = {
  '/admin': '仪表盘',
  '/admin/articles': '文章管理',
  '/admin/categories': '分类管理',
  '/admin/tags': '标签管理',
  '/admin/comments': '评论管理',
  '/admin/users': '用户管理'
}

// 根据当前响应式状态计算派生数据。
const pageTitle = computed(() => {
  if (route.path.includes('/admin/article/new')) return '新建文章'
  if (route.path.includes('/admin/article/edit/')) return '编辑文章'
  return pageTitles[route.path] || '管理后台'
})

// 处理用户操作或浏览器事件。
function handleLogout() {
  userStore.logout()
  router.push('/')
}
</script>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
}

.admin-sidebar {
  width: 240px;
  background: rgba(15, 20, 25, 0.95);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  position: fixed;
  inset: 0 auto 0 0;
  z-index: 100;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid var(--border);
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
}

.logo-icon {
  width: 34px;
  height: 34px;
  border-radius: 10px;
  background: linear-gradient(135deg, #2dd4bf, #3b82f6);
  color: #081018;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
}

.logo-text {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.sidebar-nav {
  flex: 1;
  padding: 16px 12px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  color: var(--text-secondary);
  border-radius: var(--radius-sm);
  transition: all 0.2s ease;
  text-decoration: none;
}

.nav-item:hover {
  background: var(--bg-card);
  color: var(--text-primary);
}

.nav-item.active {
  background: var(--accent-dim);
  color: var(--accent);
}

.nav-icon {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  border: 1px solid var(--border);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 700;
}

.admin-main {
  flex: 1;
  margin-left: 240px;
  display: flex;
  flex-direction: column;
}

.admin-header {
  height: 64px;
  padding: 0 24px;
  background: rgba(15, 20, 25, 0.9);
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: sticky;
  top: 0;
  z-index: 50;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.admin-user {
  color: var(--text-secondary);
  font-size: 14px;
}

.admin-content {
  flex: 1;
  padding: 24px;
}

.btn-sm {
  padding: 6px 12px;
  font-size: 13px;
}

@media (max-width: 768px) {
  .admin-sidebar {
    width: 68px;
  }

  .logo-text,
  .nav-item span:last-child {
    display: none;
  }

  .nav-item {
    justify-content: center;
    padding: 12px;
  }

  .admin-main {
    margin-left: 68px;
  }
}
</style>
