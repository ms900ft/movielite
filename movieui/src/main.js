import Vue from 'vue'
import './plugins/vuetify'
import App from './App.vue'
import store from './store'
import router from './router'
import axios from 'axios'
// import VueLodash from 'vue-lodash'
import lodash from 'lodash'
import MobileDetect from 'mobile-detect'
import './assets/main.css'

Vue.prototype._ = lodash

Vue.config.productionTip = false
let ApiUrl = 'http://' + location.host
if (process.env.NODE_ENV === 'development') {
  ApiUrl = 'http://localhost:8001'
}

axios.defaults.baseURL = ApiUrl
Vue.prototype.$baseURL = ApiUrl
Vue.prototype.$localViewURL = 'http://localhost:8081'
Vue.prototype.$hitspp = 30
function isMobile () {
  var md = new MobileDetect(window.navigator.userAgent)
  if (md.mobile()) {
    return true
  }
  return false
};
const mobile = isMobile()
Vue.prototype.$isMobile = mobile
Vue.filter('truncate', function (text, stop, clamp) {
  if (stop < 0) {
    return text
  } else if (stop === 0) {
    return ''
  }
  if (!text) {
    return ''
  }
  return text.substr(0, text.lastIndexOf(' ', stop)) + ' ...'
})

new Vue({
  render: h => h(App),
  store,
  router
}).$mount('#app')
