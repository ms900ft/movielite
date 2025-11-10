<script setup>
import { RouterView } from 'vue-router';
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { authService } from './services/auth.js';

const router = useRouter();
const isAuthenticated = ref(false);

const checkAuth = () => {
  isAuthenticated.value = authService.isAuthenticated();
};

const logout = () => {
  authService.logout();
  isAuthenticated.value = false;
  router.push('/login');
};

onMounted(() => {
  checkAuth();
});
</script>

<template>
  <div id="app">
    <header v-if="isAuthenticated">
      <h1>Movie Database</h1>
      <nav>
        <router-link to="/">Movies</router-link>
        <button @click="logout" class="logout-btn">Logout</button>
      </nav>
    </header>

    <main>
      <RouterView />
    </main>
  </div>
</template>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  min-height: 100vh;
  background-color: #f8f9fa;
}

header {
  background-color: #42b883;
  color: white;
  padding: 1rem;
  text-align: center;
}

header h1 {
  margin: 0 0 1rem 0;
  font-size: 2rem;
}

nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 2rem;
}

nav a {
  color: white;
  text-decoration: none;
  font-weight: bold;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  transition: background-color 0.3s;
}

nav a:hover,
nav a.router-link-active {
  background-color: rgba(255, 255, 255, 0.2);
}

main {
  min-height: calc(100vh - 120px);
  background-color: #f8f9fa;
}
</style>
