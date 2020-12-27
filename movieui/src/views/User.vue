<template>
  <v-layout row>

    <add-user v-show="ShowAdd" v-model="ShowAdd"> </add-user>
    <v-flex xs12 sm6 offset-sm3>
           <v-alert :color="type" icon="check_circle" value="true" v-if="type">
        User saved
      </v-alert>
      <v-form v-model="valid">
        <v-card>

          <v-toolbar color="white">
            <v-toolbar-title>Users</v-toolbar-title>

            <v-spacer></v-spacer>

            <v-btn icon>
              <v-icon @click="ShowAdd = !ShowAdd">add</v-icon>
            </v-btn>
          </v-toolbar>

          <v-list two-line subheader>
            <v-list-tile v-for="(item, index) in Users" :key="item.id" avatar>
              <v-list-tile-avatar>
                <v-icon :class="[item.iconClass]">{{ item.icon }}</v-icon>
              </v-list-tile-avatar>

              <v-list-tile-content>
                <v-text-field
                  v-model="item.UserName"
                  key="'username_'+index"

                  label="Username"
                  required
                  :rules="nameRules"
                  v-bind:disabled="!(edit == index)"
                ></v-text-field>
              </v-list-tile-content>
              <v-list-tile-content>
                <v-text-field
                  v-if="edit == index"
                  v-model="item.password"
                  key="'password_'+index"
                  label="Password"
                  type="password"
                  :rules="passwdRules"
                  required
                />
              </v-list-tile-content>
              <v-list-tile-action>
                <v-btn icon delete>
                  <v-icon
                    color="grey lighten-1"
                    @click="deleteUser(item, index)"
                    >delete</v-icon
                  >
                </v-btn>
              </v-list-tile-action>
              <v-list-tile-action>
                <v-btn icon edit >
                  <v-icon v-if="edit==index && !valid"  color="grey lighten-1" @click="cancelUser(item)"
                    >cancel</v-icon
                  >
                  <v-icon v-else-if="edit==index"   color="grey lighten-1" @click="updateUser(item)"
                    >save</v-icon
                  >
                  <v-icon v-else color="grey lighten-1" @click="editUser(index)"
                    >edit</v-icon
                  >
                </v-btn>
              </v-list-tile-action>
            </v-list-tile>
          </v-list>
        </v-card>
      </v-form>
    </v-flex>
  </v-layout>
</template>

<script>
// import User from '../models/user'
import movieApi from '@/services/MovieApi'
import AddUser from '@/components/dialog/AddUser'

export default {
  name: 'User',
  components: {
    AddUser
  },
  data () {
    return {
      Users: [],
      User: 'Users',
      ShowAdd: false,
      edit: -1,
      valid: false,
      type: null,
      elapse: null,

      nameRules: [
        v => !!v || 'Name is required',
        v => v.length <= 10 || 'Name must be less than 10 characters'
      ],
      passwdRules: [
        v => !!v || 'Password is required',
        v => (v && v.length >= 8) || 'Password must be more than 8 characters'
      ],
      items: [
        { name: 'Show User', to: '/user' }
        // { name: 'No Title', to: '/?show=notitle' },
        // { name: 'Duplicates', to: '/?show=duplicate&orderby=name' }
      ]

    }
  },
  mounted () {
    this.User = this.$store.state.auth.user.user_name
    this.getUser()
  },
  computed: {
    loggedIn () {
      return this.$store.state.auth.status.loggedIn
    }
  },
  methods: {
    editUser (index) {
      this.User = this.$store.state.auth.user.user_name
      this.getUser()
      this.edit = index
    },
    deleteUser (item, index) {
      movieApi
        .deleteUser(item)
        .then(response => {
          this.loading = false
          // this.user = {}
        })
        .catch(error => {
          console.log(error)
        })
      this.$delete(this.Users, index)
    },
    updateUser (item, index) {
      movieApi
        .updateUser(item)
        .then(response => {
          this.loading = false
          // this.user = {}
        })
        .catch(error => {
          console.log(error)
        })
      this.edit = undefined
      this.showAlert('success')
    },
    cancelUser (index) {
      this.edit = null
      this.getUser()
    },
    getUser () {
      movieApi
        .fetchUsers()
        .then(response => {
          this.Users = response.data
        })
        .catch(error => {
          console.log(error)
        })
    },
    showAlert (type) {
      this.type = type

      let timer = this.showAlert.timer
      if (timer) {
        clearTimeout(timer)
      }
      this.showAlert.timer = setTimeout(() => {
        this.type = null
      }, 3000)

      this.elapse = 1
      let t = this.showAlert.t
      if (t) {
        clearInterval(t)
      }

      this.showAlert.t = setInterval(() => {
        if (this.elapse === 3) {
          this.elapse = 0
          clearInterval(this.showAlert.t)
        } else {
          this.elapse++
        }
      }, 1000)
    }
  }

}
</script>

<style scoped>
label {
  display: block;
  margin-top: 10px;
}

.card-container.card {
  max-width: 350px !important;
  padding: 40px 40px;
}

.card {
  background-color: #f7f7f7;
  padding: 20px 25px 30px;
  margin: 0 auto 25px;
  margin-top: 50px;
  -moz-border-radius: 2px;
  -webkit-border-radius: 2px;
  border-radius: 2px;
  -moz-box-shadow: 0px 2px 2px rgba(0, 0, 0, 0.3);
  -webkit-box-shadow: 0px 2px 2px rgba(0, 0, 0, 0.3);
  box-shadow: 0px 2px 2px rgba(0, 0, 0, 0.3);
}

.profile-img-card {
  width: 96px;
  height: 96px;
  margin: 0 auto 10px;
  display: block;
  -moz-border-radius: 50%;
  -webkit-border-radius: 50%;
  border-radius: 50%;
}
</style>
