<template>
  <div class="h-full w-full overflow-y-hidden flex flex-col">
    <h2 class="text-2xl font-bold text-gray-800 mb-6">站点设置</h2>

    <!-- 系统信息 -->
    <div class="bg-white/30 mb-6 rounded-xl shadow-sm border border-gray-100 p-3 md:p-6">
      <h3 class="text-lg font-semibold text-gray-800 mb-4">系统信息</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="flex justify-between items-center">
          <span class="text-sm text-gray-600">系统版本</span>
          <span class="text-sm font-medium text-gray-900">1.0.0</span>
        </div>
        <div class="flex justify-between items-center">
          <span class="text-sm text-gray-600">数据库版本</span>
          <span class="text-sm font-medium text-gray-900">SQLite 3.38</span>
        </div>
        <div class="flex justify-between items-center">
          <span class="text-sm text-gray-600">服务器环境</span>
          <span class="text-sm font-medium text-gray-900">生产环境</span>
        </div>
        <div class="flex justify-between items-center">
          <span class="text-sm text-gray-600">最后备份时间</span>
          <span class="text-sm font-medium text-gray-900">2023-08-25 14:30</span>
        </div>
      </div>
    </div>

    <!-- 站点信息设置 -->
    <div
      class="bg-white/30 flex flex-col flex-1 overflow-y-auto rounded-xl shadow-sm border border-gray-100 p-3 md:p-6"
    >
      <div class="mb-6">
        <h3 class="text-lg font-semibold text-gray-800 py-2">基本信息</h3>
        <p class="text-sm text-gray-500">配置站点的基本信息</p>
      </div>

      <form class="space-y-6 flex flex-col flex-1">
        <div>
          <label for="site-title" class="block text-sm font-medium text-gray-700 mb-2">
            站点标题 <span class="text-red-500">*</span>
          </label>
          <input
            type="text"
            id="site-title"
            v-model="settings.siteTitle"
            class="w-full win11-input"
            placeholder="请输入站点标题"
          />
        </div>

        <div>
          <label for="site-description" class="block text-sm font-medium text-gray-700 mb-2">
            站点描述
          </label>
          <input
            type="text"
            id="site-description"
            v-model="settings.siteDescription"
            class="w-full win11-input"
            placeholder="请输入站点描述"
          />
        </div>

        <div class="pt-4 border-t border-gray-300">
          <div class="w-full flex">
            <button
              type="button"
              @click="resetForm"
              class="flex-1 text-sm font-medium text-gray-700 focus:outline-none"
            >
              重置
            </button>
            <button
              type="submit"
              @click="confirmSave"
              class="flex-1 border-l border-gray-300 text-sm font-medium text-indigo-500 focus:outline-none"
            >
              保存设置
            </button>
          </div>
        </div>
      </form>
    </div>

    <!-- 保存确认对话框 -->
    <GenericDialog
      :open="showConfirmDialog"
      title="保存确认"
      :loading="uiStore.loading"
      size-class="max-w-md"
      @cancel="cancelSave"
    >
      <template #content>
        <p class="text-gray-700">您确定要保存这些设置吗？</p>
      </template>
      <template #actions>
        <div class="w-full flex">
          <div class="flex-1 text-center cursor-pointer" @click="cancelSave">取消</div>
          <div
            class="w-1/2 border-l border-gray-300 text-center cursor-pointer text-indigo-500"
            @click="saveSettingsConfirmed"
          >
            确定保存
          </div>
        </div>
      </template>
    </GenericDialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'

import GenericDialog from '@/components/ui/GenericDialog.vue'
import { type SiteSettings, systemApi } from '@/services/system'
import { useUIStore } from '@/stores/ui'
import { useToast } from '@/utils/toastUtils'

// 设置数据
const settings = reactive<SiteSettings>({
  siteTitle: '',
  siteDescription: '',
})

// 原始设置数据，用于重置
const originalSettings = ref<SiteSettings>({
  siteTitle: '',
  siteDescription: '',
})

const uiStore = useUIStore()
const showToast = useToast()

// 保存确认对话框
const showConfirmDialog = ref(false)

// 加载设置数据
const loadSettings = async () => {
  uiStore.setLoading(true)
  try {
    const response = await systemApi.getSiteSettings()
    if (response.data.code === 0) {
      const data = response.data.data
      Object.assign(settings, data)
      Object.assign(originalSettings.value, data)
    } else {
      showToast('获取设置失败', 'error')
    }
  } catch {
    showToast('获取设置失败', 'error')
  } finally {
    uiStore.setLoading(false)
  }
}

// 显示保存确认对话框
const confirmSave = (e: Event) => {
  e.preventDefault()
  if (!settings.siteTitle.trim()) {
    showToast('站点标题不能为空', 'error')
    return
  }
  showConfirmDialog.value = true
}

// 确认保存设置
const saveSettingsConfirmed = async () => {
  closeConfirmDialog()
  uiStore.setLoading(true)
  try {
    const response = await systemApi.saveSiteSettings(settings)
    if (response.data.code === 0) {
      showToast('设置保存成功', 'success')
      Object.assign(originalSettings.value, { ...settings })
    } else {
      showToast('保存设置失败', 'error')
    }
  } catch {
    showToast('保存设置失败', 'error')
  } finally {
    uiStore.setLoading(false)
  }
}

// 取消保存设置
const cancelSave = () => {
  closeConfirmDialog()
}

// 关闭确认对话框
const closeConfirmDialog = () => {
  showConfirmDialog.value = false
}

// 重置表单
const resetForm = () => {
  Object.assign(settings, originalSettings.value)
}

onMounted(() => {
  loadSettings()
})
</script>
