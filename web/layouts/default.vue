<template>
  <div class="min-h-screen bg-gray-50">
    <nav v-if="authStore.isAuthenticated" class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <h1 class="text-xl font-semibold text-gray-900">ARK服务器管理器</h1>
          </div>
          <div class="flex items-center space-x-4">
            <span class="text-sm text-gray-700">欢迎，{{ authStore.user?.username }}</span>
            <button
              @click="logout"
              class="bg-red-600 hover:bg-red-700 text-white px-3 py-2 rounded-md text-sm font-medium"
            >
              退出登录
            </button>
          </div>
        </div>
      </div>
    </nav>
    
    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <slot />
    </main>
  </div>
</template>

<script setup>
const authStore = useAuthStore()

const logout = () => {
  authStore.logout()
  navigateTo('/login')
}
</script> 