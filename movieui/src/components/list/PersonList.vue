<template>
  <div class="column crewlist">
    <personsModal
      v-show="fullimage"
      @close="fullimage = false"
      :data="orderedPersons"
      v-model="fullimage"
      :startIndex="startIndex"
    />

    <v-list>
      <v-flex v-for="(item, index) in orderedPersons" :key="index">
        <div class="rowreg" style="margin-bottom: 10px">
          <div class="rowsmall">
            <img
            v-if=item.profile_path
              :src="castimage(item.profile_path)"
              width="45px"
              height="68px"
              style="margin-right: 10px"
              @click="showfullimage(item, index)"
            />
            <v-icon  v-else size="45px" color="#337ab7" >mdi-account</v-icon>
          </div>
          <div class="rowsmall">
            <!-- <img
              v-if="showfullimage(index)"
              :src="castimage(item.profile_path,300)"

              class="fullimage"
            /> -->
          </div>
          <div
            @click="open(item, index)"
            class="rowsmall cast"
            style="flex: 30%; margin: auto; display: block"
          >
            <p style="margin: 0">{{ item.Name }}</p>
            <p style="color: black; margin: 0; font-size: small">
              {{ desc(item) }}
            </p>
          </div>
        </div>
      </v-flex>
    </v-list>
  </div>
</template>

<script>
import personsModal from '@/components/PersonsModal'

export default {
  name: 'PersonList',
  components: { personsModal },

  data () {
    return {
      items: [],
      // collapsed: true
      // showdir: true
      fullimage: false,
      startIndex: 0,
      person: {}
    }
  },
  props: ['persons', 'type'],

  mounted () {},
  computed: {
    orderedPersons: function () {
      const order = {
        Director: 1,
        Producer: 2,
        Writer: 3,
        Story: 4,
        Screenplay: 5
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
    }
  },

  methods: {
    castimage (path, size) {
      var pic
      if (!size) {
        size = 45
      }
      path
        ? (pic = this.$baseURL + '/images/w' + size + path)
        : (pic = this.$baseURL + '/movie2/na.gif')
      return pic
    },
    open (item) {
      this.$emit('closeModal')
      this.$router.push('?person=' + item.ID + '&orderby=name')
      this.show = false
    },
    desc (item) {
      if (this.type === 'crew') {
        return item.Job
      } else {
        return item.Character
      }
    },
    showfullimage (item, index) {
      if (item.profile_path !== '') {
        this.fullimage = !this.fullimage
        this.person = item
        this.startIndex = index
      }
    }
  }
}
</script>

<style scoped>
.fullimage {
  overflow: visible;
  position: absolute;
  z-index: 100;
  float: left;
  max-width: 100%;
  max-height: 100vh;
  overflow-y: auto;
  /* buttom: 20% */

  margin-bottom: 300px;
}
</style>
