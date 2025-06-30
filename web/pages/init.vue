<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          系统初始化
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          首次使用，请设置管理员账户
        </p>
      </div>
      
      <form class="mt-8 space-y-6" @submit.prevent="handleInit">
        <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
          {{ error }}
        </div>
        
        <div v-if="success" class="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded">
          {{ success }}
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
              placeholder="管理员用户名"
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
              minlength="6"
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm"
              placeholder="密码 (至少6位)"
            />
          </div>
          <div>
            <label for="confirmPassword" class="sr-only">确认密码</label>
            <input
              id="confirmPassword"
              v-model="form.confirmPassword"
              name="confirmPassword"
              type="password"
              required
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm"
              placeholder="确认密码"
            />
          </div>
        </div>

        <div>
          <button
            type="submit"
            :disabled="authStore.isLoading"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50"
          >
            {{ authStore.isLoading ? '初始化中...' : '初始化系统' }}
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

const form = reactive({
  username: '',
  password: '',
  confirmPassword: ''
})

const error = ref('')
const success = ref('')

onMounted(async () => {
  // 检查是否已初始化
  const initialized = await authStore.checkInit()
  if (initialized) {
    await navigateTo('/login')
  }
})

const handleInit = async () => {
  error.value = ''
  success.value = ''
  
  // 验证密码
  if (form.password !== form.confirmPassword) {
    error.value = '两次输入的密码不一致'
    return
  }
  
  if (form.password.length < 6) {
    error.value = '密码至少需要6位'
    return
  }
  
  const result = await authStore.init({
    username: form.username,
    password: form.password
  })
  
  if (result.success) {
    success.value = result.message
    setTimeout(() => {
      navigateTo('/')
    }, 1500)
  } else {
    error.value = result.message
  }
}
</script> 