<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <nav v-if="authStore.isAuthenticated" class="bg-white shadow-sm border-b">
      <div class="w-full px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center space-x-8">
            <NuxtLink to="/" class="text-xl font-semibold text-gray-900 hover:text-blue-600 transition-colors">
              {{ $t('auth.loginTitle') }}
            </NuxtLink>
            
            <!-- 桌面端导航菜单 -->
            <nav class="hidden md:flex space-x-6">
              <NuxtLink
                to="/"
                class="text-gray-600 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors"
                :class="{ 'text-blue-600 bg-blue-50': $route.path === '/' }"
              >
                {{ $t('navigation.dashboard') }}
              </NuxtLink>
              <NuxtLink
                to="/servers"
                class="text-gray-600 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors"
                :class="{ 'text-blue-600 bg-blue-50': $route.path === '/servers' }"
              >
                {{ $t('navigation.servers') }}
              </NuxtLink>
            </nav>
          </div>
          
          <div class="flex items-center space-x-4">
            <span class="hidden sm:block text-sm text-gray-700">{{ $t('navigation.welcome') }}，{{ authStore.user?.username }}</span>
            
            <!-- 语言切换器 -->
            <LanguageSwitcher />
            
            <!-- 移动端菜单按钮 -->
            <button
              @click="mobileMenuOpen = !mobileMenuOpen"
              class="md:hidden p-2 rounded-md text-gray-600 hover:text-gray-900 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path v-if="!mobileMenuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path>
                <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
            </button>
            <button
              @click="logout"
              class="bg-red-600 hover:bg-red-700 text-white px-3 py-2 rounded-md text-sm font-medium transition-colors"
            >
              {{ $t('navigation.logout') }}
            </button>
          </div>
        </div>
        
        <!-- 移动端导航菜单 -->
        <div v-if="mobileMenuOpen" class="md:hidden py-4 border-t border-gray-200">
          <div class="space-y-2">
            <NuxtLink
              to="/"
              @click="mobileMenuOpen = false"
              class="block px-3 py-2 rounded-md text-base font-medium text-gray-600 hover:text-blue-600 hover:bg-gray-50 transition-colors"
              :class="{ 'text-blue-600 bg-blue-50': $route.path === '/' }"
            >
              {{ $t('navigation.dashboard') }}
            </NuxtLink>
            <NuxtLink
              to="/servers"
              @click="mobileMenuOpen = false"
              class="block px-3 py-2 rounded-md text-base font-medium text-gray-600 hover:text-blue-600 hover:bg-gray-50 transition-colors"
              :class="{ 'text-blue-600 bg-blue-50': $route.path === '/servers' }"
            >
              {{ $t('navigation.servers') }}
            </NuxtLink>
            <div class="sm:hidden px-3 py-2 text-sm text-gray-500">
              {{ $t('navigation.user') }}：{{ authStore.user?.username }}
            </div>
          </div>
        </div>
      </div>
    </nav>
    
    <main class="flex-1 w-full py-4 sm:py-6 px-4 sm:px-6" :class="{ 'lg:px-8': $route.path !== '/servers', 'lg:px-4': $route.path === '/servers' }">
      <slot />
    </main>

    <!-- 固定底部商标信息 -->
    <AppFooter 
      :app-version="appVersion"
      @show-privacy-policy="showPrivacyPolicy = true"
      @show-terms-of-service="showTermsOfService = true"
    />

    <!-- 隐私政策模态框 -->
    <div v-if="showPrivacyPolicy" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click="showPrivacyPolicy = false">
      <div class="bg-white rounded-lg p-6 max-w-md mx-4" @click.stop>
        <h3 class="text-lg font-semibold mb-4">{{ $t('modals.privacyPolicy') }}</h3>
        <p class="text-sm text-gray-600 mb-4">
          {{ $t('modals.privacyPolicyContent') }}
        </p>
        <button 
          @click="showPrivacyPolicy = false"
          class="w-full bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700 transition-colors"
        >
          {{ $t('common.confirm') }}
        </button>
      </div>
    </div>

    <!-- 服务条款模态框 -->
    <div v-if="showTermsOfService" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click="showTermsOfService = false">
      <div class="bg-white rounded-lg p-6 max-w-md mx-4" @click.stop>
        <h3 class="text-lg font-semibold mb-4">{{ $t('modals.termsOfService') }}</h3>
        <p class="text-sm text-gray-600 mb-4">
          {{ $t('modals.termsOfServiceContent') }}
        </p>
        <button 
          @click="showTermsOfService = false"
          class="w-full bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700 transition-colors"
        >
          {{ $t('common.confirm') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth'

const authStore = useAuthStore()
const mobileMenuOpen = ref(false)
const showPrivacyPolicy = ref(false)
const showTermsOfService = ref(false)
const appVersion = ref('1.0.0')

// 在客户端初始化时从cookie恢复登录状态
onMounted(() => {
  authStore.initFromStorage()
})

const logout = () => {
  authStore.logout()
  navigateTo('/login')
}

// 监听路由变化，关闭移动端菜单
watch(() => useRoute().path, () => {
  mobileMenuOpen.value = false
})
</script> 