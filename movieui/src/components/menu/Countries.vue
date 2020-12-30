<template>
  <div class="text-xs-center">
    <v-menu offset-y max-height="300">
      <template v-slot:activator="{ on }">
        <v-btn text v-on="on">Countries</v-btn>
      </template>
      <v-list>
        <v-list-item v-for="(item, index) in orderedCountries" :key="index" @click="country(item)">
          <v-list-item-title to="getLink(item)">{{ item.name }}</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-menu>
  </div>
</template>

<script>
import movieApi from '@/services/MovieApi'
export default {
  name: 'MenuCountries',
  components: {},

  data () {
    return {
      items: [{ name: 'test' }, { name: 'test2' }]
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
      this.$router.push('/?orderby=name&country=' + item.iso_id)
    }
  }
}
</script>

<style scoped>
</style>
