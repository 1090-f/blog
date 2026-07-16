import request from './request'

// 处理当前模块的相关逻辑。
export const getSiteStats = () => request.get('/site-stats')
// 处理当前模块的相关逻辑。
export const getSiteActivity = (params) => request.get('/site-activity', { params })
