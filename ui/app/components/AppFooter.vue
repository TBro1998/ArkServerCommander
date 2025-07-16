<template>
  <footer class="bg-white border-t border-gray-200">
    <div class="w-full px-4 sm:px-6 lg:px-8">
      <div class="py-4 sm:py-6">
        <div class="flex flex-col sm:flex-row justify-between items-center space-y-2 sm:space-y-0">
          <!-- 商标和版权信息 -->
          <div class="text-center sm:text-left">
            <p class="text-sm text-gray-600">
              {{ $t('footer.copyright', { year: new Date().getFullYear(), company: companyName }) }}
            </p>
          </div>
        </div>
        
        <!-- 附加链接 -->
        <div class="flex justify-center sm:justify-start space-x-6 mt-3 pt-3 border-t border-gray-100">
          <a 
            v-if="showPrivacyPolicy"
            href="#" 
            class="text-xs text-gray-500 hover:text-gray-700 transition-colors"
            @click.prevent="showPrivacyPolicyModal = true"
          >
            {{ $t('footer.privacyPolicy') }}
          </a>
          <a 
            v-if="showTermsOfService"
            href="#" 
            class="text-xs text-gray-500 hover:text-gray-700 transition-colors"
            @click.prevent="showTermsOfServiceModal = true"
          >
            {{ $t('footer.termsOfService') }}
          </a>
          <a 
            v-if="githubUrl"
            :href="githubUrl" 
            target="_blank"
            class="text-xs text-gray-500 hover:text-gray-700 transition-colors"
          >
            {{ $t('footer.github') }}
          </a>
          <span v-if="pageInfo" class="text-xs text-gray-400">
            {{ pageInfo }}
          </span>
        </div>
      </div>
    </div>

    <!-- 隐私政策模态框 -->
    <UModal 
      v-model:open="showPrivacyPolicyModal" 
      :title="$t('modals.privacyPolicy')"
    >
      <template #body>
        <p class="text-sm text-gray-600">
          {{ $t('modals.privacyPolicyContent') }}
        </p>
      </template>
      
      <template #footer>
        <div class="flex justify-end">
          <UButton 
            color="primary" 
            :label="$t('common.confirm')"
            @click="showPrivacyPolicyModal = false"
          />
        </div>
      </template>
    </UModal>

    <!-- 服务条款模态框 -->
    <UModal 
      v-model:open="showTermsOfServiceModal" 
      :title="$t('modals.termsOfService')"
    >
      <template #body>
        <p class="text-sm text-gray-600">
          {{ $t('modals.termsOfServiceContent') }}
        </p>
      </template>
      
      <template #footer>
        <div class="flex justify-end">
          <UButton 
            color="primary" 
            :label="$t('common.confirm')"
            @click="showTermsOfServiceModal = false"
          />
        </div>
      </template>
    </UModal>
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

// 模态框状态
const showPrivacyPolicyModal = ref(false)
const showTermsOfServiceModal = ref(false)
</script> 