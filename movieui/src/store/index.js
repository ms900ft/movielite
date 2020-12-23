import Vue from 'vue'
import Vuex from 'vuex'
import { auth } from './auth.module'

Vue.use(Vuex)

const store = new Vuex.Store({
  state: {
    searchString: '',
    resultsFound: '',
    genreMap: {}
  },
  actions: {
    getSearchString: function (state) {
      return state.searchString
    },
    getResultsFound: function (state) {
      return state.resultsFound
    },
    getGenreMap: function (state) {
      return state.genreMap
    }
  },
  mutations: {
    setSearchString: function (state, value) {
      state.searchString = value
    },
    setResultsFound: function (state, value) {
      state.resultsFound = value
    },
    setGenreMap: function (state, value) {
      state.genreMap = value
    }
  },
  getters: {},
  modules: {
    auth
  }
})

export default store
