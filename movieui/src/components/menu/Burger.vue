
<template>
  <v-dialog v-if="dialog" v-model="dialog" persistent width="500px">
    <v-card>
      <v-card-title class="headline">
        Really delete:

      </v-card-title>
       <v-card-text>
<span class="movietitle">{{movie.title}}</span>
       </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="primary"  @click="dialog = false">Cancel</v-btn>
        <v-btn color="primary"  @click="deleteMovie(getMovieItem)">Delete</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <div v-else>
    <dir-list v-if="showdirlist" :show="showdirlist" :movie="movie" @close="showdirlist = false"></dir-list>
    <v-menu offset-y >
      <template v-slot:activator="{ on }">
        <div class="hamburger" v-on="on">
          <v-icon>menu</v-icon>
        </div>
      </template>
      <v-list>
        <v-list-item v-for="(item, index) in getItems" :key="index">
          <v-list-item-title
            @click="item.action(getMovieItem)"
            style="cursor: pointer"
          >{{ item.name() }}</v-list-item-title>
        </v-list-item>
        <v-list-item>
          <a style="color: black" :href="downloadURL()">Download</a>
        </v-list-item>
      </v-list>
    </v-menu>
  </div>
</template>

<script>
import movieApi from '@/services/MovieApi'
import DirList from '@/components/menu/DirList'
import EventBus from '@/store/eventBus.js'
export default {
  name: 'Burger',
  components: { DirList },
  props: ['movie'],
  data () {
    return {
      dialog: false,
      showdirlist: false
    }
  },
  mounted () {
    // this.getDirList();
  },
  computed: {
    getMovieItem () {
      return this.movie
    },
    getItems () {
      let _this = this
      return [
        {
          name: function () {
            return 'Is TV Show'
          },
          action: function (item) {
            console.log('child')
            _this.isTvShow(item)
          }
        },
        {
          name: function () {
            return 'Delete Movie'
          },
          action: function (item) {
            console.log('delete')
            _this.dialog = true
          }
        },
        // {
        //   name: function () {
        //     if (_this.movie.watchlist) {
        //       return 'Remove from watchlist'
        //     }
        //     return 'Add to watchlist'
        //   },
        //   action: function (item) {
        //     console.log('watchlsit')
        //     _this.toggleWatchlist(item)
        //   }
        // },
        {
          name: function () {
            return 'Move Movie'
          },
          action: function (item) {
            console.log('move')
            _this.showdirlist = true
          }
        },
        {
          name: function () {
            return 'Show file'
          },
          action: function (item) {
            console.log('show')
            _this.showlocal()
          }
        },
        {
          name: function () {
            return 'Rescan'
          },
          action: function (item) {
            console.log('Rescan')
            _this.reScan(item)
          }
        }
      ]
    }
  },
  methods: {
    getLink (item) {
      return '/?genre=' + item.tmdb_id
    },
    openMenu () {},
    isTvShow (item) {
      item.is_tv = true
      EventBus.$emit('ISTVSHOW', item)
    },
    toggleWatchlist (item) {
      item.watchlist = !item.watchlist
      EventBus.$emit('TOGGLEWATCHLIST', item)
    },
    deleteMovie (item) {
      this.dialog = false
      // this.delete(item);
      EventBus.$emit('DELETEMOVIE', item)
    },
    reScan (item) {
      EventBus.$emit('RESCAN', item)
    },
    showlocal () {
      movieApi
        .playLocal(this.movie, { showdir: 1 })
        .then(response => {})
        .catch(error => {
          console.log(error)
        })
    },
    downloadURL () {
      return this.$baseURL + '/api/file/' + this.movie.file_id + '/download'
      // http://192.168.1.4:8000/file/161021/download
    }
  }
}
</script>

<style lang="stylus" scoped>
.hamburger {
  float: right;
  cursor: pointer;
  position: absolute;
  right: 0px;
  top: 0px;
}

li:hover {
  cursor: pointer;
}
a {
  text-decoration : none
  color :black

}
</style>
