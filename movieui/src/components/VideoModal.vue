<template>
  <v-dialog v-model="show"  style="position: relative">
    dkdskdsdksdskdsk
   <vue-player :src="videoURL()"
   :poster="imageURL(data)"
   title="this is a title"
   preview-on-mouse = true
   controls = true
   class="vplayer"
   v-model="playing">

   </vue-player>

  </v-dialog>
</template>

<script>
import vuePlayer from '@algoz098/vue-player'

export default {
  components: {
    vuePlayer
  },
  props: {
    value: Boolean,
    data: {},
    playing: Boolean
  },
  computed: {
    show: {
      get () {
        return this.value
      },
      set (value) {
        this.$emit('input', value)
      },
      videoURL () {
        console.log(this.$baseURL + '/file/' + this.movie.file_id + '/download')

        return this.$baseURL + '/file/' + this.movie.file_id + '/download'
      // http://192.168.1.4:8000/file/161021/download
      }
    }
  },
  methods: {
    imageURL () {
      var pic
      this.data && this.data.meta.backdrop_path
        ? (pic = this.$baseURL + '/images/w500' + this.data.meta.backdrop_path)
        : (pic = this.$baseURL + '/html/images/na.gif')
      console.log('------------------------------------')
      console.log(pic)
      console.log('------------------------------------')
      return pic
    },
    close (value) {
      this.$emit('input', value)
    },
    videoURL () {
      console.log(this.$baseURL + '/file/' + this.data.file_id + '/download')

      return this.$baseURL + '/file/' + this.data.file_id + '/download'
      // http://192.168.1.4:8000/file/161021/download
    }
    // playing () {
    //   console.log('------------------------------------')
    //   console.log('palying')
    //   console.log('------------------------------------')
    // }
  }
}
</script>

<style scoped>
.close {
  position: absolute;
  right: 10px;
  top: 10px;
  opacity: 1;
  background-color: white;
}
.vplayer {
  width: 500px;
}
</style>‚
