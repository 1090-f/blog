import request from './request'

// 处理当前模块的相关逻辑。
export const getDashboardStats = () => request.get('/admin/dashboard')
