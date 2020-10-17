<template>
  <div class="text-xs-center">
    <v-menu offset-y max-height="300">
      <template v-slot:activator="{ on }">
        <v-btn flat v-on="on">Work</v-btn>
      </template>
      <v-list>
        <v-list-tile v-for="(item, index) in items" :key="index" @click="gotoItem(item)">
          <v-list-tile-title>{{ item.name }}</v-list-tile-title>
        </v-list-tile>
      </v-list>
    </v-menu>
  </div>
</template>

<script>

export default {
  name: 'MenuWork',
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
