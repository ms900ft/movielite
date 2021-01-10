<template>
  <span>
    <sidebar :show="showsidebar" @closesidebar="showsidebar=false"></sidebar>
    <v-app-bar app color="blue darken-4" dark dense :scroll-threshold="3" scroll-off-screen>
      <!-- <template v-if="showExtended()" #extension>
        <v-toolbar-items>
          <v-chip v-model="chip1" outline dark style="color: white;" close>search: {{searchstring}}</v-chip>
          <v-btn text>Link 2</v-btn>
          <v-btn text>Link 3</v-btn>
        </v-toolbar-items>
      </template>-->
      <v-app-bar-nav-icon  @click.stop="showsidebar=!showsidebar"></v-app-bar-nav-icon>
      <v-btn class="hidden-sm-and-down" text exact to="/?orderby=name">
        <v-icon>mdi-home</v-icon>
      </v-btn>
      <v-btn class="hidden-sm-and-down" text exact to="/?show=watchlist">Watchlist</v-btn>
      <v-btn class="hidden-sm-and-down" text exact to="/?orderby=recent">Recently</v-btn>
      <menu-genres class="hidden-sm-and-down"></menu-genres>
      <menu-countries class="hidden-sm-and-down"></menu-countries>

      <v-text-field
        v-model="searchstring"
        text

        clearable
        autofocus
        height="30"
        placeholder="Search..."
        class="search"
        single-line
        hide-details
      ></v-text-field>

      <v-btn icon class="hidden-sm-and-down" v-on:click="keybordVisible">
        <v-icon>mdi-keyboard</v-icon>
      </v-btn>
      <v-spacer></v-spacer>
      <transition name="fade">
        <div v-if="showkeybord" class="fixed">
          <keyboard
            class="keybord"
            v-model="searchstring"
            @custom="custom"
            @input="changed"
            :layouts="[
    '1234567890{delete:backspace}|qwertzuiop|asdfghjkl|{shift:goto:1}yxcvbnm|{space:space}{custom:custom}',
    '!@#$%^&*(){delete:backspace}|QWERTZUIOP|ASDFGHJKL|{shift:goto:0}YXCVBNM|{space:space}{custom:custom}'
  ]"
            :maxlength="16"
          ></keyboard>
        </div>
      </transition>
<v-btn class="hidden-sm-and-down" text exact @click.prevent="logOut">logout</v-btn>
      <menu-user class="hidden-sm-and-down"></menu-user>

      <v-chip outlined dark style="color: white;">found: {{total()}}</v-chip>
    </v-app-bar>
  </span>
</template>

<script>
import MenuGenres from '@/components/menu/Genres'
import MenuCountries from '@/components/menu/Countries'
import MenuUser from '@/components/menu/User'
import keyboard from 'vue-keyboard'
import Sidebar from '@/components/Sidebar'
export default {
  name: 'AppNavigation',
  components: {
    MenuGenres,
    MenuCountries,
    MenuUser,
    keyboard,
    Sidebar
  },
  data () {
    return {
      showkeybord: false,
      showsidebar: false,
      extendedSlot: true,
      chip1: false
    }
  },
  props: ['genre', 'person'],
  methods: {
    showExtended () {
      if (this.$store.state.searchString) {
        this.chip1 = true
        return true
      }
      return false
    },
    closeSearch () {
      this.chip1 = false
      this.$store.state.searchString = ''

      return false
    },
    changed (value) {
      console.log('Value ' + value)
    },
    custom (value) {
      console.log('Value ' + value)
    },
    keybordVisible () {
      this.showkeybord = !this.showkeybord
    },
    total () {
      return this.$store.state.resultsFound
    },
    logOut () {
      this.$store.dispatch('auth/logout')
      this.$router.push('/login')
    }
  },
  computed: {
    searchstring: {
      get () {
        return this.$store.state.searchString
      },
      set (value) {
        if (value === null) {
          value = ''
        }
        this.$store.commit('setSearchString', value)
      }
    }
  }
}
</script>

<style scoped>
.my-input.v-input .v-input__slot {
  border-radius: 100px;
  border-width: 10px;
  border-color: black;
}

.fixed {
  position: fixed;
  top: 60px;
  right: 10px;
  bottom: 0;
  overflow: visible;

  z-index: 100;
}

.keybord {
  background-color: white;
  padding: 20px;
  border-width: 2px;
  border-style: solid;
  border-color: black;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 1s;
}
.fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */ {
  opacity: 0;
}

.search{
  display: block;
  right: 0;
  bottom: 0;
  position: flex;
  width: 150px;

}

a{
    text-decoration: none !important;
}

</style>
