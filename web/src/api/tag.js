import request from './request'

// 处理当前模块的相关逻辑。
export const getTags = () => request.get('/tags')
export const getAdminTags = () => request.get('/admin/tags')
// 处理当前模块的相关逻辑。
export const createTag = (data) => request.post('/admin/tags', data)
// 处理当前模块的相关逻辑。
export const updateTag = (id, data) => request.put(`/admin/tags/${id}`, data)
// 处理当前模块的相关逻辑。
export const deleteTag = (id) => request.delete(`/admin/tags/${id}`)
