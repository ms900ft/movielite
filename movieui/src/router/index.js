import Vue from 'vue'
import VueRouter from 'vue-router'
import ListMovie from '@/components/ListMovie'

Vue.use(VueRouter)

export default new VueRouter({
  routes: [{
    path: '/',
    name: 'ListMovie',
    props: (route) => ({
      show: route.query.show,
      orderby: route.query.orderby,
      genre: route.query.genre,
      country: route.query.country,
      cast: route.query.cast,
      crew: route.query.crew,
      person: route.query.person
    }),
    component: ListMovie
  }
  ],
  mode: 'history',
  base: '/movie2'
})
