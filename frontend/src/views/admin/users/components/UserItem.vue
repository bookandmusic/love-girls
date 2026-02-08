<template>
  <CardItem>
    <template #header>
      <div class="flex items-center justify-between">
        <div></div>
        <div class="flex flex-shrink-0">
          <span
            :class="[
              'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
              user.role === 'admin' ? 'bg-red-100 text-red-800' : 'bg-green-100 text-green-800',
            ]"
          >
            {{ user.role === 'admin' ? '管理员' : '普通用户' }}
          </span>
        </div>
      </div>
    </template>

    <template #content>
      <div class="flex items-center pb-4">
        <div class="flex-shrink-0 h-12 w-12">
          <img
            v-if="user.avatar?.thumbnail"
            :src="user.avatar.thumbnail"
            :alt="user.name"
            class="h-12 w-12 rounded-full object-cover"
            @error="handleAvatarError($event, user)"
          />
          <div
            v-else-if="user.avatar?.url"
            class="h-12 w-12 rounded-full bg-indigo-100 flex items-center justify-center"
            :style="getAvatarBgStyle(user.name)"
          >
            <span class="text-indigo-800 font-medium">{{ getUserInitial(user.name) }}</span>
          </div>
          <div
            v-else
            class="h-12 w-12 rounded-full bg-indigo-100 flex items-center justify-center"
            :style="getAvatarBgStyle(user.name)"
          >
            <span class="text-indigo-800 font-medium">{{ getUserInitial(user.name) }}</span>
          </div>
        </div>
        <div class="ml-4 flex-1 min-w-0">
          <div class="text-base font-medium text-gray-900 truncate">{{ user.name }}</div>
          <div class="text-sm text-gray-500 truncate">{{ user.email }}</div>
          <div class="text-sm text-gray-500">{{ user.role }} | {{ user.joinDate }}</div>
        </div>
        <div class="ml-4">
          <button @click="onEdit">
            <BaseIcon name="edit" size="w-6" color="text-[#FFB61E]" />
          </button>
        </div>
      </div>
    </template>
  </CardItem>
</template>

<script setup lang="ts">
import BaseIcon from '@/components/ui/BaseIcon.vue'
import CardItem from '@/components/ui/CardItem.vue'
import type { User } from '@/services/userApi'

interface Props {
  user: User
}

interface Emits {
  (e: 'edit', user: User): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const onEdit = () => {
  emit('edit', props.user)
}

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
</script>
