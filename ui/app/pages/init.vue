<template>
  <div class="min-h-screen flex flex-col bg-gradient-to-br from-green-50 to-emerald-100">
    <div class="flex-1 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div class="w-full max-w-md">
        <!-- Logo 和标题区域 -->
        <div class="text-center mb-8">
          <div class="mx-auto w-16 h-16 bg-green-600 rounded-full flex items-center justify-center mb-4">
            <UIcon name="i-lucide-settings" class="w-8 h-8 text-white" />
          </div>
          <h1 class="text-3xl font-bold text-gray-900 mb-2">
            系统初始化
          </h1>
          <p class="text-gray-600">
            首次使用，请设置管理员账户
          </p>
        </div>
        
        <!-- 初始化表单卡片 -->
        <UCard class="shadow-xl">
          <form @submit.prevent="handleInit" class="space-y-6">
            <!-- 错误提示 -->
            <UAlert
              v-if="error"
              :title="error"
              color="red"
              variant="soft"
              icon="i-lucide-alert-circle"
            />
            
            <!-- 成功提示 -->
            <UAlert
              v-if="success"
              :title="success"
              color="green"
              variant="soft"
              icon="i-lucide-check-circle"
            />
            
            <!-- 用户名输入 -->
            <UFormGroup label="管理员用户名" name="username">
              <UInput
                v-model="form.username"
                placeholder="请输入管理员用户名"
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
                placeholder="请输入密码 (至少6位)"
                icon="i-lucide-lock"
                size="lg"
                required
                minlength="6"
              />
              <template #help>
                <p class="text-xs text-gray-500 mt-1">
                  密码至少需要6位字符
                </p>
              </template>
            </UFormGroup>
            
            <!-- 确认密码输入 -->
            <UFormGroup label="确认密码" name="confirmPassword">
              <UInput
                v-model="form.confirmPassword"
                type="password"
                placeholder="请再次输入密码"
                icon="i-lucide-lock"
                size="lg"
                required
              />
            </UFormGroup>

            <!-- 初始化按钮 -->
            <UButton
              type="submit"
              :loading="authStore.isLoading"
              :disabled="authStore.isLoading"
              color="green"
              variant="solid"
              size="lg"
              block
              class="mt-8"
            >
              <UIcon 
                v-if="!authStore.isLoading" 
                name="i-lucide-settings" 
                class="w-4 h-4 mr-2" 
              />
              {{ authStore.isLoading ? '初始化中...' : '初始化系统' }}
            </UButton>
          </form>
        </UCard>
        
        <!-- 底部提示 -->
        <div class="text-center mt-6">
          <div class="flex items-center justify-center space-x-2 text-sm text-gray-500">
            <UIcon name="i-lucide-info" class="w-4 h-4" />
            <span>初始化完成后将自动跳转到主页面</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 初始化页面底部商标信息 -->
    <AppFooter 
      :show-privacy-policy="false"
      :show-terms-of-service="false"
      page-info="系统初始化向导"
    />
  </div>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth'

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