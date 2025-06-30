export default defineNuxtRouteMiddleware((to, from) => {
  const authStore = useAuthStore()
  
  // 如果用户未认证且不是在登录页面，重定向到登录页面
  if (!authStore.isAuthenticated && to.path !== '/login' && to.path !== '/init') {
    return navigateTo('/login')
  }
  
  // 如果用户已认证且在登录页面，重定向到首页
  if (authStore.isAuthenticated && (to.path === '/login' || to.path === '/init')) {
    return navigateTo('/')
  }
})