<script setup>
import { RouterView } from 'vue-router';
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { authService } from './services/auth.js';
import { moviesService } from './services/movies.js';
import { useMovieStore } from './stores/movie.js';
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';
import InputText from 'primevue/inputtext';
import Menubar from 'primevue/menubar';

const router = useRouter();
const movieStore = useMovieStore();
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
      icon: 'pi pi-heart',
      command: () => {
        router.push({ path: '/', query: { show: 'watchlist' } });
      }
    },
   {
      label: 'Recently',
      icon: 'pi pi-clock',
      command: () => {
        router.push({ path: '/', query: { orderby: 'recent' } });
      }
    },
  {
    label: 'Settings',
    icon: 'pi pi-cog',
    items: [
      {
        label: 'Users',
        icon: 'pi pi-users',
        command: () => {
          router.push('/users');
        }
      }
    ]
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

const keyboardVisible = ref(false);
const shiftActive = ref(false);
const searchInput = ref();

const toggleKeyboard = () => {
  keyboardVisible.value = !keyboardVisible.value;
  if (keyboardVisible.value) {
    searchInput.value?.$el?.focus();
  }
};

const toggleShift = () => {
  shiftActive.value = !shiftActive.value;
};

const typeKey = (key) => {
  const char = shiftActive.value && key.length === 1 ? key.toUpperCase() : key;
  searchQuery.value += char;
  shiftActive.value = false;
  onSearchInput();
  searchInput.value?.$el?.focus();
};

const backspace = () => {
  searchQuery.value = searchQuery.value.slice(0, -1);
  onSearchInput();
  searchInput.value?.$el?.focus();
};

const searchFromKeyboard = () => {
  router.push({ path: '/', query: { q: searchQuery.value } });
  keyboardVisible.value = false;
};

const onSearchFocus = () => {
  // keyboard stays visible on focus
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
    <div v-if="isAuthenticated" class="menubar-container">
      <Menubar :model="menuItems" class="mb-2">
        <template #end>
          <div class="end-slot">
            <IconField>
              <InputIcon>
                <i class="pi pi-search" />
              </InputIcon>
              <InputText ref="searchInput" v-model="searchQuery" placeholder="Search movies..." @input="onSearchInput" @keyup="onSearchKeyup" @focus="onSearchFocus" />
              <InputIcon class="keyboard-toggle" @click="toggleKeyboard">
                <i class="pi pi-desktop" />
              </InputIcon>
            </IconField>
            <span v-if="movieStore.totalResults > 0" class="results-count">{{ movieStore.totalResults }} results</span>
          </div>
        </template>
      </Menubar>

      <!-- On-screen keyboard -->
      <div v-if="keyboardVisible" class="keyboard-overlay">
        <div class="keyboard" @mousedown.prevent>
          <div class="keyboard-row">
            <button v-for="key in ['1','2','3','4','5','6','7','8','9','0']" :key="key" @click="typeKey(key)" class="key">{{ key }}</button>
          </div>
          <div class="keyboard-row">
            <button v-for="key in ['q','w','e','r','t','z','u','i','o','p','ü']" :key="key" @click="typeKey(key)" class="key">{{ shiftActive ? key.toUpperCase() : key }}</button>
          </div>
          <div class="keyboard-row">
            <button v-for="key in ['a','s','d','f','g','h','j','k','l','ö','ä']" :key="key" @click="typeKey(key)" class="key">{{ shiftActive ? key.toUpperCase() : key }}</button>
          </div>
          <div class="keyboard-row">
            <button @click="toggleShift" :class="['key', 'key-mod', { active: shiftActive }]">⇧</button>
            <button v-for="key in ['y','x','c','v','b','n','m']" :key="key" @click="typeKey(key)" class="key">{{ shiftActive ? key.toUpperCase() : key }}</button>
            <button @click="typeKey('ß')" class="key">ß</button>
            <button @click="backspace" class="key key-backspace">⌫</button>
          </div>
          <div class="keyboard-row">
            <button @click="typeKey(' ')" class="key key-space">Space</button>
            <button @click="searchFromKeyboard" class="key key-search">Search</button>
            <button @click="toggleKeyboard" class="key key-close">✕</button>
          </div>
        </div>
      </div>

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

/* Ensure PrimeVue menubar dropdowns above card overlays */
.p-menubar-sublist {
  z-index: 9999 !important;
}

.menubar-container {
  width: 100%;
}

.menubar-container .p-menubar {
  width: 100%;
}

.end-slot {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  position: relative;
}

.results-count {
  margin-left: 0;
  font-weight: 600;
  color: #64748b;
  white-space: nowrap;
}

.keyboard-toggle {
  cursor: pointer;
  opacity: 0.6;
  transition: opacity 0.2s;
}

.keyboard-toggle:hover {
  opacity: 1;
}

.keyboard-overlay {
  position: fixed;
  top: 60px;
  right: 20px;
  padding: 8px;
  z-index: 10000;
}

.keyboard {
  max-width: 800px;
  margin: 0 auto;
  background: rgba(30, 41, 59, 0.6);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border-radius: 12px;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.keyboard-row {
  display: flex;
  justify-content: center;
  gap: 4px;
}

.key {
  min-width: 40px;
  height: 40px;
  border: none;
  border-radius: 6px;
  background: #334155;
  color: #f1f5f9;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.15s;
  display: flex;
  align-items: center;
  justify-content: center;
  flex: 1;
  max-width: 50px;
}

.key:hover {
  background: #475569;
}

.key:active {
  background: #64748b;
}

.key-mod {
  max-width: 60px;
  background: #475569;
}

.key-mod.active {
  background: #3b82f6;
}

.key-backspace {
  max-width: 70px;
  background: #475569;
}

.key-space {
  max-width: 250px;
  flex: 3;
}

.key-search {
  max-width: 80px;
  background: #3b82f6;
}

.key-search:hover {
  background: #2563eb;
}

.key-close {
  max-width: 50px;
  background: #ef4444;
}

.key-close:hover {
  background: #dc2626;
}

</style>
