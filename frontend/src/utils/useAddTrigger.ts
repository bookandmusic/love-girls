// composables/useAddTrigger.ts
import { nextTick, ref } from 'vue'

export const useAddTrigger = () => {
  const trigger = ref(false)

  const fire = () => {
    trigger.value = true
    nextTick(() => (trigger.value = false))
  }

  return { trigger, fire }
}
