import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware((to, from) => {
  const authStore = useAuthStore()
  
  // 如果auth store没有初始化，先从存储中恢复
  if (!authStore.isAuthenticated) {
    authStore.initFromStorage()
  }
  
  // 再次检查认证状态
  const isAuthenticated = authStore.isAuthenticated
  
  // 调试信息
  console.log('认证中间件检查:', {
    path: to.path,
    isAuthenticated,
    hasToken: !!authStore.token
  })
  
  // 如果用户未认证且不是在登录页面，重定向到登录页面
  if (!isAuthenticated && to.path !== '/login' && to.path !== '/init') {
    console.log('用户未认证，重定向到登录页面')
    return navigateTo('/login')
  }
  
  // 如果用户已认证且在登录页面，重定向到首页
  if (isAuthenticated && (to.path === '/login' || to.path === '/init')) {
    console.log('用户已认证，重定向到首页')
    return navigateTo('/')
  }
})