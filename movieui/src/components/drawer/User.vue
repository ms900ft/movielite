
<template>
  <!-- <v-list-group prepend-icon="account_circle" :value="false"> -->
    <!-- <template v-slot:activator> -->
        <v-list-item >
          <v-list-item-action>
            <v-icon>account_circle</v-icon>
          </v-list-item-action>

          <v-list-item-content>
            <v-list-item-title>{{User}}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
    <!-- </template> -->

    <!-- <v-list>
      <v-list-item v-for="(item, index) in items" :key="index" @click="gotoItem(item)">
        <v-list-item-action></v-list-item-action>
        <v-list-item-title>{{ item.name }}</v-list-item-title>
      </v-list-item>
    </v-list> -->
  <!-- </v-list-group> -->
</template>

<script>
import VueCookies from 'vue-cookies'
import movieApi from '@/services/MovieApi'
export default {
  name: 'DrawerUser',
  components: {},

  data () {
    return {
      Users: [],
      User: 'Users',
      items: [
        // { name: 'Show User', to: '/user' }
        // { name: 'No Title', to: '/?show=notitle' },
        // { name: 'Duplicates', to: '/?show=duplicate&orderby=name' }
      ]
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
      VueCookies.set('movieuser', this.User, '365d')
      this.$router.go(this.$router.currentRoute)
    },
    gotoItem (item) {
      this.$router.push(item.to)
    }
  }
}
</script>

<style scoped>

</style>
