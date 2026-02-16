<template>
  <div class="h-full w-full overflow-y-auto flex flex-col pb-4">
    <h2 class="text-2xl font-bold admin-text-primary mb-6 font-(family-name:--font-signature)">
      站点设置
    </h2>

    <div class="admin-card p-6 mb-6">
      <h3 class="text-lg font-semibold admin-text-primary mb-4">系统信息</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="flex justify-between items-center">
          <span class="text-sm admin-text-secondary">系统版本</span>
          <span class="text-sm font-medium admin-text-primary">1.0.0</span>
        </div>
        <div class="flex justify-between items-center">
          <span class="text-sm admin-text-secondary">数据库版本</span>
          <span class="text-sm font-medium admin-text-primary">SQLite 3.38</span>
        </div>
        <div class="flex justify-between items-center">
          <span class="text-sm admin-text-secondary">服务器环境</span>
          <span class="text-sm font-medium admin-text-primary">生产环境</span>
        </div>
        <div class="flex justify-between items-center">
          <span class="text-sm admin-text-secondary">最后备份时间</span>
          <span class="text-sm font-medium admin-text-primary">2023-08-25 14:30</span>
        </div>
      </div>
    </div>

    <div class="admin-card p-6">
      <div class="mb-6">
        <h3 class="text-lg font-semibold admin-text-primary py-2">基本信息</h3>
        <p class="text-sm admin-text-muted">配置站点的基本信息</p>
      </div>

      <form class="space-y-6">
        <div>
          <label for="site-title" class="block text-sm font-medium admin-text-primary mb-2">
            站点标题 <span class="text-[#E8A8A8]">*</span>
          </label>
          <input
            type="text"
            id="site-title"
            v-model="settings.siteTitle"
            class="admin-input"
            placeholder="请输入站点标题"
          />
        </div>

        <div>
          <label for="site-description" class="block text-sm font-medium admin-text-primary mb-2">
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

        <div>
          <label for="start-date" class="block text-sm font-medium admin-text-primary mb-2">
            故事开始日期
          </label>
          <input type="date" id="start-date" v-model="settings.startDate" class="admin-input" />
        </div>

        <div class="pt-4 border-t border-white/60 flex gap-4">
          <button type="button" @click="resetForm" class="admin-btn-secondary flex-1">重置</button>
          <button type="submit" @click="confirmSave" class="admin-btn flex-1">保存设置</button>
        </div>
      </form>
    </div>

    <GenericDialog
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
