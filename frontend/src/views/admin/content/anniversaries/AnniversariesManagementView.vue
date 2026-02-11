<template>
  <div class="h-full flex flex-col">
    <h3 class="text-2xl font-bold text-gray-800 mb-6">纪念日管理</h3>
    <AnniversariesManagement :trigger-add="addTrigger" />
  </div>
</template>

<script setup lang="ts">
// 组件名称 AnniversariesManagementView
import { watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useAddTrigger } from '@/utils/useAddTrigger'

import AnniversariesManagement from '../components/AnniversariesManagement.vue'

const route = useRoute()
const router = useRouter()

const { trigger: addTrigger, fire: handleAddClick } = useAddTrigger()

// 处理添加按钮点击
const handleAddAction = () => {
  // 监听查询参数来触发添加
  if (route.query.action === 'add') {
    handleAddClick()
    // 清除查询参数
    router.replace({ query: {} })
  }
}

// 监听路由变化
watch(
  () => route.query,
  () => {
    handleAddAction()
  },
  { immediate: true }
)
</script>
