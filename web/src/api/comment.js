import request from './request'

// 处理当前模块的相关逻辑。
export const getComments = (articleId) => request.get(`/articles/${articleId}/comments`)
// 处理当前模块的相关逻辑。
export const createComment = (data) => request.post('/comments', data)
// 处理当前模块的相关逻辑。
export const deleteMyComment = (id) => request.delete(`/comments/${id}`)

// 处理当前模块的相关逻辑。
export const getAdminComments = (params) => request.get('/admin/comments', { params })
// 处理当前模块的相关逻辑。
export const updateCommentStatus = (id, data) => request.put(`/admin/comments/${id}/status`, data)
// 处理当前模块的相关逻辑。
export const deleteAdminComment = (id) => request.delete(`/admin/comments/${id}`)
