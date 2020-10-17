<template>
  <div class="text-xs-center">
    <v-menu offset-y max-height="300">
      <template v-slot:activator="{ on }">
        <v-btn flat v-on="on">{{User}}</v-btn>
      </template>
      <v-list>
        <v-list-tile v-for="(item, index) in Users" :key="index" @click="changeUser(item)">
          <v-list-tile-title>{{ item.name }}</v-list-tile-title>
        </v-list-tile>
      </v-list>
    </v-menu>
  </div>
</template>

<script>
import VueCookies from 'vue-cookies'
export default {
  name: 'MenuUser',
  components: {},

  data () {
    return {
      Users: [{ name: 'Marc' }, { name: 'Kaja' }],
      User: 'Users'
    }
  },
  mounted () {
    this.User = VueCookies.get('movieuser')
  },
  computed: {},
  methods: {
    changeUser (item) {
      this.User = item.name
      VueCookies.set('movieuser', this.User, '365d')
      this.$router.go(this.$router.currentRoute)
    }
  }
}
</script>

<style scoped>
</style>
