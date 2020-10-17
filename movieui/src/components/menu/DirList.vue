<template>
  <v-dialog v-model="showdirs" @close="showdirs(false)" width="350">
    <v-layout row>
      <v-flex>
        <v-card>
          <v-list two-line subheader>
            <v-subheader class="headline">Move movie to:</v-subheader>

            <v-list-tile
              ripple
              v-for="(item) in orderedTargets"
              :key="item.name"
              v-on:click="move(item.name)"
            >
              <v-list-tile-avatar>
                <v-icon class="blue lighten-1 white--text">folder</v-icon>
              </v-list-tile-avatar>

              <v-list-tile-content>
                <v-list-tile-title class="movietitle">{{ item.name }}</v-list-tile-title>
              </v-list-tile-content>
            </v-list-tile>
          </v-list>
        </v-card>
      </v-flex>
    </v-layout>
  </v-dialog>
</template>

<script>
import movieApi from '@/services/MovieApi'
export default {
  name: 'DirList',
  components: {},

  data () {
    return {
      items: []
      // showdir: true
    }
  },
  props: ['show', 'movie'],

  mounted () {
    movieApi
      .fetchTargets()
      .then(response => {
        this.items = response.data
      })
      .catch(error => {
        console.log(error)
      })
  },
  computed: {
    orderedTargets: function () {
      return this._.orderBy(this.items, [item => item.name.toLowerCase()])
    },
    showdirs: {
      get () {
        return this.show
      },
      set (value) {
        this.$emit('close')
      }
    }
  },
  methods: {
    move (where) {
      movieApi
        .moveMovie(this.movie, where)
        .then(response => {
          this.loading = false
        })
        .catch(error => {
          console.log(error)
        })
      this.$emit('close')
    }
  }
}
</script>

<style scoped>
</style>
