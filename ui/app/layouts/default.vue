<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <nav v-if="authStore.isAuthenticated" class="bg-white shadow-sm border-b">
      <div class="w-full px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center space-x-8">
            <NuxtLink to="/" class="text-xl font-semibold text-gray-900 hover:text-blue-600 transition-colors">
              {{ $t('auth.loginTitle') }}
            </NuxtLink>
            
            <!-- 桌面端导航菜单 - 使用 Nuxt UI NavigationMenu -->
            <UNavigationMenu 
              :items="navigationItems" 
              class="hidden md:flex"
              orientation="horizontal"
              variant="pill"
            />
          </div>
          
          <div class="flex items-center space-x-4">
            <span class="hidden sm:block text-sm text-gray-700">{{ $t('navigation.welcome') }}，{{ authStore.user?.username }}</span>
            
            <!-- 语言切换器 -->
            <LanguageSwitcher />
            
            <!-- 移动端菜单按钮 -->
            <UButton
              @click="mobileMenuOpen = !mobileMenuOpen"
              variant="ghost"
              size="sm"
              class="md:hidden"
              :icon="mobileMenuOpen ? 'i-heroicons-x-mark' : 'i-heroicons-bars-3'"
            />
            
            <UButton
              @click="logout"
              color="red"
              variant="solid"
              size="sm"
              :label="$t('navigation.logout')"
            />
          </div>
        </div>
        
        <!-- 移动端导航菜单 - 使用 Nuxt UI NavigationMenu 垂直布局 -->
        <UCollapsible v-model:open="mobileMenuOpen" class="md:hidden">
          <div class="py-4 border-t border-gray-200">
            <UNavigationMenu 
              :items="mobileNavigationItems" 
              orientation="vertical"
              variant="link"
              class="space-y-1"
            />
            <div class="sm:hidden px-3 py-2 text-sm text-gray-500 mt-2">
              {{ $t('navigation.user') }}：{{ authStore.user?.username }}
            </div>
          </div>
        </UCollapsible>
      </div>
    </nav>
    
    <main class="flex-1 w-full py-4 sm:py-6 px-4 sm:px-6" :class="{ 'lg:px-8': $route.path !== '/servers', 'lg:px-4': $route.path === '/servers' }">
      <slot />
    </main>

    <!-- 固定底部商标信息 -->
    <AppFooter :app-version="appVersion" />
  </div>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth'

const authStore = useAuthStore()
const route = useRoute()
const { t } = useI18n()
const mobileMenuOpen = ref(false)
const appVersion = ref('1.0.0')

// 桌面端导航菜单配置
const navigationItems = computed(() => [
  {
    label: t('navigation.dashboard'),
    to: '/',
    icon: 'i-heroicons-home',
    active: route.path === '/'
  },
  {
    label: t('navigation.servers'),
    to: '/servers',
    icon: 'i-heroicons-server',
    active: route.path === '/servers'
  }
])

// 移动端导航菜单配置
const mobileNavigationItems = computed(() => [
  {
    label: t('navigation.dashboard'),
    to: '/',
    icon: 'i-heroicons-home',
    active: route.path === '/',
    onSelect: () => {
      mobileMenuOpen.value = false
      navigateTo('/')
    }
  },
  {
    label: t('navigation.servers'),
    to: '/servers',
    icon: 'i-heroicons-server',
    active: route.path === '/servers',
    onSelect: () => {
      mobileMenuOpen.value = false
      navigateTo('/servers')
    }
  }
])

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