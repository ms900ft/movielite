
<template>
  <v-list-group prepend-icon="build" :value="false">
    <template v-slot:activator>
      <v-list-tile>
        <v-list-tile-title>Layout</v-list-tile-title>
      </v-list-tile>
    </template>

    <v-list>
      <v-list-tile
        v-for="(item, index) in Sizes"
        :key="index"
        @click="changeSize(item)"
      >
        <v-list-tile-action></v-list-tile-action>
        <v-list-tile-title>{{ item.name }}</v-list-tile-title>
      </v-list-tile>
    </v-list>
  </v-list-group>
</template>

<script>
import VueCookies from 'vue-cookies'
export default {
  name: 'DrawerLayout',
  components: {},

  data () {
    return {
      Sizes: [
        { name: '6 Columns', size: '2' },
        { name: '4 Columns', size: '3' },
        { name: '3 Columns', size: '4' },
        { name: '2 Columns', size: '6' },
        { name: '1 Column', size: '12' }
      ],
      Size: '6'
    }
  },
  mounted () {
    this.Size = VueCookies.get('size')
  },
  computed: {},
  methods: {
    changeSize (item) {
      this.Size = item.size
      VueCookies.set('size', this.Size, '365d', null, null, null, 'Lax')
      this.$router.go(this.$router.currentRoute)
    }
  }
}
</script>

<style scoped>
</style>
