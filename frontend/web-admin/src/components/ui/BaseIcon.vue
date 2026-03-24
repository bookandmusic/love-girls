<template>
  <component
    :is="SvgComponent"
    v-if="SvgComponent"
    :class="iconClasses"
    aria-hidden="true"
  />
</template>

<script setup lang="ts">
import { type Component, computed, shallowRef, watch } from "vue";

const props = defineProps<{
  name: string;
  size?: string;
  color?: string;
  clickable?: boolean;
}>();

const SvgComponent = shallowRef<Component | null>(null);

watch(
  () => props.name,
  async (name) => {
    try {
      SvgComponent.value = (await import(`@/assets/icons/${name}.svg`)).default;
    } catch {
      console.warn(`Icon not found: ${name}`);
      SvgComponent.value = null;
    }
  },
  { immediate: true },
);

const iconClasses = computed(() => [
  props.size ?? "w-4 h-4",
  props.color ?? "text-current",
  props.clickable ? "cursor-pointer" : "",
  "inline-block",
]);
</script>
