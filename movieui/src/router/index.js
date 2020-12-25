import Vue from 'vue'
import VueRouter from 'vue-router'
import ListMovie from '@/components/ListMovie'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'

Vue.use(VueRouter)

export const router = new VueRouter({
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
  },
  {
    path: '/login',
    component: Login
  },
  {
    path: '/register',
    component: Register
  }
  ],
  mode: 'history',
  base: '/movie2'
})

router.beforeEach((to, from, next) => {
  const publicPages = ['/login', '/register', '/home']
  const authRequired = !publicPages.includes(to.path)
  const loggedIn = localStorage.getItem('user')
  // trying to access a restricted page + not logged in
  // redirect to login page
  if (authRequired && !loggedIn) {
    next('/login')
  } else {
    next()
  }
})
