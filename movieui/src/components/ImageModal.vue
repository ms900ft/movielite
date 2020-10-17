<template>
  <v-dialog v-model="show" width="800px" style="position: relative">
    <img :src="image()" @click="close()" style="padding: 10px;background-color: white;">
    <div class="close">
      <v-icon size="30" @click="show=false">close</v-icon>
    </div>
  </v-dialog>
</template>

<script>
export default {
  props: {
    value: Boolean,
    data: {}
  },
  data () {
    return {
      // data: {},
      // movie: this.data,

      fullimage: false
      // TMDBID: this.data.meta.ID
    }
  },
  computed: {
    show: {
      get () {
        return this.value
      },
      set (value) {
        this.$emit('input', value)
      }
    }
  },
  methods: {
    image () {
      var pic
      this.data.meta && this.data.meta.poster_path
        ? (pic = this.$baseURL + '/images/w780' + this.data.meta.poster_path)
        : (pic = this.$baseURL + '/movie2/nocover2.jpg')
      return pic
    },
    close (value) {
      this.$emit('input', value)
    }
  }
}
</script>

<style scoped>
.close {
  position: absolute;
  right: 10px;
  top: 10px;
  opacity: 0.5;
  background-color: white;
}
</style>‚
