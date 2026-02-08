<template>
  <div class="avatrar-main-container">
    <!-- Heart Animation -->
    <div class="heart-connection">
      <Vue3Lottie :animationData="HeartJson" :height="180" :width="180" />
    </div>

    <!-- Boy -->
    <figure class="boy">
      <img
        :src="getAvatarSrc(boy)"
        :alt="boy.name"
        @error="handleAvatarError($event, boySrc)"
        draggable="false"
      />
      <figcaption>{{ boy.name }}</figcaption>
    </figure>

    <!-- Girl -->
    <figure class="girl">
      <img
        :src="getAvatarSrc(girl)"
        :alt="girl.name"
        @error="handleAvatarError($event, girlSrc)"
        draggable="false"
      />
      <figcaption>{{ girl.name }}</figcaption>
    </figure>
  </div>
</template>

<script setup lang="ts">
import { Vue3Lottie } from 'vue3-lottie'

import boySrc from '@/assets/images/boy.png'
import girlSrc from '@/assets/images/girl.png'
import HeartJson from '@/data/lottie/doubleHeart.json'

interface CoupleInfo {
  boy: {
    name: string
    avatar?: { thumbnail?: string; url?: string }
  }
  girl: {
    name: string
    avatar?: { thumbnail?: string; url?: string }
  }
}
const props = defineProps<CoupleInfo>()

const { boy, girl } = props

const handleAvatarError = (event: Event, src: string): void => {
  const target = event.target as HTMLImageElement
  target.src = src
  target.onerror = null // 防止无限循环
}

const getAvatarSrc = (person: { avatar?: { thumbnail?: string; url?: string } }): string => {
  return person.avatar?.thumbnail || person.avatar?.url || ''
}
</script>

<style scoped>
figure {
  width: 100%;
  aspect-ratio: 1;
  margin: 0 0 30px;
  padding: 5px 10px 0;
  box-sizing: border-box;
  display: grid;
  grid-template-rows: 100%;
  cursor: pointer;
  position: relative;
  border-radius: 50%;
  overflow: hidden;
}

figure:hover {
  overflow: visible;
}

figure::before {
  content: '';
  position: absolute;
  inset: 0;
  transform-origin: bottom;
  filter: brightness(0.9);
  transition: 0.5s;
  background: #f0ada0;
}

img {
  grid-area: 1 / 1;
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: top;
  filter: contrast(0.8) brightness(0.7);
  place-self: end center;
  transition: 0.5s;
}

figcaption {
  grid-area: 1 / 1;
  width: calc(100% + 20px);
  color: #fff;
  font-size: min(20px, 5vmin);
  text-align: center;
  place-self: end center;
  transform: perspective(500px) translateY(100%) rotateX(-90deg);
  backface-visibility: hidden;
  transform-origin: top;
  background: #984f31;
  transition: 0.5s;
}

figure:hover img {
  width: 180%;
  height: 110%;
  filter: contrast(1);
}

figure:hover::before {
  filter: brightness(0.3);
  transform: perspective(500px) rotateX(60deg);
}

figure:hover figcaption {
  transform: perspective(500px) translateY(100%) rotateX(-30deg);
}

.avatrar-main-container {
  height: 300px;
  display: grid;
  grid-auto-flow: column;
  grid-auto-columns: max(110px, 18vmin);
  place-content: end center;
  gap: min(300px, 30vmin);
  background: transparent;
  position: relative;
}

.heart-connection {
  position: absolute;
  left: 50%;
  top: 70%;
  transform: translate(-50%, -50%);
  width: 180px;
  height: 180px;
  z-index: 0;
}
</style>
