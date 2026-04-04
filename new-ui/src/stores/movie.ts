import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useMovieStore = defineStore('movie', () => {
  const totalResults = ref(0)

  function setTotalResults(total: number) {
    totalResults.value = total
  }

  return { totalResults, setTotalResults }
})
