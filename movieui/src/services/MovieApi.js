import axios from 'axios'
// axios.defaults.withCredentials = true
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
    const response = await axios.get('movie?' + queryString)
    console.log(response)
    return response.data
  },

  fetchSingleMovie (id) {
    return axios.get('movie/' + id)
      .then(response => {
        return response.data
      })
  },
  moveMovie (movie, where) {
    return axios.put('file/' + movie.file_id + '/move/' + encodeURIComponent(where))
      .then(response => {
        return response.data
      })
  },
  addMeta (movie, metaid) {
    return axios.put('movie/' + movie.id + '/addMeta/' + metaid, movie)
      .then(response => {
        return response.data
      })
  },
  async fetchGenres () {
    const response = await axios.get('genre')
    return response
  },
  async fetchCountries () {
    const response = await axios.get('country')
    return response
  },
  async fetchTargets () {
    const response = await axios.get('targets')
    return response
  },
  async fetchUsers () {
    const response = await axios.get('user')
    return response
  },
  playLocal (item, args = {}) {
    const queryString = Object.keys(args).map(key => key + '=' + args[key]).join('&')
    return axios.put('movie/' + item.id + '/play?' + queryString, {})
      .then(response => {
        return response
      })
  },
  updateMovie (item) {
    return axios.put('movie/' + item.id, item)
      .then(response => {
        return response
      })
  },
  deleteMovie (item) {
    return axios.delete('movie/' + item.id, item)
      .then(response => {
        return response
      })
  }

}
