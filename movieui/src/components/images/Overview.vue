<template>
  <div v-bind:style="styleObject" @mouseover="showByIndex = index" @mouseout="showByIndex = null">
    <div class="imagecontainer">

      <img

        :src="image(movie)"

        class="image"
        v-bind:style="{maxWidth:maxWidth(), height:height  }"
        :class="{'selected': showByIndex === index}"
        @click="openImageModal(movie)"
      />
      <div v-show="showButton(index) && size <3" class="burger">
        <burger-menu :movie="movie"></burger-menu>
      </div>
      <div v-show="showButton(index) && size <3" class="burger">
        <burger-menu :movie="movie"></burger-menu>
      </div>
      <div class="middle" v-show="showButton(index)">
        <play
          v-if="isLocal() || localHelper"
          :movie="movie"
          :size="size"
          overlay="true"
          :helper="localHelper"
        />
        <download v-else :movie="movie" overlay="true" :size="size" />
      </div>
      <div v-show="(showButton(index) && size <3) || imageWithTitle(movie) || mobile"
      class="imagetitle"
      :class="{'selected': showByIndex === index}" >{{movie.title}}</div>
    </div>
  </div>
</template>

<script>
import Play from '@/components/buttons/Play'

import Download from '@/components/buttons/Download'
import BurgerMenu from '@/components/menu/Burger'

export default {
  name: 'ImageOverview',
  components: { Play, Download, BurgerMenu },

  data () {
    let float = ''
    let background = 'white'
    if (this.$isMobile === true) {
      float = 'auto'
      background = 'radial-gradient(transparent 30%, black 50%)'
    } else if (this.size > 2) {
      float = 'left'
    }
    return {
      showByIndex: null,
      playHelper: false,
      height: '484px',
      // xxx flot left to show text
      styleObject: {
        float: float,
        height: '100%',
        'margin-left': 'auto',
        'margin-right': 'auto',
        'background': background
      }
    }
  },
  props: ['movie', 'size', 'index', 'localHelper'],
  mounted () {},
  computed: {
    mobile () {
      return this.$isMobile
    }
  },
  methods: {
    image (item) {
      var pic
      if (this.mobile === true) {
        if (item.meta && item.meta.backdrop_path) {
          pic = this.$baseURL + '/images/w300' + item.meta.backdrop_path
          this.height = 'auto'
        } else if (item.meta && item.meta.poster_path) {
          pic = this.$baseURL + '/images/w342' + item.meta.poster_path
          this.height = '484px'
        } else {
          pic = this.$baseURL + '/movie2/nocover.jpg'
          this.height = '484px'
        }
        return pic
        //   ? (pic = this.$baseURL + '/images/w500' + item.meta.backdrop_path)
        //   : (pic = this.$baseURL + '/movie2/nocover.jpg')
        // return pic
      }
      item.meta && item.meta.poster_path
        ? (pic = this.$baseURL + '/images/w342' + item.meta.poster_path)
        : (pic = this.$baseURL + '/movie2/nocover.jpg')
      return pic
    },
    openImageModal (item) {
      if (this.size < 4 || this.$vuetify.breakpoint.name === 'xs') {
        this.$emit('openModal', this.item)
      } else {
        this.$emit('imageModal', this.item)
      }
    },
    imageWithTitle (item) {
      if (!item.meta) {
        return true
      }
      if (item.meta && !item.meta.poster_path) {
        return true
      }
      return false
    },
    isLocal () {
      if (location.hostname === 'localhost') {
        return true
      }
      return false
    },

    showButton (index) {
      if (this.mobile) {
        return true
      }
      if (this.showByIndex === index) {
        return true
      }
      return false
    },
    sizex () {
      return this.size
    },
    maxWidth () {
      if (this.mobile === true) {
        return '100%'
      }
      return '180px'
    }
  }
}
</script>

<style scoped>
.text-block {
  bottom: 0px;
  top: 20px;
  left: 0px;
  margin-left: 0px;
  font-size: large;
  color: #337ab7;
  background-color: white;
  padding-left: 10px;
  padding-right: 0px;
  padding-top: 30px;

  padding: 10px;
  position: inherit;
  /* height: 250px;
  width: 180px; */
}
.imagetitle {
  bottom: 0px;
  margin: 0, auto;
  font-size: large;
  color: white;
  opacity: 1;
  background-color: #0d47a1;
  padding-left: 10px;
transition: 1.5s ease;

  padding-top: 30px;
  width: 100%;
  max-height: 50%;
  padding: 10px;
  position: absolute;

}
.imagetitle.selected {

  opacity: 0.9;

}
.image {
  opacity: 1;
  display: block;
  transition: 1.5s ease;
  backface-visibility: hidden;
  margin: auto;
  margin-left: auto;
  margin-right: auto;
  width: 100%;
  object-fit: contain;
  max-width: 100%;
  max-height: 240px;
  width: auto;
  height: auto;
  transform: scale(1.08, 1);
}

.image.selected {
  opacity: 0.6;
  transform: scale(1.2, 1.2);

}

.middle {
  transition: 0.5s ease;
  opacity: 1;
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  -ms-transform: translate(-50%, -50%);
  text-align: center;
}

.burger {
  transition: 0.5s ease;
  opacity: 1;
  position: absolute;
  top: 10px;
  right: 10px;
  transform: translate(-50%, -50%);
  -ms-transform: translate(-50%, -50%);
  text-align: center;
}
.imagecontainer {
  position: relative;
  overflow: hidden;

}
</style>
