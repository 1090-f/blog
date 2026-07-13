import request from './request'

export const getSiteStats = () => request.get('/site-stats')
export const getSiteActivity = (params) => request.get('/site-activity', { params })
