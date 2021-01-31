<template>
  <v-dialog v-model="show" width="90%" >
    <v-layout v-if="show"  style="background-color: white;">
        <div class="close">
            <v-icon size="30" outlined @click="show = false">mdi-close</v-icon>
          </div>

       <v-container  >
      <v-row class="mr-6" dense >
        <v-col
          v-for="(item, index) in data.multiplechoice.Results"
          :key="index"
          cols="12"
          class="ma-6"

        >
          <v-card
            color="#337ab7"
            dark
            class="d-flex flex-column"

          >
            <div class="d-flex flex-no-wrap justify-space-between">
              <v-container >
                <v-card-title
                  class="headline"
                  v-text="title(item) + release(item) "
                ></v-card-title>

                <v-card-text  v-text="item.overview"></v-card-text>
  <v-card-text ></v-card-text >
                <v-card-actions style="position:absolute; bottom:0px;">

                 <v-btn    @click="saveTMDBID(item)">
                    Add Metadata
                    <v-icon right >add_circle_outline</v-icon>
                  </v-btn>
                  <v-btn   @click="openTmdb(item)">
                    View in TMDB
                    <v-icon right >movie</v-icon>
                  </v-btn>
                  <v-btn   @click="playMovie(movie)">
                    Paly
                    <v-icon right >play_circle_outline</v-icon>
                  </v-btn>
                </v-card-actions>
              </v-container>

              <v-avatar
                class="ma-3"
                height="400px"
                width="400px"
                tile
              >
                <v-img contain :src="image(item)"></v-img>
              </v-avatar>
            </div>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
      <!-- <v-flex
        xs12
        v-for="(item, index) in data.multiplechoice.Results"
        :key="index"
      >
        <v-card style="overflow: hidden">
          <div class="close">
            <v-icon size="30" @click="show = false">mdi-close</v-icon>
          </div>
          <div class="choice-row">
            <div >
              <img
                :src="image(item)"
                class="column"
                style="margin: 10px, padding: 10px,float: left"
                width="338"
              />
            </div>
            <div class="column, title" style="width: 100%">
              <v-row>
              <p style="margin: 10px">{{ title(item) }} {{ release(item) }}</p>
              <p style="margin: 10px" class="moviedesc">{{ item.overview }}</p>

              </v-row>

  <v-container fill-height>
    <v-row

    align="center"
    justify="space-around"
    style="margin-top: 50px"

  >
                  <v-btn    @click="saveTMDBID(item)">
                    Add Metadata
                    <v-icon right >add_circle_outline</v-icon>
                  </v-btn>
                  <v-btn   @click="openTmdb(item)">
                    View in TMDB
                    <v-icon right >movie</v-icon>
                  </v-btn>
                  <v-btn   @click="playMovie(movie)">
                    Paly
                    <v-icon right >play_circle_outline</v-icon>
                  </v-btn>
                  </v-row>
  </v-container>

                </div>
            </div>

        </v-card>
      </v-flex> -->
    </v-layout>
  </v-dialog>
</template>

<script>
import movieApi from '@/services/MovieApi'
// import Play from '@/components/buttons/Play'
export default {
  components: {},
  data () {
    return {
      movie: this.data

      // TMDBID: this.data.meta.ID
    }
  },

  props: {
    value: Boolean,
    data: {}
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
      var pic
      item.poster_path
        ? (pic = this.$baseURL + '/images/w500' + item.poster_path)
        : (pic = this.$baseURL + '/movie2/nocover.jpg')
      return pic
    },
    imdburl () {
      return 'http://www.imdb.com/title/' + this.movie.meta.imdb_id
    },
    title (item) {
      let txt = item.title
      // txt = txt.replace(/\(/g, '<br>(')
      return txt
    },
    release (item) {
      if (item != null) {
        const date = new Date(item.release_date)
        return '(' + (1900 + date.getYear()) + ')'
      }
      return ''
    },
    openTmdb (item) {
      const id = item.id ? item.id : item.ID
      window.open('https://www.themoviedb.org/movie/' + id, '_blank')
    },
    saveTMDBID (item) {
      const id = item.id ? item.id : item.ID
      this.data.watchlist = true
      movieApi
        .addMeta(this.data, id)
        .then((response) => {
          this.loading = false
          if (response.meta.ID > 0) {
            this.movie = response
            this.$emit('changeMovie', this.movie)
            this.show = false
          }
        })

        .catch((error) => {
          console.log(error)
        })
    },
    addMeta (item) {},
    playMovie (movie) {
      movieApi
        .playLocal(movie)
        .then((response) => {})
        .catch((error) => {
          console.log(error)
        })
    }
  }
}
</script>
<style scoped >

.close {
  position: absolute;
  right: 0px;
  top: 0px;
  background-color: white;
}
</style>
