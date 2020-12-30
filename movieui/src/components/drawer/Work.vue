
<template>
  <v-list-group prepend-icon="work" :value="false">
    <template v-slot:activator>
        <v-list-item-title>Work</v-list-item-title>
    </template>

    <v-list>
      <v-list-item v-for="(item, index) in items" :key="index" @click="gotoItem(item)">
        <v-list-item-action></v-list-item-action>

        <v-list-item-title>{{ item.name }}</v-list-item-title>
      </v-list-item>
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
