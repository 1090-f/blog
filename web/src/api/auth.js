import request from './request'

// 处理当前模块的相关逻辑。
export const login = (data) => request.post('/auth/login', data)
export const adminLogin = (data) => request.post('/admin/login', data)
// 处理当前模块的相关逻辑。
export const register = (data) => request.post('/auth/register', data)
// 处理当前模块的相关逻辑。
export const getCurrentUser = (config) => request.get('/user/session', config)
export const getAdminCurrentUser = (config) => request.get('/admin/session', config)
