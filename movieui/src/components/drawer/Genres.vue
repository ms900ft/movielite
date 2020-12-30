<template>
  <v-list-group prepend-icon="movie" :value="false">
    <template v-slot:activator>

        <v-list-item-title>Genres</v-list-item-title>

    </template>

    <v-list>
      <v-list-item v-for="(item, index) in orderedGenres" :key="index" @click="genre(item)">
        <v-list-item-action></v-list-item-action>
        <v-list-item-title to="getLink(item)">{{ item.name }}</v-list-item-title>
      </v-list-item>
    </v-list>
  </v-list-group>
</template>

<script>
import movieApi from '@/services/MovieApi'
export default {
  name: 'DrawerGenres',
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
