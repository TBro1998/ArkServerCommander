<template>
  <UApp>
    <NuxtPage />
  </UApp>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth'

// 全局初始化auth store
const authStore = useAuthStore()

// 在应用启动时从cookie恢复登录状态
onMounted(() => {
  authStore.initFromStorage()
})
</script>