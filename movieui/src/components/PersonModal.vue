<template>
  <v-dialog v-model="show"  overlay-color="black"
   overlay-opacity="1">

    <v-container>
  <div class="title">
    {{data.Name}}
  </div>
    <v-img
  :contain="contain"
   @click="close()"
  :src="image()"
   :max-height="maxheight"

></v-img>
    </v-container>
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

      fullimage: false,
      maxheight: '90vh',
      contain: true
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
      var size = 500
      this.data.profile_path
        ? (pic = this.$baseURL + '/images/w' + size + this.data.profile_path)
        : (pic = this.$baseURL + '/movie2/na.gif')
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
.fullimagecontainer{
 margin-left: auto;
  margin-right: auto;
  position: relative;
  display: flex;
}
.title{
  color: aliceblue;
  width: 50%;
  margin: 0 auto;
  text-align: center;
}
</style>‚
