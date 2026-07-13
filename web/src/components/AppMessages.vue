<template>
  <teleport to="body">
    <div class="message-stack">
      <transition-group name="message">
        <div
          v-for="item in state.items"
          :key="item.id"
          :class="['message-item', `message-${item.type}`]"
        >
          {{ item.text }}
        </div>
      </transition-group>
    </div>
  </teleport>
</template>

<script setup>
import { state } from '../utils/message'
</script>

<style scoped>
.message-stack {
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 2000;
  display: flex;
  flex-direction: column;
  gap: 10px;
  pointer-events: none;
}

.message-item {
  min-width: 220px;
  max-width: min(90vw, 420px);
  padding: 12px 16px;
  border-radius: 12px;
  box-shadow: 0 16px 40px rgba(0, 0, 0, 0.28);
  backdrop-filter: blur(14px);
  color: #f8fafc;
  font-size: 14px;
  line-height: 1.4;
  text-align: center;
}

.message-success {
  background: rgba(16, 185, 129, 0.9);
}

.message-error {
  background: rgba(239, 68, 68, 0.92);
}

.message-warning {
  background: rgba(245, 158, 11, 0.92);
}

.message-info {
  background: rgba(59, 130, 246, 0.92);
}

.message-enter-active,
.message-leave-active {
  transition: all 0.22s ease;
}

.message-enter-from,
.message-leave-to {
  opacity: 0;
  transform: translateY(-8px) scale(0.96);
}
</style>
