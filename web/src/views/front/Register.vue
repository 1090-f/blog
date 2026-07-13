<template>
  <div class="auth-page">
    <div class="auth-card card">
      <h2 class="auth-title">创建账号</h2>
      <p class="auth-subtitle">注册一个新账号，开始发布你的内容</p>
      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <label class="form-label">用户名</label>
          <input
            v-model="form.username"
            class="input"
            type="text"
            placeholder="3-50 位字符，支持字母、数字和下划线"
            maxlength="50"
          />
          <span class="form-hint">用户名用于登录，注册后通常不再修改。</span>
        </div>
        <div class="form-group">
          <label class="form-label">昵称</label>
          <input
            v-model="form.nickname"
            class="input"
            type="text"
            placeholder="2-50 位字符"
            maxlength="50"
          />
          <span class="form-hint">昵称会显示在文章和评论中。</span>
        </div>
        <div class="form-group">
          <label class="form-label">密码</label>
          <input
            v-model="form.password"
            class="input"
            type="password"
            placeholder="至少 6 位"
            maxlength="20"
          />
          <span class="form-hint">建议同时包含字母和数字。</span>
        </div>
        <div class="form-group">
          <label class="form-label">确认密码</label>
          <input
            v-model="form.confirmPassword"
            class="input"
            type="password"
            placeholder="请再次输入密码"
          />
          <span v-if="form.confirmPassword && form.password !== form.confirmPassword" class="form-error">两次输入的密码不一致。</span>
        </div>
        <button type="submit" class="btn btn-primary btn-block" :disabled="loading">
          {{ loading ? '注册中...' : '注册' }}
        </button>
      </form>
      <div class="auth-footer">
        已有账号？<router-link to="/login">立即登录</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { register } from '../../api/auth'
import { message } from '../../utils/message'

const router = useRouter()
const form = ref({ username: '', nickname: '', password: '', confirmPassword: '' })
const loading = ref(false)

async function handleRegister() {
  if (!form.value.username.trim() || !form.value.nickname.trim() || !form.value.password.trim()) {
    message.warning('请填写完整的注册信息')
    return
  }

  if (form.value.password !== form.value.confirmPassword) {
    message.warning('两次输入的密码不一致')
    return
  }

  loading.value = true
  try {
    await register({
      username: form.value.username,
      nickname: form.value.nickname,
      password: form.value.password
    })
    message.success('注册成功，请登录')
    router.push('/login')
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

.form-label {
  display: block;
  margin-bottom: 6px;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
}

.form-hint {
  display: block;
  margin-top: 6px;
  font-size: 12px;
  color: var(--text-muted);
}

.form-error {
  display: block;
  margin-top: 6px;
  font-size: 12px;
  color: #ef4444;
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
