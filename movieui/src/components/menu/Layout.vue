<template>
  <div class="text-xs-center">
    <v-menu offset-y max-height="300">
      <template v-slot:activator="{ on }">
        <v-btn flat v-on="on">Layout</v-btn>
      </template>
      <v-list>
        <v-list-tile
          v-for="(item, index) in Sizes"
          :key="index"
          @click="changeSize(item)"
        >
          <v-list-tile-title>xxxx{{ item.name }}</v-list-tile-title>
        </v-list-tile>
      </v-list>
    </v-menu>
  </div>
</template>

<script>
import VueCookies from 'vue-cookies'
export default {
  name: 'Layout',
  components: {},

  data () {
    return {
      Sizes: [
        { name: '12' },
        { name: '6' },
        { name: '4' },
        { name: '3' },
        { name: '2' }
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
      this.Size = item.name
      VueCookies.set('size', this.Size, '365d', null, null, null, 'Lax')
      this.$router.go(this.$router.currentRoute)
    }
  }
}
</script>

<style scoped>
</style>
