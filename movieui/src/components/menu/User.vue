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
// import VueCookies from 'vue-cookies'
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
      this.$store.dispatch('auth/logout')
      this.$router.push('/login?username=' + item.UserName)
    }
  }
}
</script>

<style scoped>
</style>
