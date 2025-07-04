<template>
  <div class="flex items-center gap-3">
    <div class="relative">
      <input
        :id="inputId"
        :checked="modelValue"
        @change="$emit('update:modelValue', $event.target.checked)"
        type="checkbox"
        class="sr-only"
      />
      <label 
        :for="inputId" 
        class="flex items-center cursor-pointer"
      >
        <div class="relative">
          <div class="w-10 h-6 bg-gray-200 rounded-full shadow-inner transition-colors duration-200 ease-in-out"
               :class="modelValue ? 'bg-green-500' : 'bg-gray-300'">
          </div>
          <div class="absolute left-1 top-1 bg-white w-4 h-4 rounded-full shadow-md transform transition-transform duration-200 ease-in-out"
               :class="modelValue ? 'translate-x-4' : 'translate-x-0'">
          </div>
        </div>
        <span v-if="label" class="ml-3 text-sm font-medium transition-colors duration-200"
              :class="modelValue ? 'text-green-600' : 'text-gray-700'">
          {{ label }}
        </span>
      </label>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

// Props
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  label: {
    type: String,
    default: ''
  },
  id: {
    type: String,
    default: ''
  }
})

// Emits
const emit = defineEmits(['update:modelValue'])

// 计算唯一的输入ID
const inputId = computed(() => {
  return props.id || `toggle-${Math.random().toString(36).substr(2, 9)}`
})
</script> 