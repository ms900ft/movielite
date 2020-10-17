<template>
  <v-list-group prepend-icon="map" :value="false">
    <template v-slot:activator>
      <v-list-tile>
        <v-list-tile-title>Countries</v-list-tile-title>
      </v-list-tile>
    </template>

    <v-list>
      <v-list-tile v-for="(item, index) in orderedCountries" :key="index" @click="country(item)">
        <v-list-tile-action></v-list-tile-action>
        <v-list-tile-title to="getLink(item)">{{ item.name }}</v-list-tile-title>
      </v-list-tile>
    </v-list>
  </v-list-group>
</template>

<script>
import movieApi from '@/services/MovieApi'
export default {
  name: 'DrawerCountries',
  components: {},

  data () {
    return {
      items: []
    }
  },
  mounted () {
    movieApi
      .fetchCountries()
      .then(response => {
        this.items = response.data
      })
      .catch(error => {
        console.log(error)
      })
  },
  computed: {
    orderedCountries: function () {
      return this._.orderBy(this.items, 'name')
    }
  },
  methods: {
    getLink (item) {
      return '/?country=' + item.iso_id
    },
    country (item) {
      this.$router.push('?orderby=name&country=' + item.iso_id)
    }
  }
}
</script>

<style scoped>
</style>
