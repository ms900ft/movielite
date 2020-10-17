import Vue from 'vue'
import VueRouter from 'vue-router'
import LatestMovie from '@/components/LatestMovie'

Vue.use(VueRouter)

export default new VueRouter({
  routes: [{
    path: '/',
    name: 'LatestMovie',
    props: (route) => ({
      show: route.query.show,
      orderby: route.query.orderby,
      genre: route.query.genre,
      country: route.query.country,
      cast: route.query.cast,
      crew: route.query.crew,
      person: route.query.person
    }),
    component: LatestMovie
  }],
  mode: 'history',
  base: '/movie2'
})
