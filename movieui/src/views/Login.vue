<template>
  <v-app>
    <v-card width="400" class="mx-auto mt-5">
      <form @submit.prevent="handleLogin" @keyup.enter="handleLogin">
      <v-card-title class="pb-0">
        <v-layout justify-center>
        <v-icon   size="100px">account_circle</v-icon>
        </v-layout>
         <v-alert color="error" icon="check_circle" value="true" dismissible outline v-if="message">
        {{message.error}}
      </v-alert>
      </v-card-title>
      <v-card-text>

          <v-text-field
            label="Username"
            v-model="user.username"
          />
          <v-text-field
            type="password"
            label="Password"
            v-model="user.password"

          />

      </v-card-text>
      <v-divider></v-divider>
      <v-card-actions>
<v-layout justify-center>
        <v-btn type="submit" color="info" >Login</v-btn>
</v-layout>
      </v-card-actions>
</form>
    </v-card>
  </v-app>
</template>
<script>
import User from '../models/user'

export default {
  name: 'Login',
  data () {
    return {
      user: new User('', ''),
      loading: false,
      message: ''
    }
  },
  computed: {
    loggedIn () {
      return this.$store.state.auth.status.loggedIn
    }
  },
  created () {
    if (this.loggedIn) {
      this.$router.push('/')
    }
  },
  methods: {
    handleLogin () {
      this.loading = true
      this.$validator.validateAll().then(isValid => {
        if (!isValid) {
          this.loading = false
          return
        }

        if (this.user.username && this.user.password) {
          this.$store.dispatch('auth/login', this.user).then(
            () => {
              this.$router.push('/')
            },
            error => {
              this.loading = false
              this.message =
                (error.response && error.response.data) ||
                error.message ||
                error.toString()
            }
          )
        }
      })
    }
  }
}
</script>

<style scoped>

</style>
