<template>

  <v-layout row>
    <add-user
    v-show=ShowAdd
    v-model=ShowAdd
    >
    </add-user>
    <v-flex xs12 sm6 offset-sm3>
      <v-card>
        <v-toolbar color="white" >

          <v-toolbar-title>Users</v-toolbar-title>

          <v-spacer></v-spacer>

          <v-btn icon>
            <v-icon @click="ShowAdd=!ShowAdd">add</v-icon>
          </v-btn>
        </v-toolbar>

        <v-list two-line subheader>

          <v-list-tile
            v-for="item in Users"
            :key="item.id"
            avatar

          >
            <v-list-tile-avatar>
              <v-icon :class="[item.iconClass]">{{ item.icon }}</v-icon>
            </v-list-tile-avatar>

            <v-list-tile-content>
              <v-list-tile-title>{{ item.UserName }}</v-list-tile-title>
              <v-list-tile-sub-title>{{ item.subtitle }}</v-list-tile-sub-title>
            </v-list-tile-content>

            <v-list-tile-action>
              <v-btn icon delete>
                <v-icon color="grey lighten-1">delete</v-icon>
              </v-btn>
               </v-list-tile-action>
                  <v-list-tile-action>
               <v-btn icon edit>
                <v-icon color="grey lighten-1">edit</v-icon>
              </v-btn>
            </v-list-tile-action>
          </v-list-tile>

        </v-list>
      </v-card>
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
      items: [
        { name: 'Show User', to: '/user' }
        // { name: 'No Title', to: '/?show=notitle' },
        // { name: 'Duplicates', to: '/?show=duplicate&orderby=name' }
      ]
    }
  },
  mounted () {
    this.User = this.$store.state.auth.user.user_name
    movieApi
      .fetchUsers()
      .then(response => {
        this.Users = response.data
      })
      .catch(error => {
        console.log(error)
      })
  },
  computed: {
    loggedIn () {
      return this.$store.state.auth.status.loggedIn
    }
  },
  methods: {
    handleRegister () {
      this.message = ''
      this.submitted = true
      this.$validator.validate().then(isValid => {
        if (isValid) {
          this.$store.dispatch('auth/register', this.user).then(
            data => {
              this.message = data.message
              this.successful = true
            },
            error => {
              this.message =
                (error.response && error.response.data) ||
                error.message ||
                error.toString()
              this.successful = false
            }
          )
        }
      })
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
