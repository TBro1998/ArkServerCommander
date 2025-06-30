<template>
  <div class="px-4 py-6 sm:px-0">
    <div class="border-4 border-dashed border-gray-200 rounded-lg min-h-96 p-8">
      <div class="text-center">
        <h1 class="text-3xl font-bold text-gray-900 mb-4">
          欢迎使用 ARK 服务器管理器
        </h1>
        <p class="text-lg text-gray-600 mb-8">
          您已成功登录系统，可以开始管理您的 ARK 服务器了。
        </p>
        
        <div class="bg-white p-6 rounded-lg shadow-sm">
          <h2 class="text-xl font-semibold text-gray-800 mb-4">系统信息</h2>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-left">
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
        
        <div class="mt-8">
          <p class="text-sm text-gray-500">
            更多功能正在开发中...
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