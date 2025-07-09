<template>
  <div class="min-h-screen flex flex-col bg-gradient-to-br from-blue-50 to-indigo-100">
    <div class="flex-1 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div class="w-full max-w-md">
        <!-- Logo 和标题区域 -->
        <div class="text-center mb-8">
          <div class="mx-auto w-16 h-16 bg-blue-600 rounded-full flex items-center justify-center mb-4">
            <UIcon name="i-lucide-server" class="w-8 h-8 text-white" />
          </div>
          <h1 class="text-3xl font-bold text-gray-900 mb-2">
            ARK 服务器管理器
          </h1>
          <p class="text-gray-600">
            安全登录您的管理账户
          </p>
        </div>
        
        <!-- 登录表单卡片 -->
        <UCard class="shadow-xl">
          <form @submit.prevent="handleLogin" class="space-y-6">
            <!-- 错误提示 -->
            <UAlert
              v-if="error"
              :title="error"
              color="red"
              variant="soft"
              icon="i-lucide-alert-circle"
            />
            
            <!-- 用户名输入 -->
            <UFormGroup label="用户名" name="username">
              <UInput
                v-model="form.username"
                placeholder="请输入用户名"
                icon="i-lucide-user"
                size="lg"
                required
              />
            </UFormGroup>
            
            <!-- 密码输入 -->
            <UFormGroup label="密码" name="password">
              <UInput
                v-model="form.password"
                type="password"
                placeholder="请输入密码"
                icon="i-lucide-lock"
                size="lg"
                required
              />
            </UFormGroup>

            <!-- 登录按钮 -->
            <UButton
              type="submit"
              :loading="authStore.isLoading"
              :disabled="authStore.isLoading"
              color="blue"
              variant="solid"
              size="lg"
              block
              class="mt-8"
            >
              <UIcon 
                v-if="!authStore.isLoading" 
                name="i-lucide-log-in" 
                class="w-4 h-4 mr-2" 
              />
              {{ authStore.isLoading ? '登录中...' : '登录' }}
            </UButton>
          </form>
        </UCard>
        
        <!-- 底部提示 -->
        <div class="text-center mt-6">
          <p class="text-sm text-gray-500">
            首次使用？系统将自动引导您完成初始化
          </p>
        </div>
      </div>
    </div>

    <!-- 登录页面底部商标信息 -->
    <AppFooter 
      :show-privacy-policy="false"
      :show-terms-of-service="false"
      page-info="安全登录系统"
    />
  </div>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth'

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