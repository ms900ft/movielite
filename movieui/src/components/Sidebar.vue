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

        <v-list-item v-for="item in items" :key="item.title" :to="{path:item.to}" exact>
          <v-list-item-icon>
            <v-icon>{{ item.icon }}</v-icon>
          </v-list-item-icon>

          <v-list-item-content>
            <v-list-item-title>{{ item.title }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>

        <drawer-genres></drawer-genres>
        <drawer-countries></drawer-countries>
        <drawer-work></drawer-work>
        <drawer-users  ></drawer-users>
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
  mounted: function () {
    this.User = this.$store.state.auth.user
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
      ],
      User: {}
    }
  },
  props: ['show']
}
</script>
