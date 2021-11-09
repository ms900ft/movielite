<template>
  <v-container grid-list-xl fluid>
    <v-container v-show="loading">
      <div class="center-screen">
        <v-progress-circular
          indeterminate
          :size="150"
          :width="8"
          color="#0d47a1"
        ></v-progress-circular>
      </div>
    </v-container>

    <movieModal
      v-show="modalVisible"
      @close="modalVisible = false"
      :movie="modalData"
      :isLocal="isLocal()"
      fullscreen
      v-on:close="modalVisible = false"
      v-on:changeMovie="changeList"
      v-model="modalVisible"
    />

    <choiceModal
      v-show="choiceModalVisible"
      @close="choiceModalVisible = false"
      :data="modalData"
      v-on:changeMovie="changeList"
      v-model="choiceModalVisible"
    />
    <imageModal
      v-show="imageModalVisible"
      @close="imageModalVisible = false"
      :data="modalData"
      v-model="imageModalVisible"
    />
    <v-layout wrap row style="margin: auto 40px">
      <v-flex
        style="flex: 0"
        v-for="(item, index) in wholeResponse"
        :key="index"
      >
        <v-card width="180px" height="100%">
          <image-overview
            v-on:openModal="openModal(item)"
            v-on:imageModal="openImageModal(item)"
            v-on:openChoiceModal="openChoiceModal(item)"
            :movie="item"
            :index="index"
            :localHelper="playHelper"
          ></image-overview>
        </v-card>
      </v-flex>
      <infinite-loading @infinite="infiniteHandler">
        <div slot="no-more"></div>
        <div slot="no-results"></div>
        <span slot="spinner">
          <div class="center-screen">
            <v-progress-circular
              indeterminate
              :size="150"
              :width="8"
              color="#0d47a1"
            ></v-progress-circular>
          </div>
        </span>
        <!-- <template slot="spinner"></template> -->
      </infinite-loading>
    </v-layout>
  </v-container>
</template>

<script>
import movieApi from '@/services/MovieApi'
import localApi from '@/services/LocalApi'
import movieModal from '@/components/MovieModal'
import choiceModal from '@/components/ChoiceModal'
import imageModal from '@/components/ImageModal'
// import BurgerMenu from '@/components/menu/Burger'
import ImageOverview from '@/components/images/Overview'
// import VueCookies from 'vue-cookies'
import InfiniteLoading from 'vue-infinite-loading'
import EventBus from '@/store/eventBus.js'

export default {
  components: {
    movieModal,
    imageModal,
    choiceModal,
    InfiniteLoading,
    ImageOverview
  },

  data () {
    return {
      wholeResponse: [],
      loading: true,
      page: 1,
      modalVisible: false,
      choiceModalVisible: false,
      imageModalVisible: false,
      modalData: '',
      showByIndex: null,
      // searchstring: this.$store.state.searchString,
      truncateDesc: 250,
      playHelper: false
    }
  },
  props: [
    'show',
    'orderby',
    'genre',
    'country',
    'cast',
    'crew',
    'person'
    // "searchstring",
  ],
  mounted () {
    this.page = 1
    this.getMovies()
    console.log(this.$vuetify.breakpoint)

    EventBus.$on('DELETEMOVIE', item => {
      this.deleteMovie(item)
    })
    EventBus.$on('TOGGLEWATCHLIST', item => {
      this.toggleWatchlist(item)
    })
    EventBus.$on('ISTVSHOW', item => {
      this.isTvShow(item)
    })
    EventBus.$on('RESCAN', item => {
      console.log('------------------------------------')
      console.log('rescanxxx')
      console.log(item)
      this.rescanMovie(item)
    })
    EventBus.$on('FULLIMAGE', item => {
      this.imageModalVisible = !this.imageModalVisible
    })

    this.localHelper()
    // this.descLength();
    // this.searchboxVisible = false;
  },
  computed: {
    searchstring () {
      return this.$store.state.searchString
    },
    mobile () {
      return this.$isMobile
    }
  },
  methods: {
    fontSize () {
      let ret = 14
      if (this.$vuetify.breakpoint.name === 'xs') {
        return ret
      }
      switch (this.size) {
        case '6':
          ret = 24
          break
        case '12':
          ret = 28
          break
        case '4':
          ret = 16
          break
        case '3':
          ret = 16
          break
      }
      return ret
    },

    descLength () {
      let ret = 0
      if (this.$vuetify.breakpoint.name === 'xs') {
        return ret
      }

      return ret
    },
    changeList (item) {
      const index = this.wholeResponse.map(movie => movie.id).indexOf(item.id)

      this.wholeResponse.splice(index, 1, item)
      this.modalData = item
    },
    singleMovie (id) {
      this.$router.push('/movie/' + id)
    },

    desc (item) {
      var desc
      item.meta ? (desc = item.meta.Overview) : (desc = '')
      return desc
    },
    title (item) {
      let txt = item.title
      // txt = txt.replace(/\(/g, "<br>(");
      return txt
    },
    release (item) {
      if (item.meta != null) {
        const date = new Date(item.meta.release_date)
        return '(' + (1900 + date.getYear()) + ')'
      }
      return ''
    },
    prevPage () {
      if (this.page === 1) return
      --this.page
      this.getMovies()
    },
    nextPage () {
      ++this.page
      this.getMovies()
    },
    infiniteHandler ($state) {
      ++this.page

      if (this.total + this.$hitspp > this.page * this.$hitspp) {
        this.getMovies($state)
      } else {
        $state.complete()
      }

      // $state.loaded;
    },
    getMovies ($state) {
      this.loading = true
      if (!$state) {
        this.$store.state.resultsFound = 0
      }
      let args = {
        page: this.page,
        show: this.show,
        orderby: this.orderby,
        genre: this.genre,
        title: this.searchstring,
        country: this.country,
        cast: this.cast,
        crew: this.crew,
        person: this.person
      }
      movieApi
        .fetchMovieCollection(this, args)
        .then(response => {
          this.wholeResponse.push(...response.data)
          this.$store.commit('setResultsFound', response.meta.total)
          this.total = response.meta.total
          this.loading = false
          if ($state) {
            if (!response.data.length) {
              $state.complete()
            }
            $state.loaded()
          }
        })
        .catch(error => {
          console.log(error)
        })
      if ($state) {
        // $state.complete();
      }
    },
    openModal (item) {
      this.modalVisible = true
      this.modalData = item
    },
    openChoiceModal (item) {
      this.modalData = item
      this.choiceModalVisible = true
    },
    closeModal () {
      this.modalVisible = false
    },
    openImageModal (item) {
      this.imageModalVisible = true
      this.modalData = item
    },
    changed (value) {
      console.log('Value ' + value)
    },

    custom (keyboard) {
      console.log(keyboard.value)
    },
    toggleWatchlist (item) {
      this.updateMovie(item, { watchlist: true })
    },
    isTvShow (item) {
      this.updateMovie(item, { isTvShow: true })
    },
    updateMovie (item, args) {
      movieApi
        .updateMovie(item)
        .then(response => {
          this.loading = false
        })
        .catch(error => {
          console.log(error)
        })
      if (args.isTvShow) {
        this.wholeResponse.splice(this.wholeResponse.indexOf(item), 1)
      }
      if (args.watchlist && this.show === 'watchlist') {
        this.wholeResponse.splice(this.wholeResponse.indexOf(item), 1)
      }
    },

    deleteMovie (item) {
      movieApi
        .deleteMovie(item)
        .then(response => {
          this.loading = false
        })
        .catch(error => {
          console.log(error)
        })
      this.$store.commit('setResultsFound', this.total - 1)
      this.total = this.total - 1
      this.wholeResponse.splice(this.wholeResponse.indexOf(item), 1)
    },
    multipleChoice (item) {
      if (item.multiplechoice != null && item.meta == null) {
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
    localHelper () {
      localApi
        .ping(this.$localViewURL)
        .then(response => {
          if (response.data.alive === 'ok') {
            this.playHelper = true
          } else {
            this.playHelper = false
          }
        })
        .catch(error => {
          console.log(error)
        })
    },
    rescanMovie (item) {
      const id = item.meta.ID
      // item.watchlist = false
      movieApi
        .addMeta(item, id)
        .then(response => {
          this.loading = false
          if (response.meta.ID > 0) {
            this.movie = response
            // this.$emit("changeMovie", this.movie);
            // this.show = false;
          }
        })
        .catch(error => {
          console.log(error)
        })
    }
  },
  watch: {
    '$route.query' () {
      this.page = 1
      this.wholeResponse = []
      this.getMovies()
    },

    searchstring () {
      if (this.timer) {
        clearTimeout(this.timer)
        this.timer = null
      }
      this.timer = setTimeout(() => {
        this.page = 1
        this.wholeResponse = []
        this.getMovies() // your code
      }, 500)
    }
  }
}
</script>

<style lang="stylus" scoped>
.v-progress-circular {
  margin: 1rem;
}

.movieaction {
  position: absolute;
  bottom: 5px;
  right: 10px;
  display: flex;
  flex-direction: row;
  margin-left: 10px;
}

.moviesmall {
  position: absolute;
  top: 15px;
  left: 180px;
}

.centered {
  position: fixed; /* or absolute */
  top: 50%;
  left: 50%;
}

.moviedesc {
  padding-right: 20px;
}

.center-screen {
position: fixed;
top: 50%;
left: 50%;
margin-top: -50px;
margin-left: -50px;
width: 100px;
height: 100px;
}

</style>
