<template>
  <v-layout wrap>
    <v-navigation-drawer
      v-model="visible"
      temporary
      dark
      width="250px"
      style="position:fixed; top:0; left:0; overflow-y:scroll;"
    >
      <v-list class="pt-0" dense>
        <v-divider></v-divider>

        <v-list-tile v-for="item in items" :key="item.title" :to="{path:item.to}" exact>
          <v-list-tile-action>
            <v-icon>{{ item.icon }}</v-icon>
          </v-list-tile-action>

          <v-list-tile-content>
            <v-list-tile-title>{{ item.title }}</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <drawer-genres></drawer-genres>
        <drawer-countries></drawer-countries>
        <drawer-work></drawer-work>
        <drawer-users></drawer-users>
        <drawer-user></drawer-user>
      </v-list>
    </v-navigation-drawer>
  </v-layout>
</template>

<script>
import DrawerGenres from '@/components/drawer/Genres'
import DrawerCountries from '@/components/drawer/Countries'
import DrawerWork from '@/components/drawer/Work'
import DrawerUser from '@/components/drawer/User'
import DrawerUsers from '@/components/drawer/Users'

export default {
  components: {
    DrawerGenres,
    DrawerCountries,
    DrawerWork,
    DrawerUser,
    DrawerUsers
  },
  computed: {
    visible: {
      get () {
        return this.show
      },
      set (value) {
        if (value === false) {
          this.$emit('closesidebar')
        }
      }
    }
  },
  data: function () {
    return {
      // visible: this.show,

      items: [
        { title: 'Home', icon: 'home', to: '/?orderby=name' },
        { title: 'Watchlist', icon: 'star', to: '?show=watchlist' },
        { title: 'Recently', icon: 'star', to: '?orderby=recent' }
      ]
    }
  },
  props: ['show']
}
</script>
