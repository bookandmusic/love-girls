<template>
  <div class="bg-white/30 shadow overflow-hidden sm:rounded-md">
    <table class="min-w-full divide-y divide-gray-200">
      <thead class="bg-gray-50">
        <tr>
          <th
            scope="col"
            class="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider"
          >
            用户
          </th>
          <th
            scope="col"
            class="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider"
          >
            角色
          </th>
          <th
            scope="col"
            class="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider"
          >
            加入时间
          </th>
          <th
            scope="col"
            class="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider"
          >
            操作
          </th>
        </tr>
      </thead>
      <tbody class="bg-white/30 divide-y divide-gray-200">
        <tr v-for="user in users" :key="user.id" class="hover:bg-gray-50">
          <td class="px-6 py-4 whitespace-nowrap text-center">
            <div class="flex items-center justify-center">
              <div class="flex-shrink-0 h-10 w-10">
                <img
                  v-if="user.avatar?.thumbnail"
                  :src="user.avatar.thumbnail"
                  :alt="user.name"
                  class="h-10 w-10 rounded-full object-cover"
                  @error="handleAvatarError($event, user)"
                />
                <img
                  v-else-if="user.avatar?.url"
                  :src="user.avatar.url"
                  :alt="user.name"
                  class="h-10 w-10 rounded-full object-cover"
                  @error="handleAvatarError($event, user)"
                />
                <div
                  v-else
                  class="h-10 w-10 rounded-full bg-indigo-100 flex items-center justify-center"
                  :style="getAvatarBgStyle(user.name)"
                >
                  <span class="text-indigo-800 font-medium">{{ getUserInitial(user.name) }}</span>
                </div>
              </div>
              <div class="ml-4 text-left">
                <div class="text-sm font-medium text-gray-900">{{ user.name }}</div>
                <div class="text-sm text-gray-500">{{ user.email }}</div>
              </div>
            </div>
          </td>
          <td class="px-6 py-4 whitespace-nowrap text-center">
            <div class="text-sm text-gray-900">{{ user.role }}</div>
          </td>
          <td class="px-6 py-4 whitespace-nowrap text-center text-sm text-gray-500">
            {{ user.joinDate }}
          </td>
          <td class="px-6 py-4 whitespace-nowrap text-center text-sm font-medium">
            <button @click="onEdit(user)" class="text-blue-600 hover:text-blue-900">
              <BaseIcon name="edit" size="w-6" color="text-[#FFB61E]" />
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import BaseIcon from '@/components/ui/BaseIcon.vue'
import type { User } from '@/services/userApi'

interface Props {
  users: User[]
}

interface Emits {
  (e: 'edit', user: User): void
}

defineProps<Props>()
const emit = defineEmits<Emits>()

const onEdit = (user: User) => {
  emit('edit', user)
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
