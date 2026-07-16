import request from './request'

// 处理当前模块的相关逻辑。
export const getAdminUsers = (params) => request.get('/admin/users', { params })
// 处理当前模块的相关逻辑。
export const updateUserStatus = (id, data) => request.put(`/admin/users/${id}/status`, data)
export const createAdminUser = (data) => request.post('/admin/users/admin', data)
export const updateUserRole = (id, data) => request.put(`/admin/users/${id}/role`, data)
