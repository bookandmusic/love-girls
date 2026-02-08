<template>
  <div class="h-full flex flex-col overflow-y-auto">
    <div class="sticky top-0 z-10 generic-card p-4">
      <button @click="onBack" class="flex items-center mr-4">
        <BaseIcon name="left" size="w-6 h-6" />
        返回相册
      </button>
    </div>

    <div v-if="photos.length > 0" class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4 p-4">
      <vue-easy-lightbox
        :visible="visibleRef"
        :imgs="imgsRef"
        :index="indexRef"
        @hide="onHide"
      ></vue-easy-lightbox>
      <div
        v-for="photo in photos"
        :key="photo.id"
        @click="preview(photo.file?.url || '')"
        class="aspect-square overflow-hidden rounded-lg border border-white/80 bg-white/50 hover:bg-white/80 transition-all cursor-pointer transform hover:scale-105 shadow-md hover:shadow-lg"
      >
        <img
          :src="photo.file?.thumbnail || photo.file?.url || ''"
          :alt="photo.alt"
          class="w-full h-full object-cover"
        />
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="flex-1 flex flex-col items-center justify-center text-[#FF7500]">
      <BaseIcon name="photo-heart" size="w-24" color="text-[#FF7500]" />
      <p class="text-xl font-medium mt-4">暂无照片</p>
      <p class="text-md mt-2">还没有添加任何照片</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import VueEasyLightbox from 'vue-easy-lightbox'

import BaseIcon from '@/components/ui/BaseIcon.vue'
import type { Photo } from '@/services/albumApi'

interface Props {
  photos: Photo[]
}

defineProps<Props>()

const emit = defineEmits<{
  (e: 'back'): void
}>()

const onBack = () => {
  emit('back')
}

const visibleRef = ref(false)
const indexRef = ref(0) // default 0
const imgsRef = ref('')
const onShow = () => {
  visibleRef.value = true
}
const onHide = () => (visibleRef.value = false)
const preview = (url: string) => {
  imgsRef.value = url
  onShow()
}
</script>
