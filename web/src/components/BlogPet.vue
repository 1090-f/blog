<template>
  <aside
    class="blog-pet"
    :class="{ 'is-dragging': dragging }"
    :style="petStyle"
    aria-label="博客小宠物"
  >
    <div class="pet-bubble" role="status">{{ message }}</div>

    <button
      class="pet-character"
      type="button"
      :aria-label="`和${petName}说话`"
      @click="talk"
      @pointerdown="startDrag"
      @mouseenter="startPetTalk"
      @mouseleave="stopPetTalk"
    >
      <span class="pet-spark pet-spark-one">✦</span>
      <span class="pet-spark pet-spark-two">✧</span>
      <span class="pet-portrait" :style="petImageStyle">
        <img class="pet-image" src="/pets/shorekeeper-chibi.png" alt="守岸人风格的小宠物" draggable="false" />
        <span class="pet-eye-glow pet-eye-glow-left" :style="eyeStyle('left')" aria-hidden="true" />
        <span class="pet-eye-glow pet-eye-glow-right" :style="eyeStyle('right')" aria-hidden="true" />
      </span>
      <!--
          <path class="pet-tail" d="M130 183c30 17 34-19 17-27-8-4-15 2-12 10 2 5 8 4 10 1" />
          <path class="pet-ear" d="M57 64 63 25c2-10 12-8 24 5l4 17Z" />
          <path class="pet-ear" d="m123 64-6-39c-2-10-12-8-24 5l-4 17Z" />
          <path class="pet-ear-inner" d="m65 57 4-21c5 3 9 7 12 12Z" />
          <path class="pet-ear-inner" d="m115 57-4-21c-5 3-9 7-12 12Z" />
          <path class="pet-star" d="m90 28 3 7 7 3-7 3-3 7-3-7-7-3 7-3Z" />
          <ellipse class="pet-face" cx="90" cy="95" rx="39" ry="42" />
          <path class="pet-hair" d="M53 79c5-30 18-42 37-42 20 0 33 12 37 42-10-7-16-13-20-22-5 8-10 12-17 14-7-2-12-6-17-14-4 9-10 15-20 22Z" />
          <g :transform="eyeTransform('left')">
            <ellipse class="pet-eye" cx="75" cy="104" rx="6" ry="9" />
            <circle class="pet-eye-light" cx="77" cy="101" r="2" />
          </g>
          <g :transform="eyeTransform('right')">
            <ellipse class="pet-eye" cx="105" cy="104" rx="6" ry="9" />
            <circle class="pet-eye-light" cx="107" cy="101" r="2" />
          </g>
          <circle class="pet-cheek" cx="62" cy="119" r="5" />
          <circle class="pet-cheek" cx="118" cy="119" r="5" />
          <path class="pet-mouth" d="M84 119c4 5 8 5 12 0" />
          <path class="pet-body" d="M53 153c11-12 63-12 74 0l10 65H43Z" />
          <path class="pet-collar" d="M65 151h50l-25 29Z" />
          <path class="pet-tie" d="m86 174 8 0-4 16Z" />
          <path class="pet-arm" d="M55 168c-8 15-8 30-2 40" />
          <path class="pet-arm" d="M125 168c8 15 8 30 2 40" />
          <path class="pet-foot" d="M54 216c-6 5-10 5-15 3M126 216c6 5 10 5 15 3" />
        </svg> -->
    </button>
  </aside>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

const messages = [
  '今天也要记得休息一下喔～',
  '要去看看「留言板」么？',
  '欢迎来到我的小角落！',
  '写完这篇文章就去喝水吧 ✦',
]
const petHoverMessages = ['被你发现啦～', '要和我一起看看文章吗？', '今天也要开心地写博客哦！', '再陪我一会儿嘛 ✦']
const petName = '小蝙蝠'
const messageIndex = ref(0)
const dragging = ref(false)
const position = ref({ x: 24, y: 28 })
const dragOffset = ref({ x: 0, y: 0 })
const mousePosition = ref({ x: 0.5, y: 0.5 })
const hoverMessage = ref('')
let petTalkTimer = null

// 根据当前响应式状态计算派生数据。
const message = computed(() => hoverMessage.value || messages[messageIndex.value])
// 根据当前响应式状态计算派生数据。
const petStyle = computed(() => ({ '--pet-x': `${position.value.x}px`, '--pet-y': `${position.value.y}px` }))
// 根据当前响应式状态计算派生数据。
const petImageStyle = computed(() => {
  const x = (mousePosition.value.x - 0.5) * 8
  const y = (mousePosition.value.y - 0.5) * 4
  const rotate = (mousePosition.value.x - 0.5) * 5
  return { transform: `translate(${x}px, ${y}px) rotate(${rotate}deg)` }
})

// 计算指定眼睛的动态样式。
function eyeStyle(side) {
  const horizontal = (mousePosition.value.x - 0.5) * 4
  const vertical = (mousePosition.value.y - 0.5) * 2.5
  const sideOffset = side === 'left' ? -0.3 : 0.3
  return { transform: `translate(${horizontal + sideOffset}px, ${vertical}px)` }
}

// 显示桌宠的随机对话。
function talk() {
  messageIndex.value = (messageIndex.value + 1) % messages.length
  hoverMessage.value = messages[messageIndex.value]
}

// 启动对应的交互或定时任务。
function startPetTalk() {
  let index = 0
  hoverMessage.value = petHoverMessages[index]
  clearInterval(petTalkTimer)
  petTalkTimer = setInterval(() => {
    index = (index + 1) % petHoverMessages.length
    hoverMessage.value = petHoverMessages[index]
  }, 1800)
}

// 停止对应的交互或定时任务。
function stopPetTalk() {
  clearInterval(petTalkTimer)
  petTalkTimer = null
  hoverMessage.value = ''
}

// 启动对应的交互或定时任务。
function startDrag(event) {
  if (event.button !== 0) return
  dragging.value = true
  const rect = event.currentTarget.closest('.blog-pet').getBoundingClientRect()
  dragOffset.value = { x: event.clientX - rect.left, y: event.clientY - rect.top }
  window.addEventListener('pointermove', onDrag)
  window.addEventListener('pointerup', stopDrag, { once: true })
}

// 处理当前模块的相关逻辑。
function onDrag(event) {
  position.value = {
    x: Math.max(8, Math.min(window.innerWidth - 155, event.clientX - dragOffset.value.x)),
    y: Math.max(8, Math.min(window.innerHeight - 205, window.innerHeight - event.clientY + dragOffset.value.y - 225)),
  }
}

// 停止对应的交互或定时任务。
function stopDrag() {
  dragging.value = false
  window.removeEventListener('pointermove', onDrag)
  localStorage.setItem('gin-blog-pet-position', JSON.stringify(position.value))
}

// 处理当前模块的相关逻辑。
function followMouse(event) {
  mousePosition.value = {
    x: event.clientX / Math.max(window.innerWidth, 1),
    y: event.clientY / Math.max(window.innerHeight, 1),
  }
}

// 查找应显示提示的目标元素。
function getHintTarget(target) {
  if (!(target instanceof Element) || target.closest('.blog-pet')) return null
  return target.closest('[data-pet-hint], a, button, h1, h2, h3, h4, h5, h6, p, li, label, time, .tag, .article-category, .article-title')
}

// 读取目标元素的提示文本。
function getHintText(element) {
  const text = element.dataset.petHint || element.textContent.replace(/\s+/g, ' ').trim()
  return text.length > 52 ? `${text.slice(0, 52)}…` : text
}

// 处理用户操作或浏览器事件。
function handleHintEnter(event) {
  const hintTarget = getHintTarget(event.target)
  if (!hintTarget || hintTarget.contains(event.relatedTarget)) return
  const text = getHintText(hintTarget)
  if (text) hoverMessage.value = text
}

// 处理用户操作或浏览器事件。
function handleHintLeave(event) {
  const hintTarget = getHintTarget(event.target)
  if (!hintTarget || hintTarget.contains(event.relatedTarget)) return
  hoverMessage.value = ''
}

onMounted(() => {
  window.addEventListener('mousemove', followMouse)
  document.addEventListener('mouseover', handleHintEnter)
  document.addEventListener('mouseout', handleHintLeave)
  try {
    const saved = JSON.parse(localStorage.getItem('gin-blog-pet-position'))
    if (saved && Number.isFinite(saved.x) && Number.isFinite(saved.y)) {
      position.value = saved
    }
  } catch {
    // Ignore malformed local storage from an older version.
  }
})

onBeforeUnmount(() => {
  stopDrag()
  stopPetTalk()
  window.removeEventListener('mousemove', followMouse)
  document.removeEventListener('mouseover', handleHintEnter)
  document.removeEventListener('mouseout', handleHintLeave)
})
</script>

<style scoped>
.blog-pet { position: fixed; left: var(--pet-x); bottom: var(--pet-y); width: 155px; height: 225px; z-index: 120; user-select: none; transition: opacity .2s; }
.pet-bubble { position: absolute; left: 0; bottom: 185px; width: 160px; min-height: 52px; padding: 11px 13px; color: #55475c; font-size: 12px; line-height: 1.5; background: #fffdfc; border: 1px solid #ece5e6; border-radius: 16px 16px 16px 4px; box-shadow: 0 8px 24px rgba(52, 33, 53, .14); }
.pet-bubble::after { content: ''; position: absolute; left: 24px; bottom: -9px; width: 16px; height: 16px; background: #fffdfc; border-right: 1px solid #ece5e6; border-bottom: 1px solid #ece5e6; transform: rotate(45deg); }
.pet-character { position: absolute; left: 10px; bottom: 0; width: 135px; height: 190px; padding: 0; border: 0; background: transparent; cursor: grab; touch-action: none; }
.is-dragging .pet-character { cursor: grabbing; }
.pet-portrait { position: absolute; left: -27px; bottom: -2px; width: 190px; height: 190px; transform-origin: 50% 24%; transition: transform .16s cubic-bezier(.2, .8, .2, 1); pointer-events: none; }
.pet-image { width: 190px; height: 190px; object-fit: contain; filter: drop-shadow(0 10px 8px rgba(19, 39, 71, .28)); }
.pet-eye-glow { position: absolute; width: 4px; height: 4px; border-radius: 50%; background: #d9ffff; box-shadow: 0 0 5px 2px rgba(92, 229, 239, .9); pointer-events: none; transition: transform .08s linear; }
.pet-eye-glow-left { left: 75px; top: 70px; }.pet-eye-glow-right { left: 94px; top: 70px; }
.pet-svg { width: 100%; height: 100%; filter: drop-shadow(0 8px 5px rgba(40, 28, 58, .18)); }
.pet-tail { fill: none; stroke: #83cbd7; stroke-width: 13; stroke-linecap: round; stroke-linejoin: round; }
.pet-ear { fill: #a8dce3; stroke: #314762; stroke-width: 3; stroke-linejoin: round; }
.pet-ear-inner { fill: #d4f2f1; }
.pet-hair { fill: #d7eef0; stroke: #314762; stroke-width: 3; stroke-linejoin: round; }
.pet-face { fill: #f5fbf7; stroke: #314762; stroke-width: 3; }
.pet-eye { fill: #244d72; }
.pet-eye-light { fill: #fff; }
.pet-cheek { fill: #9edce2; opacity: .8; }
.pet-mouth { fill: none; stroke: #527d9c; stroke-width: 2.5; stroke-linecap: round; }
.pet-body { fill: #243d63; stroke: #314762; stroke-width: 3; }
.pet-collar, .pet-tie { fill: #83cbd7; stroke: #314762; stroke-width: 2.5; }
.pet-arm, .pet-foot { fill: none; stroke: #314762; stroke-width: 7; stroke-linecap: round; }
.pet-star { fill: #d8f7f1; stroke: #5fb8c9; stroke-width: 1.5; }
.pet-spark { position: absolute; color: #a7e8ea; font-size: 18px; animation: floaty 2.4s ease-in-out infinite; text-shadow: 0 0 8px rgba(130, 224, 232, .8); }
.pet-spark-one { left: 4px; top: 56px; }.pet-spark-two { right: 1px; top: 94px; animation-delay: .8s; }
@keyframes floaty { 0%, 100% { transform: translateY(0) rotate(-8deg); } 50% { transform: translateY(-6px) rotate(8deg); } }
@media (max-width: 600px) { .blog-pet { transform: scale(.82); transform-origin: bottom left; } .pet-bubble { bottom: 185px; } }
@media (prefers-reduced-motion: reduce) { .pet-spark { animation: none; } .pet-portrait { transition: none; } }
</style>
