import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')

  let savedUser = null
  try {
    const userStr = localStorage.getItem('user')
    if (userStr && userStr !== 'undefined') {
      savedUser = JSON.parse(userStr)
    }
  } catch (e) { /* ignore */ }
  const user = ref(savedUser)

  // 根据当前响应式状态计算派生数据。
  const isLoggedIn = computed(() => !!token.value)
  // 根据当前响应式状态计算派生数据。
  const isAdmin = computed(() => user.value?.role === 'admin')

  // 更新对应的状态值。
  function setAuth(t, u) {
    token.value = t
    user.value = u
    localStorage.setItem('token', t)
    localStorage.setItem('user', JSON.stringify(u))
  }

  // 清除当前用户的登录状态。
  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  return { token, user, isLoggedIn, isAdmin, setAuth, logout }
})
