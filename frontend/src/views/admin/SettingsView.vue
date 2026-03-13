<template>
  <div class="flex flex-col pb-10">
    <h2 class="text-2xl font-bold admin-text-primary mb-6 font-(family-name:--font-signature)">
      站点设置
    </h2>

    <div class="admin-card overflow-hidden mb-8">
      <div class="px-6 py-4 border-b border-black/5 bg-black/5">
        <h3 class="text-sm font-bold admin-text-secondary uppercase tracking-wider">系统信息</h3>
      </div>
      <div class="divide-y divide-black/5">
        <div class="px-6 py-4 flex justify-between items-center hover:bg-black/5 transition-colors">
          <span class="text-sm font-medium admin-text-primary">系统版本</span>
          <span class="text-sm font-bold admin-text-secondary bg-white/50 px-2 py-0.5 rounded-md"
            >1.0.0</span
          >
        </div>
        <div class="px-6 py-4 flex justify-between items-center hover:bg-black/5 transition-colors">
          <span class="text-sm font-medium admin-text-primary">数据库版本</span>
          <span class="text-sm font-bold admin-text-secondary bg-white/50 px-2 py-0.5 rounded-md"
            >SQLite 3.38</span
          >
        </div>
        <div class="px-6 py-4 flex justify-between items-center hover:bg-black/5 transition-colors">
          <span class="text-sm font-medium admin-text-primary">服务器环境</span>
          <span class="text-sm font-bold text-green-600 bg-green-50 px-2 py-0.5 rounded-md"
            >生产环境</span
          >
        </div>
      </div>
    </div>

    <div class="admin-card overflow-hidden">
      <div class="px-6 py-4 border-b border-black/5 bg-black/5 flex justify-between items-center">
        <div>
          <h3 class="text-sm font-bold admin-text-secondary uppercase tracking-wider">基本信息</h3>
          <p class="text-xs admin-text-muted mt-0.5">配置站点的基本信息</p>
        </div>
      </div>

      <form class="p-6 space-y-6">
        <div class="space-y-2">
          <label for="site-title" class="block text-sm font-bold admin-text-primary px-1">
            站点标题 <span class="text-red-500">*</span>
          </label>
          <input
            type="text"
            id="site-title"
            v-model="settings.siteTitle"
            class="admin-input"
            placeholder="请输入站点标题"
          />
        </div>

        <div class="space-y-2">
          <label for="site-description" class="block text-sm font-bold admin-text-primary px-1">
            站点描述
          </label>
          <input
            type="text"
            id="site-description"
            v-model="settings.siteDescription"
            class="admin-input"
            placeholder="请输入站点描述"
          />
        </div>

        <div class="space-y-2">
          <label for="start-date" class="block text-sm font-bold admin-text-primary px-1">
            故事开始日期
          </label>
          <div class="relative">
            <input type="date" id="start-date" v-model="settings.startDate" class="admin-input" />
          </div>
        </div>

        <div class="pt-6 border-t border-black/5 flex gap-4">
          <button type="button" @click="resetForm" class="admin-btn-secondary flex-1 py-3">
            重置
          </button>
          <button type="submit" @click="confirmSave" class="admin-btn flex-1 py-3">保存设置</button>
        </div>
      </form>
    </div>

    <GenericDialog
      variant="admin"
      :open="showConfirmDialog"
      title="保存确认"
      :loading="uiStore.loading"
      size-class="max-w-md"
      @cancel="cancelSave"
    >
      <template #content>
        <p class="admin-text-primary">您确定要保存这些设置吗？</p>
      </template>
      <template #actions>
        <div class="w-full flex">
          <div
            class="flex-1 text-center cursor-pointer admin-text-secondary hover:admin-text-primary transition"
            @click="cancelSave"
          >
            取消
          </div>
          <div
            class="w-1/2 border-l border-white/60 text-center cursor-pointer text-[#f0ada0] hover:text-[#d89388] transition"
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

const settings = reactive<SiteSettings>({
  siteTitle: '',
  siteDescription: '',
  startDate: '',
})

const originalSettings = ref<SiteSettings>({
  siteTitle: '',
  siteDescription: '',
  startDate: '',
})

const uiStore = useUIStore()
const showToast = useToast()

const showConfirmDialog = ref(false)

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

const confirmSave = (e: Event) => {
  e.preventDefault()
  if (!settings.siteTitle.trim()) {
    showToast('站点标题不能为空', 'error')
    return
  }
  showConfirmDialog.value = true
}

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

const cancelSave = () => {
  closeConfirmDialog()
}

const closeConfirmDialog = () => {
  showConfirmDialog.value = false
}

const resetForm = () => {
  Object.assign(settings, originalSettings.value)
}

onMounted(() => {
  loadSettings()
})
</script>
