<template>
  <div>
    <a :href="downloadURL()" download style="text-decoration: none">
      <v-btn v-if="!overlay" class="hidden-sm-and-down" outline color="primary" dark>
        Download
        <v-icon right dark>arrow_downward</v-icon>
      </v-btn>
      <v-icon v-else style="color: grey" left size="60">arrow_downward</v-icon>
    </a>
  </div>
</template>

<script>

export default {
  name: 'Download',
  components: {},
  props: ['movie', 'overlay'],
  data () {
    return {}
  },

  methods: {
    downloadURL () {
      return this.$baseURL + '/file/' + this.movie.file_id + '/download'
      // http://192.168.1.4:8000/file/161021/download
    }
  }
}
</script>

<style lang="stylus" scoped>
.hamburger {
  float: right;
}
</style>
