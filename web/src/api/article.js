import request from './request'

export const getArticles = (params) => request.get('/articles', { params })
export const getLatestArticles = (limit) => request.get('/articles/latest', { params: { limit } })
export const getPopularArticles = (limit) => request.get('/articles/popular', { params: { limit } })
export const getArticleDetail = (id) => request.get(`/articles/${id}`)
export const getArticleFull = (id) => request.get(`/articles/${id}/full`)

export const getAdminArticles = (params) => request.get('/admin/articles', { params })
export const createAdminArticle = (data) => request.post('/admin/articles', data)
export const updateArticle = (id, data) => request.put(`/admin/articles/${id}`, data)
export const deleteArticle = (id) => request.delete(`/admin/articles/${id}`)

// 前台用户发布文章
export const createArticle = (data) => request.post('/articles', data)

// 图片上传
export const uploadFile = async (file) => {
  const formData = new FormData()
  formData.append('file', file)
  const response = await request.post('/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
  return response
}
