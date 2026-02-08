<template>
  <div class="w-full h-full flex flex-col overflow-hidden">
    <!-- 纪念日列表 -->
    <ul class="flex-1 min-h-0 overflow-y-auto pr-2">
      <AnniversaryItem
        v-for="anniversary in anniversaries"
        :key="anniversary.id"
        :anniversary="anniversary"
        @edit="editAnniversary"
        @delete="confirmDelete = $event"
      />
    </ul>

    <!-- 分页 -->
    <div class="pt-3 flex-shrink-0">
      <Pagination
        :current-page="currentPage"
        :total-pages="totalPages"
        @prev="handlePrevPage"
        @next="handleNextPage"
      />
    </div>

    <!-- 编辑对话框 -->
    <AnniversaryEditDialog
      v-model:open="showEditDialog"
      :anniversary="editingAnniversary"
      :loading="uiStore.loading"
      @confirm="saveAnniversary"
      @cancel="closeDialog"
    />

    <!-- 删除确认对话框 -->
    <GenericDialog
      :open="!!confirmDelete"
      title="删除确认"
      :loading="uiStore.loading"
      size-class="max-w-md"
    >
      <template #content>
        <p class="text-gray-700">您确定要删除这个纪念日吗？此操作不可恢复。</p>
      </template>
      <template #actions>
        <div class="w-full flex">
          <div class="flex-1 text-center cursor-pointer" @click="cancelDelete">取消</div>
          <div
            class="w-1/2 border-l border-gray-300 text-center cursor-pointer text-red-500"
            @click="performDelete"
          >
            确定删除
          </div>
        </div>
      </template>
    </GenericDialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'

import GenericDialog from '@/components/ui/GenericDialog.vue'
import Pagination from '@/components/ui/Pagination.vue'
import { type Anniversary, anniversaryApi } from '@/services/anniversaryApi'
import { useUIStore } from '@/stores/ui'
import { useToast } from '@/utils/toastUtils'

import AnniversaryEditDialog from './AnniversariesManagement/AnniversaryEditDialog.vue'
import AnniversaryItem from './AnniversariesManagement/AnniversaryItem.vue'

// 定义默认纪念日对象
const DEFAULT_ANNIVERSARY: Anniversary = {
  id: 0,
  title: '',
  date: new Date().toISOString().substring(0, 10),
  description: '',
  calendar: 'solar' as 'solar' | 'lunar',
}

const props = defineProps<{ triggerAdd: boolean }>()
watch(
  () => props.triggerAdd,
  val => {
    if (val) addNewAnniversary()
  }
)

const uiStore = useUIStore()
const showToast = useToast()

const anniversaries = ref<Anniversary[]>([])
const totalAnniversaries = ref(0)
const currentPage = ref(1)
const pageSize = ref(5)

// 计算总页数
const totalPages = computed(() => Math.ceil(totalAnniversaries.value / pageSize.value) || 1)

// 编辑相关
const showEditDialog = ref(false)
const editingAnniversary = ref<Anniversary | null>(null)

// 删除确认
const confirmDelete = ref<Anniversary | null>(null)

// 取消删除
const cancelDelete = () => {
  confirmDelete.value = null
}

// 加载纪念日列表
const loadAnniversaries = async () => {
  uiStore.setLoading(true)
  try {
    const response = await anniversaryApi.getAnniversaries(currentPage.value, pageSize.value)

    // 为了处理可能的API响应，我们需要确保所有Anniversary对象都有有效的id
    anniversaries.value = response.data.anniversaries.map(mem => ({
      ...mem,
      id: mem.id ? mem.id : 0, // 确保id不会是undefined
    }))

    totalAnniversaries.value =
      response.data.total || response.data.totalCount || response.data.anniversaries.length
  } catch (error) {
    console.error('加载纪念日失败:', error)
    showToast('加载纪念日失败', 'error')
  } finally {
    uiStore.setLoading(false)
  }
}

// 编辑纪念日
const editAnniversary = (anniversary: Anniversary) => {
  // 复制anniversary对象时确保id是有效数字
  editingAnniversary.value = {
    ...anniversary,
    id: anniversary.id ? anniversary.id : 0,
  }
  showEditDialog.value = true
}

// 保存纪念日（新增或编辑）
const saveAnniversary = async (anniversaryData: Anniversary) => {
  uiStore.setLoading(true)
  try {
    if (anniversaryData.id && anniversaryData.id > 0) {
      // 编辑现有纪念日
      const response = await anniversaryApi.updateAnniversary(anniversaryData.id, {
        title: anniversaryData.title,
        date: anniversaryData.date,
        description: anniversaryData.description,
        calendar: anniversaryData.calendar,
      })

      // 更新本地数据
      const index = anniversaries.value.findIndex(m => m.id === anniversaryData.id)
      if (index !== -1) {
        anniversaries.value[index] = { ...response.data }
      }
      showToast('纪念日更新成功', 'success')
    } else {
      // 添加新纪念日
      const response = await anniversaryApi.createAnniversary({
        title: anniversaryData.title,
        date: anniversaryData.date,
        description: anniversaryData.description,
        calendar: anniversaryData.calendar,
      })

      // 添加到列表开头
      anniversaries.value.unshift(response.data)
      totalAnniversaries.value++

      // 重置到第一页，确保新添加的纪念日显示
      currentPage.value = 1
      showToast('纪念日添加成功', 'success')
    }
  } catch (error) {
    console.error('保存纪念日失败:', error)
    showToast('保存纪念日失败', 'error')
  } finally {
    uiStore.setLoading(false)
    // 关闭对话框并重置表单
    closeDialog()
  }
}

// 删除纪念日
const performDelete = async () => {
  if (!confirmDelete.value || !confirmDelete.value.id) return

  const anniversaryId = confirmDelete.value.id
  uiStore.setLoading(true)
  try {
    // 发送API请求删除纪念日
    await anniversaryApi.deleteAnniversary(anniversaryId)

    const index = anniversaries.value.findIndex(m => m.id === anniversaryId)
    if (index !== -1) {
      anniversaries.value.splice(index, 1)
      totalAnniversaries.value--

      // 检查当前页是否为空，如果是，则跳转到前一页（如果存在）
      if (anniversaries.value.length === 0 && currentPage.value > 1) {
        currentPage.value--
      }
      // 重新加载当前页数据
      await loadAnniversaries()
    }
    showToast('纪念日删除成功', 'success')
  } catch (error) {
    console.error('删除纪念日失败:', error)
    showToast('删除失败', 'error')
  } finally {
    uiStore.setLoading(false)
    confirmDelete.value = null
  }
}

// 添加新纪念日
const addNewAnniversary = () => {
  editingAnniversary.value = { ...DEFAULT_ANNIVERSARY }
  showEditDialog.value = true
}

// 关闭对话框
const closeDialog = () => {
  showEditDialog.value = false
  editingAnniversary.value = null
}

// 上一页
const handlePrevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    loadAnniversaries()
  }
}

// 下一页
const handleNextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    loadAnniversaries()
  }
}

onMounted(() => {
  loadAnniversaries()
})
</script>
