import request from './request'

// 处理当前模块的相关逻辑。
export const getArticles = (params) => request.get('/articles', { params })
// 处理当前模块的相关逻辑。
export const getLatestArticles = (limit) => request.get('/articles/latest', { params: { limit } })
// 处理当前模块的相关逻辑。
export const getPopularArticles = (limit) => request.get('/articles/popular', { params: { limit } })
// 处理当前模块的相关逻辑。
export const getArticleDetail = (id) => request.get(`/articles/${id}`)
// 处理当前模块的相关逻辑。
export const getArticleFull = (id) => request.get(`/articles/${id}/full`)

// 处理当前模块的相关逻辑。
export const getAdminArticles = (params) => request.get('/admin/articles', { params })
export const getAdminArticleDetail = (id) => request.get(`/admin/articles/${id}`)
// 处理当前模块的相关逻辑。
export const createAdminArticle = (data) => request.post('/admin/articles', data)
// 处理当前模块的相关逻辑。
export const updateArticle = (id, data) => request.put(`/admin/articles/${id}`, data)
// 处理当前模块的相关逻辑。
export const deleteArticle = (id) => request.delete(`/admin/articles/${id}`)

// 图片上传
export const uploadFile = async (file) => {
  const formData = new FormData()
  formData.append('file', file)
  const response = await request.post('/admin/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
  return response
}
