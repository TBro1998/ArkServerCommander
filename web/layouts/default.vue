<template>
  <div class="min-h-screen bg-gray-50">
    <nav v-if="authStore.isAuthenticated" class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center space-x-8">
            <NuxtLink to="/" class="text-xl font-semibold text-gray-900 hover:text-blue-600 transition-colors">
              ARK服务器管理器
            </NuxtLink>
            
            <!-- 导航菜单 -->
            <nav class="flex space-x-6">
              <NuxtLink
                to="/"
                class="text-gray-600 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors"
                :class="{ 'text-blue-600 bg-blue-50': $route.path === '/' }"
              >
                控制台
              </NuxtLink>
              <NuxtLink
                to="/servers"
                class="text-gray-600 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors"
                :class="{ 'text-blue-600 bg-blue-50': $route.path === '/servers' }"
              >
                服务器管理
              </NuxtLink>
              <a
                href="http://localhost:8080/swagger/index.html"
                target="_blank"
                class="text-gray-600 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors"
              >
                API文档
                <svg class="inline w-3 h-3 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"></path>
                </svg>
              </a>
            </nav>
          </div>
          
          <div class="flex items-center space-x-4">
            <span class="text-sm text-gray-700">欢迎，{{ authStore.user?.username }}</span>
            <button
              @click="logout"
              class="bg-red-600 hover:bg-red-700 text-white px-3 py-2 rounded-md text-sm font-medium transition-colors"
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

// 在客户端初始化时从cookie恢复登录状态
onMounted(() => {
  authStore.initFromStorage()
})

const logout = () => {
  authStore.logout()
  navigateTo('/login')
}
</script> 