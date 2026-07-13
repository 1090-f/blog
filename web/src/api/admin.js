import request from './request'

export const getDashboardStats = () => request.get('/admin/dashboard')
