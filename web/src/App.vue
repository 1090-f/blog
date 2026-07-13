<template>
  <div class="dark">
    <router-view />
    <AppMessages />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import AppMessages from './components/AppMessages.vue'
import { getProfile } from './api/auth'
import { useUserStore } from './stores/user'

const userStore = useUserStore()

onMounted(async () => {
  if (!userStore.token) {
    return
  }

  try {
    const res = await getProfile()
    userStore.setAuth(userStore.token, res.data)
  } catch (error) {
    // Expired tokens are handled by the request interceptor.
  }
})
</script>
