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
            {{ $t('auth.initTitle') }}
          </h1>
          <p class="text-gray-600">
            {{ $t('auth.initSubtitle') }}
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
            <UFormGroup :label="$t('auth.adminUsername')" name="username">
              <UInput
                v-model="form.username"
                :placeholder="$t('auth.enterAdminUsername')"
                icon="i-lucide-user"
                size="lg"
                required
              />
            </UFormGroup>
            
            <!-- 密码输入 -->
            <UFormGroup :label="$t('auth.password')" name="password">
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
                  {{ $t('auth.passwordMinLength') }}
                </p>
              </template>
            </UFormGroup>
            
            <!-- 确认密码输入 -->
            <UFormGroup :label="$t('auth.confirmPassword')" name="confirmPassword">
              <UInput
                v-model="form.confirmPassword"
                type="password"
                :placeholder="$t('auth.enterConfirmPassword')"
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
              {{ authStore.isLoading ? $t('auth.initLoading') : $t('auth.initButton') }}
            </UButton>
          </form>
        </UCard>
        
        <!-- 底部提示 -->
        <div class="text-center mt-6">
          <div class="flex items-center justify-center space-x-2 text-sm text-gray-500">
            <UIcon name="i-lucide-info" class="w-4 h-4" />
            <span>{{ $t('auth.initTip') }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 初始化页面底部商标信息 -->
    <AppFooter 
      :show-privacy-policy="false"
      :show-terms-of-service="false"
      :page-info="$t('auth.initWizard')"
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
    error.value = $t('auth.passwordMismatch')
    return
  }
  
  if (form.password.length < 6) {
    error.value = $t('auth.passwordMinLengthError')
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