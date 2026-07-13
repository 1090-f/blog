import request from './request'

export const getComments = (articleId) => request.get(`/articles/${articleId}/comments`)
export const createComment = (data) => request.post('/comments', data)
export const deleteMyComment = (id) => request.delete(`/comments/${id}`)

export const getAdminComments = (params) => request.get('/admin/comments', { params })
export const updateCommentStatus = (id, data) => request.put(`/admin/comments/${id}/status`, data)
export const deleteAdminComment = (id) => request.delete(`/admin/comments/${id}`)
