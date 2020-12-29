import axios from 'axios'

// const API_URL = this.baseURL

class AuthService {
  login (user) {
    return axios
      .post('/login', {
        username: user.username,
        password: user.password
      })
      .then(response => {
        if (response.data.token) {
          localStorage.setItem('user', JSON.stringify(response.data))
        }

        return response.data
      })
  }

  logout () {
    localStorage.removeItem('user')
  }
}

export default new AuthService()
