import request from './request'

export const getProfile = () => request.get('/user/profile')
export const updateProfile = (data) => request.put('/user/profile', data)
export const getAdminUsers = (params) => request.get('/admin/users', { params })
export const updateUserStatus = (id, data) => request.put(`/admin/users/${id}/status`, data)
