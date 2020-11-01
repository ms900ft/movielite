
<template>
  <v-list-group prepend-icon="account_circle" :value="false">
    <template v-slot:activator>
      <v-list-tile>
        <v-list-tile-title>{{User}}</v-list-tile-title>
      </v-list-tile>
    </template>

    <v-list>
      <v-list-tile v-for="(item, index) in Users" :key="index" @click="changeUser(item)">
        <v-list-tile-action></v-list-tile-action>
        <v-list-tile-title>{{ item.name }}</v-list-tile-title>
      </v-list-tile>
    </v-list>
  </v-list-group>
</template>

<script>
import VueCookies from 'vue-cookies'
export default {
  name: 'DrawerUser',
  components: {},

  data () {
    return {
      Users: [{ name: 'Marc' }, { name: 'Kaja' }, { name: 'covid-19kino' }],
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
      VueCookies.sext('movieuser', this.User, '365d', null, null, null, 'Lax')
      // this.$cookies.set("use_path_argument","value",null, null, null, null, "Lax");
      this.$router.go(this.$router.currentRoute)
    }
  }
}
</script>

<style scoped>
</style>
