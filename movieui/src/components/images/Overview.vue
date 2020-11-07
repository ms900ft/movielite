<template>
  <div
    v-bind:style="styleObject"
    @mouseover="showByIndex = index"
    @mouseout="showByIndex = null"
  >
    <div class="imagecontainer">
      <img
        :src="image(movie)"
        class="image"
        v-bind:style="{ Width: maxWidth(), height: height }"
        :class="{ selected: showByIndex === index }"
        @click="openModal(movie)"
      />
      <div v-show="showButton(index)" class="burger">
        <burger-menu :movie="movie"></burger-menu>
      </div>
      <div class="middle" v-show="showButton(index)">
        <play
          v-if="isLocal() || localHelper"
          :movie="movie"
          overlay="true"
          :helper="localHelper"
        />
        <download v-else :movie="movie" overlay="true" />
      </div>
      <div v-if="multipleChoice(movie)" class="multibutton">
        <v-btn color="#0d47a1" dark @click="openChoiceModal(movie)">
          Find
          <v-icon right dark>search</v-icon>
        </v-btn>
      </div>

      <div
        v-show="imageWithTitle(movie) || mobile || showByIndex === index"
        class="imagetitle"
        :class="{ selected: showByIndex === index }"
      >
        {{ movie.title }}
      </div>
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
    } else {
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
        background: background,
        width: '180px'
      }
    }
  },
  props: ['movie', 'index', 'localHelper'],
  mounted () {},
  computed: {
    mobile () {
      return this.$isMobile
    }
  },
  methods: {
    image (item) {
      var pic

      item.meta && item.meta.poster_path
        ? (pic = this.$baseURL + '/images/w342' + item.meta.poster_path)
        : (pic = this.$baseURL + '/movie2/nocover.jpg')
      return pic
    },
    openModal (item) {
      this.$emit('openModal', this.item)
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
    maxWidth () {
      if (this.mobile === true) {
        return '100%'
      }
      return '180px'
    },
    multipleChoice (item) {
      if (item.multiplechoice != null && item.meta == null) {
        return true
      }
      return false
    },
    openChoiceModal (item) {
      this.$emit('openChoiceModal', item)
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
  /* object-fit: contain; */
  width: 180px;
  max-height: 240px;
  /* width: auto; */
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
.multibutton {
  /* transition: 0.5s ease;
  opacity: 1; */
  position: absolute;
  bottom: 10%;
  left: 50%;
  transform: translate(-50%, -50%);
  -ms-transform: translate(-50%, -50%);
  text-align: center;
}
.imagecontainer {
  position: relative;
  overflow: hidden;
  width: 180px;
  box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}
</style>
