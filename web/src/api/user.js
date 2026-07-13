import request from './request'

export const getAdminUsers = (params) => request.get('/admin/users', { params })
export const updateUserStatus = (id, data) => request.put(`/admin/users/${id}/status`, data)
