import Vue from 'vue'
import vuetify from './plugins/vuetify'

import App from './App.vue'
import store from './store'
import { router } from './router'
import axios from 'axios'
// import VueLodash from 'vue-lodash'
import lodash from 'lodash'
import MobileDetect from 'mobile-detect'
import VeeValidate from 'vee-validate'

import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import '@mdi/font/css/materialdesignicons.css'

import './assets/main.css'

Vue.prototype._ = lodash

Vue.config.productionTip = false
let BaseUrl = location.protocol

if (process.env.NODE_ENV === 'development') {
  BaseUrl = 'http://localhost:8001'
}

axios.defaults.baseURL = BaseUrl
Vue.prototype.$baseURL = BaseUrl
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

Vue.use(VeeValidate)
Vue.component('font-awesome-icon', FontAwesomeIcon)
new Vue({
  vuetify,
  render: h => h(App),
  store,
  router
}).$mount('#app')
