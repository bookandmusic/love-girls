// stores/ui.ts
import { defineStore } from 'pinia'

export const useUIStore = defineStore('ui', {
  state: () => ({
    loading: false,
    playing: true,
  }),
  actions: {
    setLoading(val: boolean) {
      this.loading = val
    },
    setPlaying(val: boolean) {
      this.playing = val
    },
  },
})
