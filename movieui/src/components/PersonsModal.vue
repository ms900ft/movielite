<template>
  <v-dialog v-model="show" overlay-color="black" overlay-opacity="1" style="position: absolute; top: 0;z-index:999">
    <v-carousel hide-delimiters v-model="startIndex" height="100vh" >
      <v-carousel-item v-for="(item, i) in items" :key="i">
        <v-container>
          <div class="title">
            {{ item.Name }} ({{desc(item)}})
          </div>
          <v-img
            :contain="contain"
            @click="close()"
            :src="image(item)"
            :max-height="maxheight"
          ></v-img>
        </v-container>
      </v-carousel-item>
    </v-carousel>
  </v-dialog>
</template>

<script>
export default {
  props: {
    value: Boolean,
    startIndex: Number,
    data: {}
  },
  mounted () {
    this.items = this.data
  },
  data () {
    return {
      // data: {},
      // movie: this.data,

      maxheight: '90vh',
      contain: true,
      items: {}
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
    image (item) {
      if (!item) {
        return ''
      }
      var pic
      var size = 500
      item.profile_path
        ? (pic = this.$baseURL + '/images/w' + size + item.profile_path)
        : (pic = this.$baseURL + '/movie2/na.gif')
      return pic
    },
    close (value) {
      this.$emit('input', value)
    },
    desc (item) {
      if (item.Job) {
        return item.Job
      } else {
        return item.Character
      }
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
.fullimagecontainer {
  margin-left: auto;
  margin-right: auto;
  position: relative;
  display: flex;
}
.title {
  color: aliceblue;
  width: 50%;
  margin: 0 auto;
  text-align: center;
}

</style>‚
