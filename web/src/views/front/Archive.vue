<template>
  <div class="archive-page">
    <section class="card archive-card">
      <header class="archive-heading">
        <h1><span>分类</span><i>/</i><strong>{{ archiveTitle }}</strong></h1>
        <span>{{ total }} 篇文章</span>
      </header>

      <div v-if="loading" class="loading-state">加载中...</div>
      <div v-else-if="Object.keys(groupedArticles).length === 0" class="empty-state">暂无文章</div>
      <div v-else class="archive-groups">
        <section v-for="(items, year) in groupedArticles" :key="year" class="archive-group">
          <h2 class="archive-year">{{ year }}</h2>
          <div class="archive-timeline">
            <div class="archive-year-summary">{{ items.length }} 篇文章</div>
            <router-link v-for="article in items" :key="article.id" :to="`/article/${article.id}`" class="archive-item">
              <time class="archive-date">{{ formatDate(article.createdAt) }}</time>
              <span class="archive-dot" aria-hidden="true"></span>
              <span class="archive-item-content">
                <span class="archive-category">{{ article.category?.name || '未分类' }}</span>
                <span class="archive-title">{{ article.title }}</span>
              </span>
            </router-link>
          </div>
        </section>
      </div>
    </section>

    <div v-if="total > pageSize" class="pagination-wrap">
      <button class="btn btn-outline" :disabled="page <= 1" @click="changePage(page - 1)">上一页</button>
      <span class="page-info">{{ page }} / {{ Math.ceil(total / pageSize) }}</span>
      <button class="btn btn-outline" :disabled="page >= Math.ceil(total / pageSize)" @click="changePage(page + 1)">下一页</button>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getArticles } from '../../api/article'
import { getCategories } from '../../api/category'
import { getTags } from '../../api/tag'

const route = useRoute()
const articles = ref([])
const categories = ref([])
const tags = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = 20
const total = ref(0)
const selectedCategory = ref(null)
const selectedTag = ref(null)

const groupedArticles = computed(() => articles.value.reduce((groups, article) => {
  const year = new Date(article.createdAt).getFullYear()
  if (!groups[year]) groups[year] = []
  groups[year].push(article)
  return groups
}, {}))

const archiveTitle = computed(() => {
  if (selectedCategory.value) return categories.value.find(item => item.id === selectedCategory.value)?.name || '分类归档'
  if (selectedTag.value) return `标签：${tags.value.find(item => item.id === selectedTag.value)?.name || ''}`
  return '全部文章'
})

function formatDate(dateStr) {
  const date = new Date(dateStr)
  return `${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

function changePage(nextPage) {
  page.value = nextPage
  fetchArticles()
}

async function fetchArticles() {
  loading.value = true
  try {
    const params = { page: page.value, pageSize }
    if (selectedCategory.value) params.categoryId = selectedCategory.value
    if (selectedTag.value) params.tagId = selectedTag.value
    const res = await getArticles(params)
    articles.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch {
    articles.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    const [categoryRes, tagRes] = await Promise.all([getCategories(), getTags()])
    categories.value = categoryRes.data || []
    tags.value = tagRes.data || []
  } catch {
    categories.value = []
    tags.value = []
  }
})

watch(
  () => [route.query.categoryId, route.query.tagId],
  ([categoryId, tagId]) => {
    selectedCategory.value = categoryId ? Number(categoryId) : null
    selectedTag.value = tagId ? Number(tagId) : null
    page.value = 1
    fetchArticles()
  },
  { immediate: true }
)
</script>

<style>
.archive-card { padding: 32px 40px 38px; background: rgba(20, 28, 40, .92); }
.archive-heading { display: flex; justify-content: space-between; gap: 20px; align-items: center; margin: 0 0 52px; }
.archive-heading h1 { margin: 0; font-size: 22px; font-weight: 500; }
.archive-heading h1 span,
.archive-heading h1 i { color: var(--text-muted); font-style: normal; font-weight: 400; }
.archive-heading h1 i { margin: 0 14px; }
.archive-heading h1 strong { color: var(--accent); }
.archive-heading > span,
.page-info { color: var(--text-muted); font-size: 14px; white-space: nowrap; }

.archive-groups { display: flex; flex-direction: column; gap: 44px; }
.archive-group { display: grid; grid-template-columns: 135px minmax(0, 1fr); gap: 28px; }
.archive-year { margin: 0; color: var(--text-primary); font-size: 38px; line-height: 1; }
.archive-timeline { position: relative; padding-left: 68px; }
.archive-timeline::before { content: ''; position: absolute; top: 14px; bottom: 8px; left: 3px; border-left: 2px dashed rgba(255, 255, 255, .14); }
.archive-year-summary { position: relative; margin: 1px 0 22px; color: var(--text-secondary); font-size: 17px; }
.archive-year-summary::before { content: ''; position: absolute; left: -74px; top: 2px; width: 13px; height: 13px; border: 4px solid #ff8677; border-radius: 50%; background: var(--bg-card); }
.archive-item { position: relative; display: grid; grid-template-columns: 68px 18px minmax(0, 1fr); align-items: center; width: 100%; min-height: 48px; padding: 5px 0; color: var(--text-primary); text-align: left; text-decoration: none; cursor: pointer; }
.archive-date { color: var(--text-muted); font-size: 15px; }
.archive-dot { width: 6px; height: 6px; border-radius: 50%; background: rgba(255, 255, 255, .4); }
.archive-item-content { display: flex; align-items: center; gap: 12px; min-width: 0; }
.archive-category { flex: 0 0 auto; border-radius: 6px; padding: 4px 8px; color: var(--accent); background: var(--accent-dim); font-size: 13px; }
.archive-title { overflow: hidden; color: var(--text-primary); font-size: 17px; font-weight: 600; text-overflow: ellipsis; white-space: nowrap; transition: color .2s ease; }
.archive-item:hover .archive-title { color: var(--accent); }
.archive-item:hover .archive-dot { background: var(--accent); transform: scale(1.5); }

.pagination-wrap { display: flex; justify-content: center; align-items: center; gap: 16px; margin-top: 20px; }
.loading-state,.empty-state { padding: 60px 0; color: var(--text-muted); text-align: center; }

@media (max-width: 768px) {
  .archive-card { padding: 24px 18px; }
  .archive-heading { margin-bottom: 34px; }
  .archive-group { display: block; }
  .archive-year { margin-bottom: 26px; font-size: 32px; }
  .archive-timeline { padding-left: 54px; }
  .archive-year-summary::before { left: -60px; }
  .archive-item { grid-template-columns: 58px 14px minmax(0, 1fr); }
  .archive-item-content { gap: 7px; }
  .archive-category { display: none; }
  .archive-title { font-size: 15px; }
}
</style>
