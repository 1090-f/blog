import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores/user'

const routes = [
  {
    path: '/',
    component: () => import('../layouts/FrontLayout.vue'),
    children: [
      { path: '', name: 'Home', component: () => import('../views/front/Home.vue') },
      { path: 'blog', redirect: '/archive' },
      { path: 'categories', name: 'Categories', component: () => import('../views/front/Categories.vue') },
      { path: 'tags', name: 'Tags', component: () => import('../views/front/Tags.vue') },
      { path: 'archive', name: 'Archive', component: () => import('../views/front/Archive.vue') },
      { path: 'about', name: 'About', component: () => import('../views/front/About.vue') },
      { path: 'article/:id', name: 'Article', component: () => import('../views/front/Article.vue') },
      { path: 'write', name: 'WriteArticle', component: () => import('../views/front/WriteArticle.vue'), meta: { requiresAuth: true } }
    ]
  },
  { path: '/login', name: 'Login', component: () => import('../views/front/Login.vue') },
  { path: '/register', name: 'Register', component: () => import('../views/front/Register.vue') },
  {
    path: '/admin',
    component: () => import('../layouts/AdminLayout.vue'),
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      { path: '', name: 'AdminDashboard', component: () => import('../views/admin/Dashboard.vue') },
      { path: 'articles', name: 'AdminArticles', component: () => import('../views/admin/ArticleList.vue') },
      { path: 'article/new', name: 'AdminArticleNew', component: () => import('../views/admin/ArticleForm.vue') },
      { path: 'article/edit/:id', name: 'AdminArticleEdit', component: () => import('../views/admin/ArticleForm.vue') },
      { path: 'categories', name: 'AdminCategories', component: () => import('../views/admin/CategoryList.vue') },
      { path: 'tags', name: 'AdminTags', component: () => import('../views/admin/TagList.vue') },
      { path: 'comments', name: 'AdminComments', component: () => import('../views/admin/CommentList.vue') },
      { path: 'users', name: 'AdminUsers', component: () => import('../views/admin/UserList.vue') }
    ]
  }
]

if (import.meta.env.MODE === 'admin') {
  routes.unshift({ path: '/', redirect: '/admin' })
}

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to) {
    if (to.hash) {
      return { el: to.hash, top: 80, behavior: 'smooth' }
    }
    return { top: 0 }
  }
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  if (to.meta.requiresAuth && !userStore.token) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else if (to.meta.requiresAdmin && userStore.user?.role !== 'admin') {
    next({ name: 'Home' })
  } else {
    next()
  }
})

export default router
