<template>
  <v-dialog v-model="show"  width="90%" class="v-dialog">
    <v-layout v-if="show" wrap row style="background-color: white">
      <v-flex xs12 v-for="(item, index) in data.multiplechoice.Results" :key="index">
        <v-card style="overflow: hidden;">
          <div class="close">
            <v-icon size="30" @click="show=false">close</v-icon>
          </div>
          <div class="row">
            <div>
              <img
                :src="image(item)"
                class="column"
                style="margin: 10px, padding: 10px,float: left"
                width="338"
              />
            </div>
            <div class="column, title" style="width: 100%">
              <p style="margin:10px">{{item.Title}}</p>
              <p style="margin:10px">{{release(item)}}</p>
              <div class="row">
                <div class="column">
                  <v-btn outline dark @click="saveTMDBID( item)">
                    Add Metadata
                    <v-icon right dark>add_circle_outline</v-icon>
                  </v-btn>
                  <v-btn outline dark @click="openTmdb( item)">
                    View in TMDB
                    <v-icon right dark>movie</v-icon>
                  </v-btn>
                  <v-btn outline dark @click="playMovie( movie)">
                    Paly
                    <v-icon right dark>play_circle_outline</v-icon>
                  </v-btn>
                </div>
              </div>
            </div>
          </div>
        </v-card>
      </v-flex>
    </v-layout>
  </v-dialog>
</template>

<script>
import movieApi from '@/services/MovieApi'
// import Play from '@/components/buttons/Play'
export default {
  components: {

  },
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
        : (pic = this.$baseURL + '/html/images/na.gif')
      return pic
    },
    imdburl () {
      return 'http://www.imdb.com/title/' + this.movie.meta.imdb_id
    },
    title () {
      let txt = this.movie.title
      txt = txt.replace(/\(/g, '<br>(')
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
        .then(response => {
          this.loading = false
          if (response.meta.ID > 0) {
            this.movie = response
            this.$emit('changeMovie', this.movie)
            this.show = false
          }
        })

        .catch(error => {
          console.log(error)
        })
    },
    addMeta (item) {},
    playMovie (movie) {
      movieApi
        .playLocal(movie)
        .then(response => {})
        .catch(error => {
          console.log(error)
        })
    }
  }
}
</script>
<style scoped >
.title {
  background: #337ab7;
  color: white;
  margin: 20px;
}

.cast {
  cursor: pointer;
  width: 30%;
  float: left;
  color: #337ab7;
}

.is-collapsed {
  div:nth-child(n + 5) {
    display: none;
  }
}

.v-dialog {
  position: absolute;
  top: 20px;
  margin-left: auto;
  margin-right: auto;
  background-color: white;
}

.row {
  display: flex;
  flex-direction: row;
  margin-left: 10px;

}

.rowsmall {
  display: flex;
  flex-direction: row;
}

.column {
  display: flex;
  flex-direction: column;
  flex-basis: 100%;
  margin: 20px;
}

.coldetail {
  display: flex;
  flex-direction: column;
  width: 20%;
  background: #337ab7;
  color: white;
  margin-left: 10px;
  margin-block-end: 10px;
  padding-block-start: 8px;
  padding-left: 10px;
}

.crew {
  background: #337ab7;
  color: white;
  font-size: large;
  margin: 10px;
  text-align: center;
}

.crewlist {
  background: white;
}

.detail {
  background: #e0e0e0;
  display: flex;
  flex-direction: column;
  width: 40%;
  margin-left: 10px;
  margin-block-end: 10px;
  padding-block-start: 8px;
}

.close {
  position: absolute;
  right: 5px;
  top: 5px;
  background-color: white;
}
</style>
