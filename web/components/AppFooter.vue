<template>
  <footer class="bg-white border-t border-gray-200">
    <div class="w-full px-4 sm:px-6 lg:px-8">
      <div class="py-4 sm:py-6">
        <div class="flex flex-col sm:flex-row justify-between items-center space-y-2 sm:space-y-0">
          <!-- 商标和版权信息 -->
          <div class="text-center sm:text-left">
            <p class="text-sm text-gray-600">
              © {{ new Date().getFullYear() }} 
              <span class="font-semibold text-gray-800">ARK服务器管理器</span>
              由 <span class="font-medium">{{ companyName }}</span> 开发
            </p>
          </div>
        </div>
        
        <!-- 附加链接 -->
        <div class="flex justify-center sm:justify-start space-x-6 mt-3 pt-3 border-t border-gray-100">
          <a 
            v-if="showPrivacyPolicy"
            href="#" 
            class="text-xs text-gray-500 hover:text-gray-700 transition-colors"
            @click.prevent="$emit('show-privacy-policy')"
          >
            隐私政策
          </a>
          <a 
            v-if="showTermsOfService"
            href="#" 
            class="text-xs text-gray-500 hover:text-gray-700 transition-colors"
            @click.prevent="$emit('show-terms-of-service')"
          >
            服务条款
          </a>
          <a 
            :href="supportEmail" 
            class="text-xs text-gray-500 hover:text-gray-700 transition-colors"
          >
            技术支持
          </a>
          <a 
            v-if="githubUrl"
            :href="githubUrl" 
            target="_blank"
            class="text-xs text-gray-500 hover:text-gray-700 transition-colors"
          >
            GitHub
          </a>
          <span v-if="pageInfo" class="text-xs text-gray-400">
            {{ pageInfo }}
          </span>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup>
// Props定义
const props = defineProps({
  companyName: {
    type: String,
    default: 'TBro'
  },
  appVersion: {
    type: String,
    default: '1.0.0'
  },
  supportEmail: {
    type: String,
    default: 'mailto:support@yourcompany.com'
  },
  githubUrl: {
    type: String,
    default: 'https://github.com/TBro1998/ArkServerManager'
  },
  showPrivacyPolicy: {
    type: Boolean,
    default: true
  },
  showTermsOfService: {
    type: Boolean,
    default: true
  },
  pageInfo: {
    type: String,
    default: ''
  }
})

// Events定义
defineEmits(['show-privacy-policy', 'show-terms-of-service'])
</script> 