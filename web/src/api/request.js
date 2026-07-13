import axios from 'axios'
import { useUserStore } from '../stores/user'
import { message } from '../utils/message'
import router from '../router'

const request = axios.create({
  baseURL: '/api',
  timeout: 10000
})

function isAdminPage() {
  return import.meta.env.MODE === 'admin' || window.location.pathname === '/admin' || window.location.pathname.startsWith('/admin/')
}

let isHandlingExpiredToken = false

function handleExpiredToken() {
  if (isHandlingExpiredToken) {
    return
  }

  isHandlingExpiredToken = true

  const userStore = useUserStore()
  const redirect = router.currentRoute.value.fullPath

  userStore.logout()
  message.warning('登录已过期，请重新登录')

  const isLoginPage = router.currentRoute.value.name === 'Login'
  if (!isLoginPage) {
    router.replace({ name: 'Login', query: { redirect } })
  }

  window.setTimeout(() => {
    isHandlingExpiredToken = false
  }, 0)
}

request.interceptors.request.use(config => {
  config.baseURL = isAdminPage() ? '/admin-api' : '/api'
  const userStore = useUserStore()
  if (userStore.token) {
    config.headers.Authorization = `Bearer ${userStore.token}`
  }
  return config
})

request.interceptors.response.use(
  res => {
    if (res.data.code !== 0) {
      message.error(res.data.message || 'Request failed')
      return Promise.reject(res.data)
    }
    return res.data
  },
  err => {
    const status = err.response?.status
    const code = err.response?.data?.code
    const userStore = useUserStore()

    if (status === 401 && code === 4010 && userStore.token) {
      handleExpiredToken()
      return Promise.reject(err)
    }

    const msg = err.response?.data?.message || 'Network error'
    message.error(msg)
    return Promise.reject(err)
  }
)

export default request
