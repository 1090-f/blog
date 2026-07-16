<template>
  <div class="taxonomy-page categories-page">
    <section class="taxonomy-hero card">
      <h1>分类</h1>
      <p>全部分类 · {{ totalArticles }} 篇文章</p>
    </section>

    <div v-if="loading" class="loading-state">分类加载中...</div>
    <div v-else-if="categories.length" class="category-grid">
      <router-link
        v-for="category in categories"
        :key="category.id"
        :to="{ path: '/archive', query: { categoryId: category.id } }"
        class="category-card card"
      >
        <span class="category-folder" aria-hidden="true">
          <svg viewBox="0 0 48 48">
            <path d="M5 13a4 4 0 0 1 4-4h11l5 5h14a4 4 0 0 1 4 4v17a4 4 0 0 1-4 4H9a4 4 0 0 1-4-4V13Z" />
            <path d="M5 17h38" />
          </svg>
        </span>
        <span class="category-info">
          <strong>{{ category.name }}</strong>
          <small>{{ category.articleCount }} 篇文章</small>
        </span>
        <span class="category-arrow" aria-hidden="true">›</span>
      </router-link>
    </div>
    <div v-else class="empty-state">暂无分类</div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { getArticles } from '../../api/article'
import { getCategories } from '../../api/category'
import { getSiteStats } from '../../api/site'

const categories = ref([])
const siteStats = ref({ articleCount: 0 })
const loading = ref(true)
// 根据当前响应式状态计算派生数据。
const totalArticles = computed(() => siteStats.value.articleCount || categories.value.reduce((sum, category) => sum + category.articleCount, 0))

// 加载当前页面所需的数据。
async function loadCategories() {
  try {
    const [categoryResponse, statsResponse] = await Promise.all([getCategories(), getSiteStats()])
    const rawCategories = categoryResponse.data || []
    siteStats.value = statsResponse.data || siteStats.value
    const countedCategories = await Promise.all(rawCategories.map(async category => {
      try {
        const response = await getArticles({ categoryId: category.id, page: 1, pageSize: 1 })
        return { ...category, articleCount: response.data?.total || 0 }
      } catch {
        return { ...category, articleCount: 0 }
      }
    }))
    categories.value = countedCategories
  } finally {
    loading.value = false
  }
}

onMounted(loadCategories)
</script>

<style scoped>
.taxonomy-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.taxonomy-hero {
  padding: 38px 46px;
}

.taxonomy-hero h1 {
  margin-bottom: 8px;
  color: var(--accent);
  font-size: 32px;
  line-height: 1.2;
}

.taxonomy-hero p {
  color: var(--text-muted);
  font-size: 15px;
}

.category-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 24px;
}

.category-card {
  display: flex;
  align-items: center;
  gap: 22px;
  min-height: 158px;
  padding: 28px 34px;
  color: var(--text-primary);
  text-decoration: none;
}

.category-folder {
  display: grid;
  place-items: center;
  width: 68px;
  height: 68px;
  flex: 0 0 auto;
  border-radius: 50%;
  color: var(--accent);
  background: var(--accent-dim);
}

.category-folder svg {
  width: 42px;
  height: 42px;
  fill: currentColor;
  stroke: currentColor;
  stroke-width: 2;
}

.category-folder svg path:last-child {
  fill: none;
}

.category-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
  gap: 5px;
}

.category-info strong {
  overflow: hidden;
  font-size: 20px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.category-info small {
  color: var(--text-muted);
  font-size: 14px;
}

.category-arrow {
  margin-left: auto;
  color: var(--text-muted);
  font-size: 30px;
  line-height: 1;
}

.category-card:hover .category-arrow {
  color: var(--accent);
}

.loading-state,
.empty-state {
  padding: 60px 0;
  color: var(--text-muted);
  text-align: center;
}

@media (max-width: 900px) {
  .category-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 600px) {
  .taxonomy-hero {
    padding: 28px 22px;
  }

  .taxonomy-hero h1 {
    font-size: 27px;
  }

  .category-grid {
    grid-template-columns: 1fr;
  }
}
</style>
