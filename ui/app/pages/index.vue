<template>
  <UContainer class="py-12">
    <main class="text-center space-y-12">
      <!-- 欢迎标题区域 -->
      <header class="space-y-6">
        <UAvatar 
          size="2xl" 
          icon="i-lucide-server" 
          :ui="{ background: 'bg-gradient-to-br from-blue-500 to-indigo-600' }"
          class="mx-auto mb-6"
        />
        <h1 class="text-5xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent">
          {{ $t('home.title') }}
        </h1>
        <p class="text-xl text-gray-600 max-w-3xl mx-auto leading-relaxed">
          {{ $t('home.subtitle') }}
        </p>
      </header>
      
      <!-- 系统信息卡片 -->
      <section>
        <UCard class="max-w-lg mx-auto">
          <template #header>
            <header class="flex items-center space-x-3">
              <UIcon name="i-lucide-user" class="w-5 h-5 text-blue-600" />
              <h2 class="text-xl font-semibold">{{ $t('home.systemInfo') }}</h2>
            </header>
          </template>
          
          <dl class="space-y-4">
            <div class="flex justify-between items-center py-2 border-b border-gray-100">
              <dt class="font-medium text-gray-700">{{ $t('home.username') }}</dt>
              <dd class="text-gray-900 font-mono">{{ authStore.user?.username }}</dd>
            </div>
            <div class="flex justify-between items-center py-2">
              <dt class="font-medium text-gray-700">{{ $t('home.userID') }}</dt>
              <dd class="text-gray-900 font-mono">#{{ authStore.user?.id }}</dd>
            </div>
          </dl>
        </UCard>
      </section>
      
      <!-- 功能卡片网格 -->
      <section class="space-y-8">
        <h2 class="text-2xl font-semibold text-gray-900">{{ $t('home.features') }}</h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-6xl mx-auto">
          <!-- 服务器管理卡片 -->
          <NuxtLink to="/servers">
            <UCard class="h-full hover:shadow-xl transition-all duration-300 cursor-pointer group border-2 border-transparent hover:border-blue-200">
              <article class="text-center space-y-6 p-6">
                <UAvatar 
                  size="lg" 
                  icon="i-lucide-server" 
                  :ui="{ background: 'bg-blue-100 group-hover:bg-blue-200' }"
                  class="mx-auto transition-colors"
                />
                <header>
                  <h3 class="text-xl font-semibold text-gray-900 mb-2">{{ $t('home.serverManagement') }}</h3>
                  <p class="text-gray-600 leading-relaxed">{{ $t('home.serverManagementDesc') }}</p>
                </header>
                <footer>
                  <UButton 
                    color="blue" 
                    variant="ghost" 
                    size="sm"
                    class="group-hover:bg-blue-50"
                  >
                    {{ $t('home.startManage') }}
                    <UIcon name="i-lucide-arrow-right" class="w-4 h-4 ml-1" />
                  </UButton>
                </footer>
              </article>
            </UCard>
          </NuxtLink>

          <!-- 占位卡片：玩家管理 -->
          <UCard class="h-full opacity-60 border-2 border-gray-100">
            <article class="text-center space-y-6 p-6">
              <UAvatar 
                size="lg" 
                icon="i-lucide-users" 
                :ui="{ background: 'bg-gray-100' }"
                class="mx-auto"
              />
              <header>
                <h3 class="text-xl font-semibold text-gray-500 mb-2">{{ $t('home.playerManagement') }}</h3>
                <p class="text-gray-400 leading-relaxed">{{ $t('home.playerManagementDesc') }}</p>
              </header>
              <footer>
                <UBadge color="gray" variant="soft" size="sm">
                  {{ $t('home.comingSoon') }}
                </UBadge>
              </footer>
            </article>
          </UCard>

          <!-- 占位卡片：日志监控 -->
          <UCard class="h-full opacity-60 border-2 border-gray-100">
            <article class="text-center space-y-6 p-6">
              <UAvatar 
                size="lg" 
                icon="i-lucide-bar-chart-3" 
                :ui="{ background: 'bg-gray-100' }"
                class="mx-auto"
              />
              <header>
                <h3 class="text-xl font-semibold text-gray-500 mb-2">{{ $t('home.logMonitoring') }}</h3>
                <p class="text-gray-400 leading-relaxed">{{ $t('home.logMonitoringDesc') }}</p>
              </header>
              <footer>
                <UBadge color="gray" variant="soft" size="sm">
                  {{ $t('home.comingSoon') }}
                </UBadge>
              </footer>
            </article>
          </UCard>
        </div>
      </section>
    </main>
  </UContainer>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  middleware: 'auth'
})

const authStore = useAuthStore()

onMounted(() => {
  authStore.initFromStorage()
})
</script> 