<template>
  <v-app class="main">
    <app-navigation v-if="currentUser"></app-navigation>
    <v-main class="content">
      <router-view></router-view>
    </v-main>
  </v-app>
</template>

<script>
import AppNavigation from '@/components/AppNavigation'

export default {
  name: 'App',
  components: {
    AppNavigation
  },

  data () {
    return {
      searchstring: '',
      searchstring2: '',
      searchboxVisible: false
    }
  },
  // props: ["searchboxVisible"],
  methods: {
    searchMovie () {
      this.$router.push('/search/' + this.searchString)
      this.searchString = ''
    }
  },
  computed: {
    dataAvailable () {
      return this.searchString !== null && this.searchString !== ''
    },
    currentUser () {
      return this.$store.state.auth.user
    }
  }
}
</script>

<style lang="stylus" scoped>
.content {

  width: 100%;
  margin-left: auto;
  margin-right: auto;
}

.main {
  // background : #9e9e9e;
}
</style>
