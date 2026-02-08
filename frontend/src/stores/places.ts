import { defineStore } from 'pinia'

import { type Place, placeApi } from '@/services/placeApi'

export const usePlacesStore = defineStore('places', {
  state: () => ({
    places: [] as Place[],
  }),

  getters: {
    getPlaces(): Place[] {
      return this.places
    },
  },

  actions: {
    async fetchPlaces() {
      const res = await placeApi.getPlaces(1, 99999) // 获取所有地点
      if (res.code === 0) {
        this.places = res.data.places
      }
    },
  },
})
