import { reactive } from 'vue'

let nextId = 1

export const state = reactive({
  items: []
})

function show(type, text, duration = 2400) {
  const id = nextId++
  state.items.push({ id, type, text })

  window.setTimeout(() => {
    const index = state.items.findIndex(item => item.id === id)
    if (index >= 0) {
      state.items.splice(index, 1)
    }
  }, duration)
}

export const message = {
  success(text, duration) {
    show('success', text, duration)
  },
  error(text, duration) {
    show('error', text, duration)
  },
  warning(text, duration) {
    show('warning', text, duration)
  },
  info(text, duration) {
    show('info', text, duration)
  }
}
