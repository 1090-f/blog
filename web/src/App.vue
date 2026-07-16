<template>
  <div class="dark">
    <router-view />
    <AppMessages />
  </div>
</template>

<script setup>
import { onBeforeUnmount, onMounted } from 'vue'
import AppMessages from './components/AppMessages.vue'
import { getAdminCurrentUser, getCurrentUser } from './api/auth'
import { useUserStore } from './stores/user'

const userStore = useUserStore()
const SESSION_CHECK_INTERVAL = 15000
let sessionCheckTimer = null
let isCheckingSession = false

function currentUserRequest() {
  const isAdminApp = window.__BLOG_APP_MODE__ === 'admin' || import.meta.env.MODE === 'admin'
  return isAdminApp ? getAdminCurrentUser : getCurrentUser
}

async function checkSession() {
  if (!userStore.token || isCheckingSession) {
    return
  }

  const checkedToken = userStore.token
  isCheckingSession = true
  try {
    const res = await currentUserRequest()({ skipErrorMessage: true })
    if (userStore.token === checkedToken) {
      userStore.setAuth(checkedToken, res.data)
    }
  } catch (error) {
    // 令牌过期或账号禁用由响应拦截器统一清理登录态。
  } finally {
    isCheckingSession = false
  }
}

function checkSessionWhenVisible() {
  if (document.visibilityState === 'visible') {
    checkSession()
  }
}

onMounted(() => {
  checkSession()
  sessionCheckTimer = window.setInterval(checkSession, SESSION_CHECK_INTERVAL)
  window.addEventListener('focus', checkSession)
  document.addEventListener('visibilitychange', checkSessionWhenVisible)
})

onBeforeUnmount(() => {
  window.clearInterval(sessionCheckTimer)
  window.removeEventListener('focus', checkSession)
  document.removeEventListener('visibilitychange', checkSessionWhenVisible)
})
</script>
