<template>
  <div class="text-xs-center">
    <v-menu offset-y max-height="300">
      <template v-slot:activator="{ on }">
        <v-btn text v-on="on">{{User}}</v-btn>
      </template>
      <v-list>
        <v-list-item v-for="(item, index) in Users" :key="index" @click="changeUser(item)">
          <v-list-item-title>{{ item.UserName }}</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-menu>
  </div>
</template>

<script>
import VueCookies from 'vue-cookies'
import movieApi from '@/services/MovieApi'
export default {
  name: 'MenuUser',
  components: {},

  data () {
    return {
      Users: [],
      User: 'Users'
    }
  },
  mounted () {
    this.User = this.$store.state.auth.user.user_name
    movieApi
      .fetchUsers()
      .then(response => {
        this.Users = response.data
      })
      .catch(error => {
        console.log(error)
      })
  },
  computed: {},
  methods: {
    changeUser (item) {
      this.User = item.UserName
      VueCookies.set('movieuser', this.User, '365d', null, null, null, 'Lax')
      this.$router.go(this.$router.currentRoute)
    }
  }
}
</script>

<style scoped>
</style>
