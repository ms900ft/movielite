import axios from 'axios'
// axios.defaults.withCredentials = true
export default {
  async ping (baseUrl) {
    try {
      const response = await axios({
        url: '/alive',
        baseURL: baseUrl,
        method: 'Get'
      })
      return response
    } catch (e) {
      console.log(e.response) // undefined
    }
  },
  async play (baseUrl, item) {
    const response = await axios({
      url: '/file/' + item.file_id + '/download/' + item.title,
      baseURL: baseUrl,
      method: 'Get'
    })
    return response
  }
}
