<template>
  <div class="text-xs-center">
    <v-menu offset-y max-height="300">
      <template v-slot:activator="{ on }">
        <v-btn text v-on="on">Genres</v-btn>
      </template>
      <v-list>
        <v-list-item v-for="(item, index) in orderedGenres" :key="index" @click="genre(item)">
          <v-list-item-title to="getLink(item)">{{ item.name }}</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-menu>
  </div>
</template>

<script>
import movieApi from '@/services/MovieApi'

export default {
  name: 'MenuGenres',
  components: {},

  data () {
    return {
      items: [{ name: 'test' }, { name: 'test2' }]
    }
  },
  mounted () {
    movieApi
      .fetchGenres()
      .then(response => {
        this.items = response.data
        const genreMap = this._.keyBy(this.items, 'tmdb_id')

        this.$store.commit('setGenreMap', genreMap)
      })
      .catch(error => {
        console.log(error)
      })
  },
  computed: {
    orderedGenres: function () {
      return this._.orderBy(this.items, 'name')
    }
  },
  methods: {
    getLink (item) {
      return '/?genre=' + item.tmdb_id
    },
    genre (item) {
      this.$router.push('/?orderby=name&genre=' + item.tmdb_id)
    }
  }
}
</script>

<style scoped>
</style>
