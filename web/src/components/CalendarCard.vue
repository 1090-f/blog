<template>
  <article class="calendar-card card" aria-label="日历">
    <header class="calendar-header">
      <button class="calendar-nav" type="button" aria-label="上个月" @click="changeMonth(-1)">‹</button>
      <h4>{{ displayedMonth.getFullYear() }}年{{ displayedMonth.getMonth() + 1 }}月</h4>
      <button class="calendar-nav" type="button" aria-label="下个月" @click="changeMonth(1)">›</button>
    </header>

    <div class="calendar-weekdays" aria-hidden="true">
      <span v-for="weekday in weekdays" :key="weekday">{{ weekday }}</span>
    </div>

    <div class="calendar-grid" role="grid">
      <span v-for="blank in firstWeekday" :key="`blank-${blank}`" class="calendar-day is-empty" aria-hidden="true"></span>
      <button
        v-for="day in daysInMonth"
        :key="day"
        class="calendar-day"
        :class="{ 'is-today': isToday(day), 'is-selected': isSelected(day) }"
        type="button"
        :aria-label="`${displayedMonth.getMonth() + 1}月${day}日`"
        @click="selectDay(day)"
      >
        {{ day }}
      </button>
    </div>

    <div class="calendar-activity" aria-hidden="true">
      <span v-for="(level, index) in activityLevels" :key="index" class="activity-cell" :class="`level-${level}`"></span>
    </div>
  </article>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { getSiteActivity } from '../api/site'

const ACTIVITY_REFRESH_INTERVAL = 30 * 60 * 1000
const now = new Date()
const displayedMonth = ref(new Date(now.getFullYear(), now.getMonth(), 1))
const selectedDate = ref(new Date(now.getFullYear(), now.getMonth(), now.getDate()))
const weekdays = ['日', '一', '二', '三', '四', '五', '六']

// 根据当前响应式状态计算派生数据。
const firstWeekday = computed(() => displayedMonth.value.getDay())
// 根据当前响应式状态计算派生数据。
const daysInMonth = computed(() => new Date(
  displayedMonth.value.getFullYear(),
  displayedMonth.value.getMonth() + 1,
  0
).getDate())

const activityByDay = ref({})
// 根据当前响应式状态计算派生数据。
const maxActivity = computed(() => Math.max(1, ...Object.values(activityByDay.value)))
// 根据当前响应式状态计算派生数据。
const activityLevels = computed(() => {
  const cellCount = Math.ceil((firstWeekday.value + daysInMonth.value) / 7) * 7
  return Array.from({ length: cellCount }, (_, index) => {
  const day = index - firstWeekday.value + 1
  if (day < 1 || day > daysInMonth.value) return 0

  const activity = activityByDay.value[day] || 0
  if (activity === 0) return 0
  const ratio = activity / maxActivity.value
  if (ratio >= .75) return 3
  if (ratio >= .4) return 2
  return 1
  })
})

// 更新当前筛选条件或页面状态。
function changeMonth(offset) {
  const nextMonth = new Date(
    displayedMonth.value.getFullYear(),
    displayedMonth.value.getMonth() + offset,
    1
  )
  displayedMonth.value = nextMonth
}

// 加载当前页面所需的数据。
async function loadActivity() {
  activityByDay.value = {}
  try {
    const response = await getSiteActivity({
      year: displayedMonth.value.getFullYear(),
      month: displayedMonth.value.getMonth() + 1
    })
    activityByDay.value = Object.fromEntries(
      (response.data?.days || []).map(item => [Number(item.date.slice(-2)), item.total || 0])
    )
  } catch {
    // The calendar remains usable when activity data is unavailable.
  }
}

// 选中指定的数据项。
function selectDay(day) {
  selectedDate.value = new Date(
    displayedMonth.value.getFullYear(),
    displayedMonth.value.getMonth(),
    day
  )
}

// 判断当前状态是否满足条件。
function isToday(day) {
  return displayedMonth.value.getFullYear() === now.getFullYear()
    && displayedMonth.value.getMonth() === now.getMonth()
    && day === now.getDate()
}

// 判断当前状态是否满足条件。
function isSelected(day) {
  return selectedDate.value.getFullYear() === displayedMonth.value.getFullYear()
    && selectedDate.value.getMonth() === displayedMonth.value.getMonth()
    && selectedDate.value.getDate() === day
}

watch(displayedMonth, loadActivity)

let activityRefreshTimer

onMounted(() => {
  loadActivity()
  activityRefreshTimer = window.setInterval(loadActivity, ACTIVITY_REFRESH_INTERVAL)
})

onUnmounted(() => {
  window.clearInterval(activityRefreshTimer)
})
</script>

<style scoped>
.calendar-card {
  padding: 22px;
}

.calendar-header {
  display: grid;
  grid-template-columns: 28px minmax(0, 1fr) 28px;
  align-items: center;
  gap: 8px;
  margin-bottom: 18px;
}

.calendar-header h4 {
  color: var(--text-primary);
  font-size: 21px;
  font-weight: 600;
  text-align: center;
}

.calendar-nav {
  width: 28px;
  height: 28px;
  border: 0;
  border-radius: 8px;
  background: transparent;
  color: var(--text-secondary);
  font-size: 28px;
  line-height: 1;
  cursor: pointer;
}

.calendar-nav:hover {
  background: var(--accent-dim);
  color: var(--accent);
}

.calendar-weekdays,
.calendar-grid {
  display: grid;
  grid-template-columns: repeat(7, minmax(0, 1fr));
  column-gap: 4px;
}

.calendar-weekdays {
  margin-bottom: 8px;
  color: var(--text-muted);
  font-size: 13px;
  text-align: center;
}

.calendar-day {
  position: relative;
  display: grid;
  place-items: center;
  min-width: 0;
  height: 34px;
  border: 1px solid transparent;
  border-radius: 7px;
  background: transparent;
  color: var(--text-secondary);
  font-size: 14px;
  cursor: pointer;
}

button.calendar-day:hover {
  background: var(--accent-dim);
  color: var(--text-primary);
}

.calendar-day.is-empty {
  cursor: default;
}

.calendar-day.is-today {
  border-color: var(--accent);
  color: var(--text-primary);
  font-weight: 600;
}

.calendar-day.is-selected {
  background: var(--accent);
  color: var(--bg-primary);
  font-weight: 600;
}

.calendar-day.is-today.is-selected {
  border-color: var(--accent);
}

.calendar-activity {
  display: grid;
  grid-template-columns: repeat(7, minmax(0, 1fr));
  gap: 4px;
  margin-top: 18px;
  padding-top: 14px;
  border-top: 1px solid var(--border);
}

.activity-cell {
  height: 20px;
  border-radius: 5px;
  background: rgba(45, 212, 191, .08);
}

.activity-cell.level-1 {
  background: rgba(45, 212, 191, .28);
}

.activity-cell.level-2 {
  background: rgba(45, 212, 191, .52);
}

.activity-cell.level-3 {
  background: var(--accent);
}

@media (max-width: 1200px) {
  .calendar-card {
    padding-inline: 16px;
  }

  .calendar-day {
    height: 32px;
    font-size: 13px;
  }

  .activity-cell {
    height: 17px;
  }
}
</style>
