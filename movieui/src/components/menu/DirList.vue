<template>
  <v-dialog v-model="showdirs" @close="showdirs(false)" max-width="50%"  scrollable>
    <v-layout row>
      <v-flex>
        <v-card>
          <v-list    min-height="800">
            <v-subheader class="justify-center headline">Move movie to</v-subheader>

            <v-list-item
              v-for="(item) in orderedTargets"
              :key="item.name"
              v-on:click="move(item.name)"
            >
              <v-list-item-icon class="ml-8">
                <v-icon color="blue darken-4" large>mdi-folder</v-icon>
              </v-list-item-icon>

              <v-list-item-content>
                <v-list-item-title class="dirlist">{{ item.name }}</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list>
        </v-card>
      </v-flex>
    </v-layout>
  </v-dialog>
</template>

<script>
import movieApi from '@/services/MovieApi'
export default {
  name: 'DirList',
  components: {},

  data () {
    return {
      items: []
      // showdir: true
    }
  },
  props: ['show', 'movie'],

  mounted () {
    movieApi
      .fetchTargets()
      .then(response => {
        this.items = response.data
      })
      .catch(error => {
        console.log(error)
      })
  },
  computed: {
    orderedTargets: function () {
      return this._.orderBy(this.items, [item => item.name.toLowerCase()])
    },
    showdirs: {
      get () {
        return this.show
      },
      set (value) {
        this.$emit('close')
      }
    }
  },
  methods: {
    move (where) {
      movieApi
        .moveMovie(this.movie, where)
        .then(response => {
          this.loading = false
        })
        .catch(error => {
          console.log(error)
        })
      this.$emit('close')
    }
  }
}
</script>

<style scoped>
.dirlist {
  color: #0D47A1;
  font-size: 24px;
  font-weight: 500;
  cursor: pointer;
  padding-left: 10px;

}

.headline {
  color: white;
  background-color: #0D47A1;
  font-size: 24px;
  font-weight: 500;
  cursor: pointer;
  padding-top: 10px;
}
</style>
