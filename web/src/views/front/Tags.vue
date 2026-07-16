<template>
  <div class="taxonomy-page tags-page">
    <section class="taxonomy-hero card">
      <h1>标签</h1>
      <p>全部标签 · {{ tags.length }} 个标签</p>
    </section>

    <section class="tag-cloud-card card">
      <div v-if="loading" class="loading-state">标签加载中...</div>
      <div v-else-if="tags.length" class="tag-cloud">
        <router-link
          v-for="tag in tags"
          :key="tag.id"
          :to="{ path: '/archive', query: { tagId: tag.id } }"
          class="tag-pill"
          :class="{ 'is-popular': tag.articleCount >= popularThreshold }"
        >
          <span>#{{ tag.name }}</span>
          <strong>{{ tag.articleCount }}</strong>
        </router-link>
      </div>
      <div v-else class="empty-state">暂无标签</div>
    </section>

    <section v-if="rankedTags.length" class="tag-ranking card">
      <h2>Top 10</h2>
      <div v-for="(tag, index) in rankedTags" :key="tag.id" class="ranking-item">
        <span class="ranking-index">{{ index + 1 }}</span>
        <div class="ranking-main">
          <div class="ranking-label">
            <span>#{{ tag.name }}</span>
            <small>{{ tag.articleCount }} 篇文章</small>
          </div>
          <div class="ranking-track">
            <span :style="{ width: `${(tag.articleCount / maxCount) * 100}%` }"></span>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { getArticles } from '../../api/article'
import { getTags } from '../../api/tag'

const tags = ref([])
const loading = ref(true)
// 根据当前响应式状态计算派生数据。
const rankedTags = computed(() => [...tags.value].filter(tag => tag.articleCount > 0).sort((a, b) => b.articleCount - a.articleCount).slice(0, 10))
// 根据当前响应式状态计算派生数据。
const maxCount = computed(() => rankedTags.value[0]?.articleCount || 1)
// 根据当前响应式状态计算派生数据。
const popularThreshold = computed(() => Math.max(2, maxCount.value * .6))

// 加载当前页面所需的数据。
async function loadTags() {
  try {
    const response = await getTags()
    const rawTags = response.data || []
    tags.value = await Promise.all(rawTags.map(async tag => {
      try {
        const articleResponse = await getArticles({ tagId: tag.id, page: 1, pageSize: 1 })
        return { ...tag, articleCount: articleResponse.data?.total || 0 }
      } catch {
        return { ...tag, articleCount: 0 }
      }
    }))
  } finally {
    loading.value = false
  }
}

onMounted(loadTags)
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

.tag-cloud-card {
  padding: 34px 42px;
}

.tag-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 14px;
}

.tag-pill {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 9px 14px;
  border-radius: 10px;
  color: var(--accent);
  background: var(--accent-dim);
  font-size: 14px;
  text-decoration: none;
}

.tag-pill strong {
  min-width: 25px;
  padding: 1px 7px;
  border-radius: 9px;
  color: var(--bg-primary);
  background: var(--accent);
  font-size: 12px;
  text-align: center;
}

.tag-pill.is-popular {
  background: rgba(45, 212, 191, .24);
}

.tag-ranking {
  padding: 36px 42px;
}

.tag-ranking h2 {
  margin-bottom: 30px;
  font-size: 22px;
}

.ranking-item {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 28px;
}

.ranking-item:last-child {
  margin-bottom: 0;
}

.ranking-index {
  width: 22px;
  padding-top: 2px;
  color: var(--accent);
  font-size: 18px;
  font-weight: 700;
  text-align: center;
}

.ranking-main {
  flex: 1;
  min-width: 0;
}

.ranking-label {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 8px;
  color: var(--text-secondary);
  font-size: 14px;
}

.ranking-label small {
  color: var(--accent);
  font-size: 14px;
  white-space: nowrap;
}

.ranking-track {
  height: 9px;
  overflow: hidden;
  border-radius: 999px;
  background: var(--accent-dim);
}

.ranking-track span {
  display: block;
  height: 100%;
  border-radius: inherit;
  background: var(--accent);
}

.loading-state,
.empty-state {
  padding: 60px 0;
  color: var(--text-muted);
  text-align: center;
}

@media (max-width: 600px) {
  .taxonomy-hero,
  .tag-cloud-card,
  .tag-ranking {
    padding: 28px 22px;
  }

  .taxonomy-hero h1 {
    font-size: 27px;
  }

  .ranking-label {
    align-items: flex-start;
    flex-direction: column;
    gap: 3px;
  }
}
</style>
