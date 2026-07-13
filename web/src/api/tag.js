import request from './request'

export const getTags = () => request.get('/tags')
export const createTag = (data) => request.post('/admin/tags', data)
export const updateTag = (id, data) => request.put(`/admin/tags/${id}`, data)
export const deleteTag = (id) => request.delete(`/admin/tags/${id}`)
