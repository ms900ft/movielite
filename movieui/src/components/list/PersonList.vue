<template>
  <div class="column crewlist">
    <v-list>
      <v-flex v-for="(item, index) in orderedPersons" :key="index" @click="open(item)">
        <div class="row" style="margin-bottom: 10px">
          <div class="rowsmall">
            <img
              :src="castimage(item.profile_path)"
              width="45px"
              height="68px"
              style="margin-right: 10px;"
            />
          </div>
          <div class="rowsmall cast" style="flex: 30%; margin: auto; display: block">
            <p style="margin: 0;">{{item.Name}}</p>
            <p style="color: black; margin:0; font-size: small">{{desc(item)}}</p>
          </div>
        </div>
      </v-flex>
    </v-list>
  </div>
</template>

<script>
export default {
  name: 'PersonList',
  components: {},

  data () {
    return {
      items: []
      // collapsed: true
      // showdir: true
    }
  },
  props: ['persons', 'type'],

  mounted () {},
  computed: {
    orderedPersons: function () {
      const order = {
        'Director': 1,
        'Producer': 2,
        'Writer': 3,
        'Story': 4,
        'Screenplay': 5
      }
      var list = this.persons
      return list.sort(function (a, b) {
        var aa, bb
        if (order[a.Job]) {
          aa = order[a.Job]
        } else {
          aa = 100
        }
        if (order[b.Job]) {
          bb = order[b.Job]
        } else {
          bb = 100
        }

        if (aa > bb) {
          return 1
        }
        if (aa < bb) {
          return -1
        }
        return 0
      })
    } },

  methods: {
    castimage (path) {
      var pic
      path
        ? (pic = this.$baseURL + '/images/w45' + path)
        : (pic = this.$baseURL + '/movie2/na.gif')
      return pic
    },
    open (item) {
      this.$emit('closeModal')
      this.$router.push('?person=' + item.ID)
      this.show = false
    },
    desc (item) {
      if (this.type === 'crew') {
        return item.Job
      } else {
        return item.Character
      }
    }

  }
}
</script>

<style scoped>
</style>
