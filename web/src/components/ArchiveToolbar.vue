<template>
  <div class="archive-toolbar">
    <nav id="category-filter" class="archive-nav" aria-label="文章分类">
      <button class="archive-nav-item archive-home" type="button" title="首页" @click="router.push('/')">⌂</button>
      <div class="archive-nav-list">
        <button class="archive-nav-item" :class="{ active: !selectedCategory && !selectedTag }" type="button" @click="router.push('/archive')">归档 <span class="archive-nav-count">{{ siteStats.articleCount }}</span></button>
        <button v-for="category in categories" :key="category.id" :ref="element => setCategoryButton(element, category.id)" class="archive-nav-item" :class="{ active: selectedCategory === category.id }" type="button" @click="selectCategory(category.id)">{{ category.name }} <span class="archive-nav-count">{{ category.articleCount || 0 }}</span></button>
      </div>
      <button class="archive-nav-item archive-more" type="button" @click="router.push('/categories')">更多 <span aria-hidden="true">›</span></button>
    </nav>
  </div>
</template>

<script setup>
import { computed, nextTick, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const props = defineProps({
  categories: { type: Array, default: () => [] },
  siteStats: { type: Object, default: () => ({ articleCount: 0 }) }
})

const router = useRouter()
const route = useRoute()
const categoryButtons = ref(new Map())
// 根据当前响应式状态计算派生数据。
const selectedCategory = computed(() => route.query.categoryId ? Number(route.query.categoryId) : null)
// 根据当前响应式状态计算派生数据。
const selectedTag = computed(() => route.query.tagId ? Number(route.query.tagId) : null)

// 更新对应的状态值。
function setCategoryButton(element, categoryId) {
  if (element) {
    categoryButtons.value.set(categoryId, element)
  } else {
    categoryButtons.value.delete(categoryId)
  }
}

// 处理当前模块的相关逻辑。
async function scrollSelectedCategory() {
  await nextTick()
  const button = selectedCategory.value ? categoryButtons.value.get(selectedCategory.value) : null
  button?.scrollIntoView({ behavior: 'smooth', block: 'nearest', inline: 'center' })
}

// 选中指定的数据项。
async function selectCategory(categoryId) {
  await router.push({ path: '/archive', query: { categoryId, ...(selectedTag.value ? { tagId: selectedTag.value } : {}) } })
  scrollSelectedCategory()
}

watch(
  () => [selectedCategory.value, props.categories.length],
  scrollSelectedCategory,
  { immediate: true }
)

</script>

<style scoped>
.archive-toolbar { margin-top: 324px; margin-bottom: 20px; }

@media (max-width: 1024px) {
  .archive-toolbar { margin-top: 0; }
}
.archive-nav { display: flex; align-items: center; gap: 10px; padding: 14px 16px; overflow: hidden; background: rgba(20, 28, 40, .93); border: 1px solid var(--border); border-radius: 20px; backdrop-filter: blur(12px); }
.archive-nav { margin-bottom: 14px; }
.archive-nav-list { display: flex; align-items: center; gap: 10px; min-width: 0; overflow-x: auto; scrollbar-width: none; -ms-overflow-style: none; }
.archive-nav-list::-webkit-scrollbar { display: none; }
.archive-nav-item { flex: 0 0 auto; border: 0; border-radius: 999px; padding: 10px 17px; background: rgba(15,20,25,.68); color: var(--text-secondary); font-size: 15px; cursor: pointer; transition: background .2s ease,color .2s ease; }
.archive-home { min-width: 44px; padding-inline: 12px; color: var(--accent); font-size: 21px; }
.archive-more { margin-left: auto; }
.archive-more span { margin-left: 6px; font-size: 22px; line-height: 0; vertical-align: -2px; }
.archive-nav-item:hover { color: var(--text-primary); }
.archive-nav-item.active { color: var(--bg-primary); background: var(--accent); }
.archive-nav-count { display: inline-flex; align-items: center; justify-content: center; min-width: 1.25em; height: 1.25em; margin-left: 7px; padding: 0 4px; border-radius: 999px; background: rgba(255,255,255,.1); font-size: 11px; font-weight: 600; line-height: 1; opacity: .85; }
.archive-nav-item.active .archive-nav-count { background: rgba(15,20,25,.14); }
</style>
