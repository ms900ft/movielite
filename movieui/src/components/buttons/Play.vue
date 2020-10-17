<template>
  <div>
    <v-btn v-if="!overlay" outline color="primary" dark @click="playMovie(movie)">
      Play
      <v-icon right dark>play_circle_outline</v-icon>
    </v-btn>
    <v-icon v-else style="color: grey" @click="playMovie(movie)" left size="60">play_circle_outline</v-icon>
  </div>
</template>

<script>
import movieApi from '@/services/MovieApi'
import localApi from '@/services/LocalApi'

export default {
  name: 'Play',
  components: {},
  props: ['movie', 'size', 'helper', 'overlay'],
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
    playMovie (movie) {
      if (this.helper) {
        localApi
          .play(this.$localViewURL, movie)
          .then(response => {})
          .catch(error => {
            console.log(error)
          })
      } else {
        movieApi
          .playLocal(movie)
          .then(response => {})
          .catch(error => {
            console.log(error)
          })
      }
    },
    small () {
      if (this.size < 3) {
        return true
      }
      return false
    }
  }
}
</script>

<style lang="stylus" scoped>
.hamburger {
  float: right;
}
</style>
