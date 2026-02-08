<template>
  <div class="flex-grow overflow-y-auto">
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 p-4">
      <CardItem v-for="album in albums" :key="album.id" @click="onSelectAlbum(album)">
        <!-- 背景图 -->
        <template #background>
          <div
            class="absolute inset-0 bg-cover bg-center"
            :style="{
              backgroundImage: `url(${album.coverImage?.file?.thumbnail || album.coverImage?.file?.url || ''})`,
            }"
          />
        </template>

        <!-- 遮罩 -->
        <template #overlay>
          <div class="absolute inset-0 bg-black/45"></div>
        </template>

        <!-- header -->
        <template #header>
          <p class="text-xl font-medium text-white truncate">
            {{ album.name }}
          </p>
        </template>

        <!-- content -->
        <template #content>
          <div class="flex items-center space-x-2">
            <div class="flex text-sm text-white/90">
              <BaseIcon name="photo-heart" size="w-5 h-5" class="mr-1" />
              {{ album.photoCount }}
            </div>
            <div class="flex text-sm text-white/90">
              <BaseIcon name="calendar" size="w-5 h-5" class="mr-1" />
              {{ album.createdAt }}
            </div>
          </div>
          <p class="mt-3 text-sm text-white/90">
            {{ album.description }}
          </p>
        </template>
      </CardItem>
    </div>
  </div>
</template>

<script setup lang="ts">
import BaseIcon from '@/components/ui/BaseIcon.vue'
import CardItem from '@/components/ui/CardItem.vue'
import type { Album } from '@/services/albumApi'

interface Props {
  albums: Album[]
}

defineProps<Props>()

const emit = defineEmits<{
  (e: 'select-album', album: Album): void
}>()

const onSelectAlbum = (album: Album) => {
  emit('select-album', album)
}
</script>
