<template>
  <v-dialog  v-model="show"  width="100%">

    <div  v-if="show" class="media">
      <div class="close">
        <v-icon size="30" @click="close()">close</v-icon>
      </div>
      <div class="media-object">
        <div>
          <img
            :src="image()"
            @click=showFullImage()
            height="250px"
            style="padding: 10px"
            width="180px"
          />
        </div>
        <div>
          <play v-if="isLocal" :movie="movie" />
        </div>
      </div>
      <div class="media-body">
        <div class="media-heading">
          <span style="margin-left: 10px;" class="movietitle" v-html="title()">{{title()}}</span>
          <span class="release">{{release()}}</span>
        </div>
        <p >
          <span class="titledetails">{{titledetails()}}</span>
        </p>
        <div
          v-if="movie.meta"
          class="moviedesc moviedescmodal"
          style="padding-left: 0px;"
        >{{movie.meta.Overview }}</div>
        <div class="rowsmall" style="margin-top: 10">
          <div class="coldetail">Scaned Title</div>
          <div class="detail">
            <editable
              v-bind:value="movie.title"
              class-name="text-editable"
              :is-edit="false"
              :on-change="saveTitle"
              :multi="false"
            ></editable>
          </div>
        </div>
        <div class="rowsmall" style="margin-top: 10">
          <div class="coldetail">TMDB ID</div>
          <div class="detail">
            <editable
              v-bind:value="TMDBID"
              class-name="text-editable"
              :is-edit="false"
              :on-change="saveTMDBID"
              :multi="false"
            ></editable>
          </div>
        </div>
        <div class="rowsmall" style="margin-top: 10" v-if="movie.meta">
          <div class="coldetail">Original Title</div>
          <div class="detail">{{movie.meta.original_title}}</div>
        </div>
        <div class="rowsmall" style="margin-top: 10" v-if="movie.meta">
          <div class="coldetail">Production Contries</div>
          <div class="detail">{{countries()}}</div>
        </div>
        <div class="rowsmall" style="margin-top: 10" v-if="movie.meta">
          <div class="coldetail">Original Language</div>
          <div class="detail">{{movie.meta.original_language}}</div>
        </div>
        <div class="rowsmall" style="margin-top: 10" v-if="movie.meta">
          <div class="coldetail">IMDB</div>
          <div class="detail">
            <a v-bind:href="imdburl()" target="_new" style="text-decoration: none;">
              <v-icon>play_circle_outline</v-icon>
            </a>
          </div>
        </div>
        <div class="rowsmall" style="margin-top: 10">
          <div class="coldetail">Filename</div>
          <div class="detail breakall">{{filename()}}</div>
        </div>
        <div class="rowsmall" style="margin-top: 10">
          <div class="coldetail">Path</div>
          <div class="detail breakall">{{fullpath()}}</div>
        </div>
        <div v-if="movie.meta">
          <div class="row">
            <div class="column  crew">Cast</div>

            <div class="column crew crewleft">Crew</div>
          </div>
          <div v-if="movie.meta.Credits.Crew" class="row">
            <person-list :persons="movie.meta.Credits.Cast" type="cast" v-on:closeModal="close()"></person-list>
            <person-list class="crewleft" :persons="movie.meta.Credits.Crew" type="crew" v-on:closeModal="close()"></person-list>
          </div>
        </div>
      </div>
    </div>
  </v-dialog>
</template>

<script>
import editable from '@/components/fields/Editable'
import movieApi from '@/services/MovieApi'
import PersonList from '@/components/list/PersonList'
import Play from '@/components/buttons/Play'
import EventBus from '@/store/eventBus.js'
export default {
  components: {
    editable,
    PersonList,
    Play
  },
  data () {
    return {
      // data: {},
      // movie: this.data,
      // item: this.movie,
      fullimage: false
      // TMDBID: this.data.meta.ID
    }
  },

  props: {
    value: Boolean,
    movie: {},
    isLocal: Boolean
  },
  computed: {
    show: {
      get () {
        return this.value
      },
      set (value) {
        this.$emit('input', value)
      }
    },
    TMDBID: {
      get () {
        if (this.movie.meta) {
          return this.movie.meta.ID
        } else {
          return 0
        }
      },
      set (value) {
        this.movie.meta.ID = value
      }
    }
  },
  methods: {
    image () {
      var pic
      this.movie.meta && this.movie.meta.poster_path
        ? (pic = this.$baseURL + '/images/w342' + this.movie.meta.poster_path)
        : (pic = this.$baseURL + '/movie2/nocover.jpg')
      return pic
    },
    imdburl () {
      return 'http://www.imdb.com/title/' + this.movie.meta.imdb_id
    },
    filename () {
      if (this.movie !== '') {
        return this.movie.File.FileName
      }
      return ''
    },
    fullpath () {
      if (this.movie !== '') {
        return this.movie.File.FullPath
      }
      return ''
    },

    title () {
      let txt = this.movie.title || ''
      txt = txt.replace(/\(/g, '<br>(')
      return txt
    },
    release () {
      if (this.movie.meta != null) {
        const date = new Date(this.movie.meta.release_date)
        return '(' + (1900 + date.getYear()) + ')'
      }
      return ''
    },
    titledetails () {
      if (this.movie.meta != null) {
        let parts = []
        parts.push(this.movie.meta.Runtime + ' min')
        const genres = this.movie.meta.Genres.map(value => value.Name).join(
          ', '
        )
        parts.push(genres)

        const result = parts.join(' | ')
        return result
      }
      return ''
    },
    countries () {
      if (this.movie.meta != null) {
        let parts = []
        const countries = this.movie.meta.production_countries
          .map(value => value.Name)
          .join(', ')
        parts.push(countries)

        const result = parts.join(' | ')
        return result
      }
      return ''
    },
    saveTitle (title) {
      this.movie.title = title
      var item = this.movie
      this.updateMovie(item)
    },

    saveTMDBID (id) {
      let movie = this.movie
      movieApi
        .addMeta(movie, id)
        .then(response => {
          this.loading = false
          if (response.meta.ID > 0) {
            movie = response
            this.$emit('changeMovie', movie)
          }
        })
        .catch(error => {
          console.log(error)
        })
    },
    updateMovie (movie) {
      movieApi
        .updateMovie(movie)
        .then(response => {
          this.loading = false
          // this.movie = response.data
          this.$emit('changeMovie', movie)
          // }
        })
        .catch(error => {
          console.log(error)
        })
    },
    close () {
      this.$emit('close', false)
      // this.show = false
    },
    fullImage () {
      var pic
      this.movie.meta && this.movie.meta.poster_path
        ? (pic = this.$baseURL + '/images/w780' + this.movie.meta.poster_path)
        : (pic = this.$baseURL + '/movie2/nocover2.jpg')
      return pic
    },
    // close (value) {
    //   this.$emit('input', value)
    // },
    showFullImage () {
      EventBus.$emit('FULLIMAGE', 1)
    }

  }
}
</script>
<style lang="stylus" >
.titledetails {
  background: #337ab7;
  color: white;

}

.cast {
  cursor: pointer;
  width: 30%;
  float: left;
  color: #337ab7;
  font-size: large;
}

.is-collapsed {
  div:nth-child(n+5) {
    display: none;
  }
}

.v-dialog {
  position: absolute;
  top: 20px;
  margin-left: auto;
  margin-right: auto;
}

.row {
  display: flex;
  flex-direction: row;
}

.rowsmall {
  display: flex;
  flex-direction: row;
}

.column {
  display: flex;
  flex-direction: column;
  flex-basis: 100%;
}

.coldetail {
  display: flex;
  flex-direction: column;
  width: 30%;
  background: #337ab7;
  color: white;
  margin-block-end: 10px;
  padding-block-start: 8px;
  font-size: large;
}

.crew {
  background: #337ab7;
  color: white;
  font-size: large;
  text-align: center;

}
.crewleft{
  margin-left : 10px;
}
.crewlist {
  background: white;
}

.detail {
  background: #e0e0e0;
  display: flex;
  flex-direction: column;
  width: 100%;
  margin-left: 10px;
  margin-block-end: 10px;
  padding-block-start: 8px;
  font-size: large;
}

.breakall {
  word-break: break-all;
}

.close {
   position: absolute;
  right: 5px;
  top: 5px;
}

.flex-container {
  display: flex;
  flex-wrap: wrap;
  flex-direction: column;
  background-color: white;
  align-items: flex-start;
}

.media {
  position: relative;
  display: flex;
  align-items: flex-start;
  background: white;
  padding: 1em;
  border-radius: 3px;
  flex-wrap: wrap;
  max-width : 1200px;
  margin-left :auto;
  margin-right : auto;
}

.media-object {
  margin-right: 1em;
}

.media-body {
  flex: 1;
}

.media-heading {
  margin: 0 0 0.5em;
}

.fullimage {
  margin-left: auto;
  margin-right: auto;
}

.moviedescmodal{
  padding-left: 0px;
  padding-bottom : 10px;
}
</style>
