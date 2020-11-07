<template>
  <v-dialog v-model="show" height="100%" fullscreen style="position: relative" content-class="imagedialog">
    <img :src="image()" @click="close()" class="fullimage">
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

<style >
.close {
  position: absolute;
  right: 10px;
  top: 10px;
  opacity: 0.5;
  background-color: white;
}
.imagedialog {
  width: 100%;
  position: relative;
  max-height: 100%;
  top: 0px;
    /* margin-left: -440px;
  left: 50%; */
  display: grid;

}
.fullimage {
  padding: 3px;
  background-color: #bf4141;
  margin-left: auto;
  margin-right: auto;

  /* max-width:100%;
max-height:100%;
object-fit: contain */
  max-width: 100%;
  max-height: 100vh;
  margin: auto;
}
</style>‚
