<template>
  <div>
    <v-icon  v-if="movie.watchlist" style="color: orange" @click="toggleWatchlist(movie)" left size="20">star</v-icon>

      <!-- <v-icon  right dark>star_outline</v-icon> -->
    <v-icon  v-else style="color: orange" @click="toggleWatchlist(movie)" left size="20">star_outline</v-icon>
  </div>
</template>

<script>
import EventBus from '@/store/eventBus.js'

export default {
  name: 'Watchlist',
  components: {},
  props: ['movie'],
  data () {
    return {
      items: [],
      x: false
    }
  },
  mounted () {
    // this.descLength();
    // this.searchboxVisible = false;
    //    this.localHelper = this.Helper();
  },
  methods: {
    toggleWatchlist (item) {
      item.watchlist = !item.watchlist
      EventBus.$emit('TOGGLEWATCHLIST', item)
    }
  }
}
</script>

<style lang="stylus" scoped>
.hamburger {
  float: right;
}
</style>
