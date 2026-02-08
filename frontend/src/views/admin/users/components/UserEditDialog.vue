<template>
  <GenericDialog :open="open" title="编辑用户信息" @cancel="closeDialog" :loading="loading">
    <template #content>
      <!-- 用户名 -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">用户名</label>
        <input
          v-model="localUser.name"
          type="text"
          class="w-full win11-input"
          placeholder="请输入用户名"
        />
      </div>

      <!-- 邮箱 -->
      <div class="mt-4">
        <label class="block text-sm font-medium text-gray-700 mb-1">邮箱</label>
        <input
          v-model="localUser.email"
          type="email"
          class="w-full win11-input"
          placeholder="请输入邮箱"
        />
      </div>

      <!-- 头像上传 -->
      <div class="mt-4">
        <label class="block text-sm font-medium text-gray-700 mb-1">头像</label>
        <div class="flex items-center">
          <div class="flex-shrink-0 h-16 w-16 cursor-pointer" @click="triggerAvatarUpload">
            <img
              v-if="localUser.avatar?.thumbnail"
              :src="localUser.avatar.thumbnail"
              :alt="localUser.name"
              class="h-16 w-16 rounded-full object-cover"
              @error="handleAvatarError($event, localUser)"
            />
            <div
              v-else-if="localUser.avatar?.url"
              class="h-16 w-16 rounded-full bg-indigo-100 flex items-center justify-center"
              :style="getAvatarBgStyle(localUser.name)"
            >
              <span class="text-indigo-800 font-medium">{{ getUserInitial(localUser.name) }}</span>
            </div>
            <div
              v-else
              class="h-16 w-16 rounded-full bg-indigo-100 flex items-center justify-center"
              :style="getAvatarBgStyle(localUser.name)"
            >
              <span class="text-indigo-800 font-medium">{{ getUserInitial(localUser.name) }}</span>
            </div>
          </div>
          <div class="ml-4">
            <input
              ref="avatarInputRef"
              type="file"
              accept="image/*"
              @change="uploadAvatar"
              class="hidden"
            />
            <p class="text-xs text-gray-500 mt-1">支持 JPG, PNG 格式</p>
          </div>
        </div>
      </div>

      <!-- 密码 -->
      <div class="mt-4">
        <label class="block text-sm font-medium text-gray-700 mb-1">新密码</label>
        <input
          v-model="newPassword"
          type="password"
          placeholder="留空则不修改密码"
          class="w-full win11-input"
        />
      </div>

      <div class="mt-4">
        <label class="block text-sm font-medium text-gray-700 mb-1">确认密码</label>
        <input
          v-model="confirmPassword"
          type="password"
          placeholder="请再次输入密码"
          class="w-full win11-input"
        />
      </div>

      <div v-if="passwordError" class="mt-2 text-sm text-red-600">
        {{ passwordError }}
      </div>
    </template>
    <template #actions>
      <div class="w-full flex">
        <div class="flex-1 text-center cursor-pointer" @click="closeDialog">取消</div>
        <div
          class="w-1/2 border-l border-gray-300 text-center cursor-pointer text-blue-500"
          @click="handleSave()"
        >
          确认
        </div>
      </div>
    </template>
  </GenericDialog>

  <!-- 确认对话框 -->
  <GenericDialog
    :open="showConfirmDialog"
    title="确认保存"
    :loading="loading"
    size-class="max-w-md"
    @cancel="cancelConfirm"
  >
    <template #content>
      <p class="text-gray-700">您确定要保存对这个用户的更改吗？</p>
    </template>
    <template #actions>
      <div class="w-full flex">
        <div class="flex-1 text-center cursor-pointer" @click="cancelConfirm">取消</div>
        <div
          class="w-1/2 border-l border-gray-300 text-center cursor-pointer text-blue-500"
          @click="confirmSave"
        >
          确定
        </div>
      </div>
    </template>
  </GenericDialog>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue'

import GenericDialog from '@/components/ui/GenericDialog.vue'
import type { User, UserFormData } from '@/services/userApi'

interface Props {
  open: boolean
  loading: boolean
  user: User
}

interface Emits {
  (e: 'update:open', open: boolean): void
  (e: 'confirm', user: UserFormData): void
  (e: 'uploadAvatar', event: Event): void
  (e: 'cancel'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const closeDialog = () => {
  emit('update:open', false)
  emit('cancel')
}

const uploadAvatar = (event: Event) => {
  emit('uploadAvatar', event)
}

// 本地状态
const newPassword = ref('')
const confirmPassword = ref('')
const passwordError = ref('')
const avatarInputRef = ref<HTMLInputElement>()
const showConfirmDialog = ref(false)

// 创建本地响应式副本
const localUser = reactive<User>({ ...props.user })

// 监听props.user的变化，更新本地副本
watch(
  () => props.user,
  newUser => {
    Object.assign(localUser, newUser)
  },
  { deep: true }
)

// 获取用户名首字母
const getUserInitial = (name: string) => {
  return name ? name.charAt(0).toUpperCase() : '?'
}

// 获取头像背景样式
const getAvatarBgStyle = (name: string) => {
  if (!name) return { backgroundColor: '#8b5cf6' }

  // 使用用户名生成一个颜色值
  let hash = 0
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash)
  }

  // 生成HSL颜色值
  const hue = hash % 360
  return {
    backgroundColor: `hsl(${hue}, 70%, 60%)`,
    color: 'white',
  }
}

// 处理头像加载错误
const handleAvatarError = (event: Event, user: User) => {
  // 当头像加载失败时，显示默认头像
  const target = event.target as HTMLImageElement
  target.style.display = 'none'

  // 重新渲染为默认头像
  user.avatar = undefined
}

// 触发头像上传
const triggerAvatarUpload = () => {
  if (avatarInputRef.value) {
    // 重置文件输入框的值，这样即使用户选择同一个文件，也会触发change事件
    avatarInputRef.value.value = ''
    // 打开文件选择对话框
    avatarInputRef.value.click()
  }
}

// 处理保存事件
const handleSave = () => {
  // 验证密码
  if (newPassword.value !== confirmPassword.value) {
    passwordError.value = '两次输入的密码不一致'
    return
  }

  if (newPassword.value && newPassword.value.length < 6) {
    passwordError.value = '密码长度至少6位'
    return
  }

  // 显示二次确认对话框
  showConfirmDialog.value = true
}

// 取消二次确认
const cancelConfirm = () => {
  showConfirmDialog.value = false
}

// 确认保存
const confirmSave = () => {
  showConfirmDialog.value = false

  // 添加新密码到表单数据
  const userData = { ...localUser } as UserFormData
  if (newPassword.value) {
    userData.newPassword = newPassword.value
  }

  // 触发保存事件
  emit('confirm', userData)
}
</script>
