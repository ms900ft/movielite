<template>
    <v-dialog v-model="show" persistent max-width="600px">

      <v-card>
        <form name="form" @submit="addUser">
        <v-card-title>

          <span class="headline">Add User</span>
        </v-card-title>
        <v-card-text>

          <v-container grid-list-md>
            <v-layout wrap>
              <v-flex xs12 >
                <v-text-field label="Username*" :rules="nameRules" required v-model="user.username" :counter="max"
            ></v-text-field>
              </v-flex>

              <v-flex xs12>
                <v-text-field label="Password*" type="password" v-model="user.password" :rules="rules" required></v-text-field>
              </v-flex>
               <v-flex xs12>
                <v-text-field label="Password repeat*" type="password" v-model="user.password2" required :rules="matchrules"></v-text-field>
              </v-flex>
              <v-flex xs12 sm6>
                <v-select
                  :items="['admin', 'user']"
                  label="Type*"
                  required
                ></v-select>
              </v-flex>
            </v-layout>
          </v-container>
          <small>*indicates required field</small>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" flat @click="show = false">Close</v-btn>
          <v-btn  color="blue darken-1" flat type="submit">Save</v-btn>
        </v-card-actions>
        </form>
      </v-card>
    </v-dialog>

</template>

<script>
import User from '@/models/user'
import movieApi from '@/services/MovieApi'

export default {
  name: 'AddUser',
  props: {
    value: Boolean
  },
  data: () => ({
    // dialog: this.show
    // on: this.show
    user: new User('', ''),
    max: 15,
    min: 8,
    match: '',
    nameRules: [
      v => !!v || 'Name is required',
      v => (v && v.length <= 10) || 'Name must be less than 10 characters'
    ]

  }),
  computed: {
    show: {
      get () {
        return this.value
      },
      set (value) {
        this.$emit('input', value)
      }
    },

    rules () {
      const rules = []
      // var match = this.user.password

      if (this.max) {
        const rule =
            v => (v || '').length <= this.max ||
              `A maximum of ${this.max} characters is allowed`

        rules.push(rule)
      }
      if (this.min) {
        const rule =
            v => (v || '').length >= this.min ||
              `A minimum of ${this.min} characters required`

        rules.push(rule)
      }

      if (!this.allowSpaces) {
        const rule =
            v => (v || '').indexOf(' ') < 0 ||
              'No spaces are allowed'

        rules.push(rule)
      }

      return rules
    },
    matchrules () {
      var match = this.user.password
      const rules = []
      if (match) {
        const rule =
            v => (!!v && v) === match ||
              'Values do not match'

        rules.push(rule)
      }

      return rules
    }
  },

  methods: {
    addUser () {
      // this.$validator.validateAll().then(isValid => {
      //   if (!isValid) {
      //     this.loading = false
      //     return
      //   }
      console.log('------------------------------------')
      console.log('submit')
      console.log('------------------------------------')
      if (this.user.username && this.user.password) {
        movieApi
          .addUser(this.user)
          .then((response) => {
            this.loading = false
            this.user = {}
          })
          .catch((error) => {
            console.log(error)
          })
      }
    },
    validateField () {
      this.$refs.form.validate()
    }
  },
  watch: {
    match: 'validateField',
    max: 'validateField',
    min: 'validateField',
    model: 'validateField'
  }

}
</script>
