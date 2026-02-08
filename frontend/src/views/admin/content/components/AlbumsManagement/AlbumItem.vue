<template>
  <CardItem>
    <!-- 背景图 -->
    <template #background>
      <div
        class="absolute inset-0 bg-cover bg-center"
        :style="{ backgroundImage: `url(${album.coverImage?.file?.thumbnail})` }"
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
          <BaseIcon name="photo" class="mr-1" size="w-4" color="text-[#FFC773]" />
          {{ album.photoCount }}
        </div>
        <div class="flex text-sm text-white/90">
          <BaseIcon name="calendar" class="mr-1" size="w-4" color="text-[#FFC773]" />
          {{ album.createdAt }}
        </div>
      </div>
      <p class="mt-3 text-md text-white/90">
        {{ album.description }}
      </p>
    </template>

    <!-- footer -->
    <template #footer>
      <div class="flex justify-end space-x-3">
        <button @click="viewAlbum">
          <BaseIcon name="open" size="w-6 h-6" color="text-[#FFB61E]" />
        </button>
        <button @click="editAlbum">
          <BaseIcon name="edit" size="w-6 h-6" color="text-[#FFB61E]" />
        </button>
        <button @click="deleteAlbum">
          <BaseIcon name="delete" size="w-6 h-6" color="text-[#FFB61E]" />
        </button>
      </div>
    </template>
  </CardItem>
</template>

<script setup lang="ts">
import BaseIcon from '@/components/ui/BaseIcon.vue'
import CardItem from '@/components/ui/CardItem.vue'
import { type Album } from '@/services/albumApi'

interface Props {
  album: Album
}

const props = defineProps<Props>()

const emit = defineEmits<{
  view: [album: Album]
  edit: [album: Album]
  delete: [album: Album]
}>()

const viewAlbum = () => {
  emit('view', props.album)
}

const editAlbum = () => {
  emit('edit', props.album)
}

const deleteAlbum = () => {
  emit('delete', props.album)
}
</script>
