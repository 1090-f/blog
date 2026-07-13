import request from './request'

export const getCategories = () => request.get('/categories')
export const createUserCategory = (data) => request.post('/categories', data)
export const createCategory = (data) => request.post('/admin/categories', data)
export const updateCategory = (id, data) => request.put(`/admin/categories/${id}`, data)
export const deleteCategory = (id) => request.delete(`/admin/categories/${id}`)
