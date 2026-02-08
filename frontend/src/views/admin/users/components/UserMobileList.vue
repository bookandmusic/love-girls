<template>
  <div class="overflow-hidden sm:rounded-md">
    <UserItem v-for="user in users" :key="user.id" :user="user" @edit="$emit('edit', user)" />
  </div>
</template>

<script setup lang="ts">
import type { User } from '@/services/userApi'

import UserItem from './UserItem.vue'

interface Props {
  users: User[]
}

interface Emits {
  (e: 'edit', user: User): void
}

defineProps<Props>()
defineEmits<Emits>()
</script>
