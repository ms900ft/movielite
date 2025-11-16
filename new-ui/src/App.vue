<script setup>
import { RouterView } from 'vue-router';
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { authService } from './services/auth.js';
import { moviesService } from './services/movies.js';
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';
import InputText from 'primevue/inputtext';
import Menubar from 'primevue/menubar';

const router = useRouter();
const isAuthenticated = ref(false);
const searchQuery = ref('');
const genres = ref([]);
const countries = ref([]);

const menuItems = ref([
  {
    label: 'Movies',
    icon: 'pi pi-video',
    command: () => {
      router.push('/');
    }
  },

  {
    label: 'Genres',
    icon: 'pi pi-tags',
    items: genres
  },
  {
    label: 'Countries',
    icon: 'pi pi-globe',
    items: countries
  },
   {
     label: 'Watchlist',
     command: () => {
       router.push({ path: '/', query: { search: 'watchlist' } });
     }
   },
  {
    label: 'Settings',
    icon: 'pi pi-cog'
  }

]);

// Update menu items when genres and countries change
import { watch } from 'vue';
watch(genres, () => {
  menuItems.value[1].items = genres.value;
}, { deep: true });

watch(countries, () => {
  menuItems.value[2].items = countries.value;
}, { deep: true });

const fetchGenres = async () => {
  try {
    const response = await moviesService.getGenres();
    const data = response.data || response || [];
    genres.value = data.map(genre => ({
      label: genre.name,
      icon: 'pi pi-tag',
      command: () => {
        router.push({ path: '/', query: { genre: genre.tmdb_id, country: '', orderby: 'name' } });
      }
    }));
  } catch (error) {
    console.error('Error fetching genres:', error);
  }
};

const fetchCountries = async () => {
  try {
    const response = await moviesService.getCountries();
    const data = response.data || response || [];
    countries.value = data.map(country => ({
      label: country.name,
      icon: 'pi pi-map-marker',
      command: () => {
        router.push({ path: '/', query: { country: country.iso_id, orderby: 'name' } });
      }
    }));
  } catch (error) {
    console.error('Error fetching countries:', error);
  }
};

const onSearchInput = () => {
  // Only search if 3+ characters or empty (clear search)
  if (searchQuery.value.length >= 3 || searchQuery.value.length === 0) {
    router.push({ path: '/', query: { q: searchQuery.value } });
  }
};

const onSearchKeyup = (event) => {
  if (event.key === 'Enter') {
    router.push({ path: '/', query: { q: searchQuery.value } });
  }
};

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
  fetchGenres();
  fetchCountries();
});
</script>

<template>
  <div id="app">
    <main>
    <div v-if="isAuthenticated" class="w-full">
      <Menubar :model="menuItems" class="w-full mb-2">
        <template #end>
          <IconField>
            <InputIcon>
              <i class="pi pi-search" />
            </InputIcon>
            <InputText v-model="searchQuery" placeholder="Search movies..." @input="onSearchInput" @keyup="onSearchKeyup" />
          </IconField>
        </template>
      </Menubar>

    </div>


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
