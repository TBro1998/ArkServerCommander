<template>
  <div class="px-2 sm:px-4 py-4 sm:py-6">
    <div class="border-4 border-dashed border-gray-200 rounded-lg min-h-96 p-4 sm:p-8">
      <div class="text-center">
        <h1 class="text-2xl sm:text-3xl font-bold text-gray-900 mb-4">
          欢迎使用 ARK 服务器管理器
        </h1>
        <p class="text-base sm:text-lg text-gray-600 mb-8">
          您已成功登录系统，可以开始管理您的 ARK 服务器了。
        </p>
        
        <div class="bg-white p-4 sm:p-6 rounded-lg shadow-sm">
          <h2 class="text-lg sm:text-xl font-semibold text-gray-800 mb-4">系统信息</h2>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 text-left">
            <div>
              <span class="font-medium text-gray-700">用户名：</span>
              <span class="text-gray-900">{{ authStore.user?.username }}</span>
            </div>
            <div>
              <span class="font-medium text-gray-700">用户ID：</span>
              <span class="text-gray-900">{{ authStore.user?.id }}</span>
            </div>
          </div>
        </div>
        
        <div class="mt-8 space-y-4">
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 w-full max-w-6xl mx-auto">
            <!-- 服务器管理卡片 -->
            <NuxtLink
              to="/servers"
              class="bg-blue-50 hover:bg-blue-100 p-4 sm:p-6 rounded-lg border border-blue-200 transition-colors group"
            >
              <div class="text-center">
                <svg class="mx-auto h-12 w-12 text-blue-600 group-hover:text-blue-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
                </svg>
                <h3 class="mt-2 text-base sm:text-lg font-medium text-gray-900">服务器管理</h3>
                <p class="mt-1 text-sm text-gray-600">添加、配置和管理您的ARK服务器</p>
              </div>
            </NuxtLink>

            <!-- 占位卡片：玩家管理 -->
            <div class="bg-gray-50 p-4 sm:p-6 rounded-lg border border-gray-200 opacity-60">
              <div class="text-center">
                <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"></path>
                </svg>
                <h3 class="mt-2 text-base sm:text-lg font-medium text-gray-500">玩家管理</h3>
                <p class="mt-1 text-sm text-gray-400">即将推出...</p>
              </div>
            </div>

            <!-- 占位卡片：日志监控 -->
            <div class="bg-gray-50 p-4 sm:p-6 rounded-lg border border-gray-200 opacity-60">
              <div class="text-center">
                <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
                </svg>
                <h3 class="mt-2 text-base sm:text-lg font-medium text-gray-500">日志监控</h3>
                <p class="mt-1 text-sm text-gray-400">即将推出...</p>
              </div>
            </div>
          </div>
          
          <p class="text-sm text-gray-500 mt-6">
            点击上方卡片开始管理您的ARK服务器
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  middleware: 'auth'
})

const authStore = useAuthStore()

onMounted(() => {
  authStore.initFromStorage()
})
</script> 