import { createRouter, createWebHistory } from 'vue-router';
import Movies from '../views/Movies.vue';
import MovieDetail from '../views/MovieDetail.vue';
import Login from '../views/Login.vue';
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
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

// Navigation guard to check authentication
router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!authService.isAuthenticated()) {
      next('/login');
    } else {
      next();
    }
  } else {
    next();
  }
});

export default router;