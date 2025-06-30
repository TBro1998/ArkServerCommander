<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          ARK 服务器管理器
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          请登录您的账户
        </p>
      </div>
      
      <form class="mt-8 space-y-6" @submit.prevent="handleLogin">
        <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
          {{ error }}
        </div>
        
        <div class="rounded-md shadow-sm -space-y-px">
          <div>
            <label for="username" class="sr-only">用户名</label>
            <input
              id="username"
              v-model="form.username"
              name="username"
              type="text"
              required
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm"
              placeholder="用户名"
            />
          </div>
          <div>
            <label for="password" class="sr-only">密码</label>
            <input
              id="password"
              v-model="form.password"
              name="password"
              type="password"
              required
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm"
              placeholder="密码"
            />
          </div>
        </div>

        <div>
          <button
            type="submit"
            :disabled="authStore.isLoading"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
          >
            {{ authStore.isLoading ? '登录中...' : '登录' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: false
})

const authStore = useAuthStore()
const router = useRouter()

const form = reactive({
  username: '',
  password: ''
})

const error = ref('')

onMounted(async () => {
  // 检查是否已初始化
  const initialized = await authStore.checkInit()
  if (!initialized) {
    await navigateTo('/init')
  }
  
  // 检查是否已经登录
  authStore.initFromStorage()
  if (authStore.isAuthenticated) {
    await navigateTo('/')
  }
})

const handleLogin = async () => {
  error.value = ''
  
  const result = await authStore.login(form)
  
  if (result.success) {
    await navigateTo('/')
  } else {
    error.value = result.message
  }
}
</script> 