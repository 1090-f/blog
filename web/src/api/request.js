import axios from 'axios'
import { useUserStore } from '../stores/user'
import { message } from '../utils/message'
import router from '../router'

const request = axios.create({
  baseURL: '/api',
  timeout: 10000
})

let isHandlingInvalidSession = false

function localizeErrorMessage(messageText) {
  const messages = {
    'invalid username or password': '用户名或密码错误',
    'user is disabled': '账号已被禁用，请联系管理员',
    'network error': '网络连接失败',
    'request failed': '请求失败'
  }
  return messages[messageText?.toLowerCase()] || messageText
}

// 处理用户操作或浏览器事件。
function handleInvalidSession(reason) {
  if (isHandlingInvalidSession) {
    return
  }

  isHandlingInvalidSession = true

  const userStore = useUserStore()
  const redirect = router.currentRoute.value.fullPath

  userStore.logout()
  message.warning(reason)

  const isLoginPage = router.currentRoute.value.name === 'Login'
  if (!isLoginPage) {
    router.replace({ name: 'Login', query: { redirect } })
  }

  window.setTimeout(() => {
    isHandlingInvalidSession = false
  }, 0)
}

request.interceptors.request.use(config => {
  // The public and admin applications are served by different origins in
  // production, so both can use their own same-origin /api endpoint.
  config.baseURL = '/api'
  const userStore = useUserStore()
  if (userStore.token) {
    config.headers.Authorization = `Bearer ${userStore.token}`
  }
  return config
})

request.interceptors.response.use(
  res => {
    if (res.data.code !== 0) {
      message.error(localizeErrorMessage(res.data.message) || '请求失败')
      return Promise.reject(res.data)
    }
    return res.data
  },
  err => {
    const status = err.response?.status
    const code = err.response?.data?.code
    const userStore = useUserStore()

    if (userStore.token && status === 401 && code === 4010) {
      handleInvalidSession('登录已过期，请重新登录')
      return Promise.reject(err)
    }

    if (userStore.token && status === 403 && code === 4031) {
      handleInvalidSession('账号已被禁用，已自动退出登录')
      return Promise.reject(err)
    }

    const msg = err.response?.data?.message || '网络连接失败'
    if (!err.config?.skipErrorMessage) {
      message.error(localizeErrorMessage(msg))
    }
    return Promise.reject(err)
  }
)

export default request
