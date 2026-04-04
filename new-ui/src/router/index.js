import { createRouter, createWebHistory } from 'vue-router';
import Movies from '../views/Movies.vue';
import MovieDetail from '../views/MovieDetail.vue';
import Login from '../views/Login.vue';
import Users from '../views/Users.vue';
import { authService } from '../services/auth.js';

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/',
    name: 'Movies',
    component: Movies,
    meta: { requiresAuth: true }
  },
  {
    path: '/movies',
    name: 'MoviesList',
    component: Movies,
    meta: { requiresAuth: true }
  },
  {
    path: '/movie/:id',
    name: 'MovieDetail',
    component: MovieDetail,
    meta: { requiresAuth: true }
  },
  {
    path: '/users',
    name: 'Users',
    component: Users,
    meta: { requiresAuth: true }
  }
];

const router = createRouter({
  history: createWebHistory('/movie2'),
  routes
});

// Navigation guard to check authentication
router.beforeEach((to) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!authService.isAuthenticated()) {
      return '/login';
    }
  }
});

export default router;