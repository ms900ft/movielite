import axios from 'axios'
import _ from 'lodash'
import authHeader from './auth-header'

axios.interceptors.response.use(function (response) {
  return response
}, function (error) {
  console.log(error.response.data)
  if (error.response.status === 401) {
    localStorage.removeItem('user')
    window.location.href = '/login'
  }
  return Promise.reject(error)
})

export default {

  async fetchMovieCollection (vue, args) {
    const limit = vue.$hitspp
    const offset = (args.page - 1) * limit

    const params = _.pickBy({
      limit: limit,
      offset: offset,
      ...args
    })

    const queryString = Object.keys(params).map(key => key + '=' + params[key]).join('&')
    const response = await axios.get('/api/movie?' + queryString, { headers: authHeader() })
    // console.log(response)
    return response.data
  },

  fetchSingleMovie (id) {
    return axios.get('/api/movie/' + id, {
      headers: authHeader()
    })
      .then(response => {
        return response.data
      })
  },
  moveMovie (movie, where) {
    return axios.put('/api/file/' + movie.file_id + '/move/' + encodeURIComponent(where), {}, {
      headers: authHeader()
    })
      .then(response => {
        return response.data
      })
  },
  addMeta (movie, metaid) {
    return axios.put('/api/movie/' + movie.id + '/addMeta/' + metaid, movie, {
      headers: authHeader()
    })
      .then(response => {
        return response.data
      })
  },
  async fetchGenres () {
    const response = await axios.get('/api/genre', {
      headers: authHeader()
    })
    return response
  },
  async fetchCountries () {
    const response = await axios.get('/api/country', {
      headers: authHeader()
    })
    return response
  },
  async fetchTargets () {
    const response = await axios.get('/api/targets', {
      headers: authHeader()
    })
    return response
  },
  async fetchUsers () {
    const response = await axios.get('/api/user', {
      headers: authHeader()
    })
    return response
  },
  playLocal (item, args = {}) {
    const queryString = Object.keys(args).map(key => key + '=' + args[key]).join('&')
    return axios.put('/api/movie/' + item.id + '/play?' + queryString, {},
      { headers: authHeader() })
      .then(response => {
        return response
      })
  },
  updateMovie (item) {
    return axios.put('/api/movie/' + item.id, item, {
      headers: authHeader()
    })
      .then(response => {
        return response
      })
  },
  deleteMovie (item) {
    return axios.delete('/api/movie/' + item.id, item, {
      headers: authHeader()
    })
      .then(response => {
        return response
      })
  },
  addUser (item) {
    return axios.post('/api/user', item, {
      headers: authHeader()
    })
      .then(response => {
        return response
      })
  },
  deleteUser (item) {
    return axios.delete('/api/user/' + item.id, {
      headers: authHeader()
    })
      .then(response => {
        return response
      })
  },
  updateUser (item) {
    return axios.put('/api/user/' + item.id, item, {
      headers: authHeader()
    })
      .then(response => {
        return response
      })
  }
}
