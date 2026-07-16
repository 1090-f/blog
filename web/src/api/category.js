import request from './request'

// 处理当前模块的相关逻辑。
export const getCategories = () => request.get('/categories')
export const getAdminCategories = () => request.get('/admin/categories')
// 处理当前模块的相关逻辑。
export const createUserCategory = (data) => request.post('/categories', data)
// 处理当前模块的相关逻辑。
export const createCategory = (data) => request.post('/admin/categories', data)
// 处理当前模块的相关逻辑。
export const updateCategory = (id, data) => request.put(`/admin/categories/${id}`, data)
// 处理当前模块的相关逻辑。
export const deleteCategory = (id) => request.delete(`/admin/categories/${id}`)
