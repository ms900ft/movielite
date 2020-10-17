
<template>
  <v-list-group prepend-icon="work" :value="false">
    <template v-slot:activator>
      <v-list-tile>
        <v-list-tile-title>Work</v-list-tile-title>
      </v-list-tile>
    </template>

    <v-list>
      <v-list-tile v-for="(item, index) in items" :key="index" @click="gotoItem(item)">
        <v-list-tile-action></v-list-tile-action>

        <v-list-tile-title>{{ item.name }}</v-list-tile-title>
      </v-list-tile>
    </v-list>
  </v-list-group>
</template>

<script>
export default {
  name: 'DrawerWork',
  components: {},

  data () {
    return {
      items: [
        { name: 'Add Meta', to: '/?show=multiple' },
        { name: 'No Title', to: '/?show=notitle' },
        { name: 'Duplicates', to: '/?show=duplicate&orderby=name' }
      ]
    }
  },
  mounted () {},
  computed: {},
  methods: {
    gotoItem (item) {
      this.$router.push(item.to)
    }
  }
}
</script>

<style scoped>
</style>
