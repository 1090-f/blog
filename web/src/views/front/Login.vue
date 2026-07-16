<template>
  <div class="auth-page">
    <div class="auth-card card">
      <h2 class="auth-title">欢迎回来</h2>
      <p class="auth-subtitle">登录你的账号，继续写作和管理内容</p>
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <input
            v-model="form.username"
            class="input"
            type="text"
            placeholder="用户名"
          />
        </div>
        <div class="form-group">
          <input
            v-model="form.password"
            class="input"
            type="password"
            placeholder="密码"
          />
        </div>
        <button type="submit" class="btn btn-primary btn-block" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </form>
      <div v-if="!isAdminApp" class="auth-footer">
        还没有账号？<router-link to="/register">立即注册</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { adminLogin, login } from '../../api/auth'
import { useUserStore } from '../../stores/user'
import { message } from '../../utils/message'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const isAdminApp = window.__BLOG_APP_MODE__ === 'admin' || import.meta.env.MODE === 'admin'

const form = ref({ username: '', password: '' })
const loading = ref(false)

// 处理用户操作或浏览器事件。
async function handleLogin() {
  if (!form.value.username.trim() || !form.value.password.trim()) {
    message.warning('请输入用户名和密码')
    return
  }

  loading.value = true
  try {
    const loginRequest = isAdminApp ? adminLogin : login
    const res = await loginRequest(form.value)
    const { token, user } = res.data
    userStore.setAuth(token, user)
    message.success('登录成功')
    router.push(route.query.redirect || '/')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-page {
  min-height: 80vh;
  display: flex;
  align-items: center;
  justify-content: center;
}

.auth-card {
  width: 100%;
  max-width: 400px;
  padding: 40px;
}

.auth-title {
  font-size: 24px;
  font-weight: 700;
  text-align: center;
  margin-bottom: 8px;
}

.auth-subtitle {
  text-align: center;
  color: var(--text-muted);
  font-size: 14px;
  margin-bottom: 32px;
}

.form-group {
  margin-bottom: 20px;
}

.btn-block {
  width: 100%;
  padding: 12px;
  font-size: 16px;
}

.auth-footer {
  text-align: center;
  margin-top: 24px;
  color: var(--text-muted);
  font-size: 14px;
}

.auth-footer a {
  color: var(--accent);
}
</style>
