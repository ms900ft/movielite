import axios from 'axios'
// axios.defaults.withCredentials = true
export default {
  async ping (baseUrl) {
    const response = await axios({
      url: '/alive',
      baseURL: baseUrl,
      method: 'Get'
    })
    return response
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
