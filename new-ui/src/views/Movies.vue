<template>
  <div class="movies-container">
    <h1>Movies</h1>
    <div v-if="currentSearch" class="current-search">{{ currentSearch }}</div>
    <div v-if="loading" class="loading">Loading movies...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else>
      <div class="movies-grid">
        <div v-for="movie in movies" :key="movie.id" class="movie-item" @click="goToMovieDetail(movie.id)">
          <div class="movie-poster">
            <img
              v-if="movie.meta && movie.meta.poster_path"
              :src="`/images/w342${movie.meta.poster_path}`"
              :alt="movie.title"
              @error="handleImageError"
            />
            <div v-else class="no-poster">{{ movie.title }}</div>
            <div v-if="openMenuMovieID !== movie.id" class="star-icon" :class="{ 'watchlist-star': movie.watchlist }" @click.stop="toggleWatchlist(movie)" aria-label="watchlist" title="Toggle watchlist">★</div>
            <button v-if="movie.multiplechoice" class="multiple-choice-button" @click.stop="showMultipleChoice(movie)" title="Choose correct movie">
              <i class="pi pi-search"></i> Find
            </button>
            <button class="menu-button" @click.stop="toggleMenu($event, movie)">
              <i class="pi pi-bars"></i>
            </button>
            <div class="play-button-overlay" @click.stop="playMovie(movie.id)">
              <img src="https://www.freeiconspng.com/uploads/play-button-icon-png-0.png" alt="play" style="width: 40px; height: 40px;" />
            </div>
            <div class="movie-title-overlay">{{ movie.title }}</div>
          </div>
        </div>
      </div>
      <div v-if="loadingMore" class="loading-more">Loading more movies...</div>
      <div v-if="movies.length === 0" class="no-movies">
        No movies found.
      </div>
    </div>

    <!-- Menu component -->
    <Menu ref="movieMenu" :model="menuItems" :popup="true" />

    <!-- Move Movie Dialog -->
    <Dialog v-model:visible="moveDialogVisible" modal header="Move Movie" :style="{ width: '50rem' }">
      <p>Select a directory to move the movie to:</p>
      <div class="target-list">
        <button v-for="target in targets" :key="target.name" @click="moveMovie(selectedMovie, target.name)" class="target-button">
          {{ target.name }}
        </button>
      </div>
    </Dialog>

    <!-- Multiple Choice Dialog -->
    <Dialog v-model:visible="multipleChoiceVisible" modal header="Select the correct movie" :style="{ width: '60rem' }">
      <div class="multiple-choice-grid">
        <div v-for="result in selectedMovie?.multiplechoice?.Results" :key="result.ID" class="multiple-choice-item" @click="selectMultipleChoice(result)">
          <img v-if="result.poster_path" :src="`/images/w185${result.poster_path}`" :alt="result.title" @error="result.poster_path = null" />
          <div v-if="!result.poster_path" class="no-poster-small">
            <i class="pi pi-image" style="font-size: 40px; opacity: 0.5;"></i>
          </div>
          <div class="multiple-choice-title">{{ result.title }}</div>
          <div class="multiple-choice-date">{{ result.release_date }}</div>
          <a :href="`https://www.themoviedb.org/movie/${result.ID}`" target="_blank" rel="noopener" class="tmdb-link" @click.stop>
            <i class="pi pi-external-link"></i> View on TMDB
          </a>
        </div>
      </div>
    </Dialog>

    <!-- Streaming Video Player Dialog -->
    <Dialog v-model:visible="streamingVisible" modal header="Now Playing" :style="{ width: '80vw', maxWidth: '1200px' }" @hide="streamingSrc = ''">
      <video v-if="streamingSrc" :src="streamingSrc" controls style="width: 100%; max-height: 70vh;">
        Your browser does not support video playback.
      </video>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { debounce } from 'lodash-es';
import { moviesService } from '../services/movies.js';
import { useMovieStore } from '../stores/movie.js';
import Menu from 'primevue/menu';
import Dialog from 'primevue/dialog';

const router = useRouter();
const route = useRoute();
const movieStore = useMovieStore();

const movies = ref([]);
const loading = ref(true);
const error = ref(null);
const currentOffset = ref(0);
const hasMore = ref(true);
const loadingMore = ref(false);
const limit = 40;
const searchQuery = ref(route.query.q || '');
const currentSearch = ref('');
const currentPerson = ref(null);

// Menu related
const movieMenu = ref();
// Track which movie's popup menu is currently open to avoid overlay conflicts
const openMenuMovieID = ref(null);
const selectedMovie = ref(null);
const targets = ref([]);
const moveDialogVisible = ref(false);
const multipleChoiceVisible = ref(false);
const streamingVisible = ref(false);
const streamingSrc = ref('');

// Menu items ref (updated dynamically when menu opens)
const menuItems = ref([]);

const updateMenuItems = () => {
  menuItems.value = [
    {
      label: 'Stream',
      icon: 'pi pi-play',
      command: () => {
        const streamUrl = moviesService.getStreamUrl(selectedMovie.value.id);
        streamingSrc.value = streamUrl;
        streamingVisible.value = true;
      }
    },
    {
      label: 'Play (Server)',
      icon: 'pi pi-desktop',
      command: () => playMovie(selectedMovie.value.id, true)
    },
    {
      label: 'Download',
      icon: 'pi pi-download',
      command: async () => {
        const movie = selectedMovie.value;
        const token = localStorage.getItem('authToken');
        const url = moviesService.getDownloadUrl(movie.id);
        window.open(`${url}?token=${token}`, '_blank');
      }
    },
    {
      label: selectedMovie.value?.watchlist ? 'Remove from Watchlist' : 'Add to Watchlist',
      icon: 'pi pi-star',
      command: () => toggleWatchlist(selectedMovie.value)
    },
    {
      label: 'Move Movie',
      icon: 'pi pi-folder',
      command: () => showMoveDialog(selectedMovie.value)
    },
    {
      label: 'View Details',
      icon: 'pi pi-info-circle',
      command: () => goToMovieDetail(selectedMovie.value.id)
    },
    {
      label: 'Show',
      icon: 'pi pi-folder-open',
      command: () => showMovie(selectedMovie.value.id)
    },
    {
      label: 'Rescan',
      icon: 'pi pi-refresh',
      command: () => rescanMovieMeta(selectedMovie.value)
    }
  ];
};

const setCurrentSearch = async () => {
  currentPerson.value = null;
  if (route.query.country) {
    try {
      const countries = await moviesService.getCountries();
      const country = countries.find(c => c.iso_id === route.query.country);
      currentSearch.value = `Country: ${country ? country.name : route.query.country}`;
    } catch (e) {
      currentSearch.value = `Country: ${route.query.country}`;
    }
  } else if (route.query.genre) {
    try {
      const genres = await moviesService.getGenres();
      const genre = genres.find(g => g.tmdb_id == route.query.genre);
      currentSearch.value = `Genre: ${genre ? genre.name : route.query.genre}`;
    } catch (e) {
      currentSearch.value = `Genre: ${route.query.genre}`;
    }
  } else if (route.query.person) {
    try {
      const person = await moviesService.getPerson(route.query.person);
      currentSearch.value = `Person: ${person.Name}`;
      currentPerson.value = person;
    } catch (e) {
      currentSearch.value = `Person: ${route.query.person}`;
    }
  } else if (route.query.show === 'watchlist') {
    currentSearch.value = 'Watchlist';
  } else if (searchQuery.value) {
    currentSearch.value = `Search: "${searchQuery.value}"`;
  } else {
    currentSearch.value = '';
  }
};

const fetchTargets = async () => {
  try {
    const response = await moviesService.getTargets();
    targets.value = response.data || response || [];
  } catch (error) {
    console.error('Error fetching targets:', error);
  }
};

const fetchMovies = async (offset = 0) => {
  await setCurrentSearch();
  try {
    if (offset === 0) {
      loading.value = true;
      error.value = null;
    } else {
      loadingMore.value = true;
    }
    const params = { limit, offset };

    // Set default ordering when country or genre is selected
    if (route.query.country || route.query.genre) {
      if (route.query.country) {
        params.country = route.query.country;
      }
      if (route.query.genre) {
        params.genre = route.query.genre;
      }
      params.orderby = route.query.orderby || 'name';
    } else {
      // Normal pagination for other cases
      if (searchQuery.value) {
        params.title = searchQuery.value;
      }
      // Check for person query parameter
      if (route.query.person) {
        params.person = route.query.person;
      }
      // Check for show query parameter
      if (route.query.show === 'watchlist') {
        params.show = 'watchlist';
      }
      // Check for orderby parameter
      if (route.query.orderby) {
        params.orderby = route.query.orderby;
      }
    }
    const response = await moviesService.getMovies(params);
    const newMovies = response.data || [];
    movieStore.setTotalResults(response.meta?.total || 0);
    if (offset === 0) {
      movies.value = newMovies;
    } else {
      movies.value.push(...newMovies);
    }
    hasMore.value = newMovies.length === limit;
    currentOffset.value = offset;
  } catch (err) {
    error.value = 'Failed to load movies. Please try again.';
    console.error('Error fetching movies:', err);
  } finally {
    loading.value = false;
    loadingMore.value = false;
  }
};

const loadMore = () => {
  if (!loadingMore.value && hasMore.value) {
    fetchMovies(currentOffset.value + limit);
  }
};

const searchMovies = () => {
  currentOffset.value = 0;
  hasMore.value = true;
  // Clear person filter when searching
  router.replace({ query: {} });
  fetchMovies(0);
};

const debouncedSearch = debounce(() => {
  // Debounced search only triggers for 3+ characters, but immediate search for clearing
  if (searchQuery.value.length === 0) {
    searchMovies();
  }
}, 300);

const onSearchInput = () => {
  debouncedSearch();
};

// Watch for query parameter changes
import { watch } from 'vue';
watch(() => route.query.q, async (newQuery) => {
  if (newQuery !== searchQuery.value) {
    searchQuery.value = newQuery || '';
    await setCurrentSearch();
    fetchMovies(0);
  }
});

watch(() => route.query.genre, async (newGenre) => {
  // Filter by genre - always reload
  await setCurrentSearch();
  fetchMovies(0);
});

watch(() => route.query.country, async (newCountry) => {
  // Filter by country - always reload
  await setCurrentSearch();
  fetchMovies(0);
});

watch(() => route.query.show, async () => {
  // Filter by show - always reload
  await setCurrentSearch();
  fetchMovies(0);
});

watch(() => route.query.orderby, async () => {
  // Filter by orderby - always reload
  await setCurrentSearch();
  fetchMovies(0);
});

const goToMovieDetail = (movieId) => {
  router.push(`/movie/${movieId}`);
};

const playMovie = async (movieId, forceServerPlay = false) => {
  const isLocalhost = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1';
  const shouldServerPlay = forceServerPlay || isLocalhost;
  
  if (shouldServerPlay) {
    try {
      await moviesService.playMovie(movieId);
    } catch (err) {
      console.error('Error playing movie:', err);
    }
  } else {
    const streamUrl = moviesService.getStreamUrl(movieId);
    streamingSrc.value = streamUrl;
    streamingVisible.value = true;
  }
};

const showMovie = async (movieId) => {
  try {
    await moviesService.showMovie(movieId);
  } catch (err) {
    console.error('Error showing movie:', err);
  }
};

const toggleWatchlist = async (movie) => {
  const wasInWatchlist = movie.watchlist;
  try {
    const updatedMovie = { ...movie, watchlist: !movie.watchlist };
    await moviesService.updateMovie(movie.id, updatedMovie);
    movie.watchlist = !movie.watchlist;
    // Remove from view if removed from watchlist and on watchlist page
    if (wasInWatchlist && route.query.show === 'watchlist') {
      movies.value = movies.value.filter(m => m.id !== movie.id);
    }
  } catch (err) {
    console.error('Error toggling watchlist:', err);
  }
};

const moveMovie = async (movie, targetDir) => {
  try {
    await moviesService.moveFile(movie.File.id, targetDir);
    alert(`Movie moved to ${targetDir}`);
    // Optionally refresh the list or remove the movie from the list
    movies.value = movies.value.filter(m => m.id !== movie.id);
    moveDialogVisible.value = false;
  } catch (err) {
    console.error('Error moving movie:', err);
    alert('Failed to move movie.');
  }
};

// Method to toggle the menu
const toggleMenu = (event, movie) => {
  selectedMovie.value = movie;
  openMenuMovieID.value = movie.id;
  updateMenuItems();
  movieMenu.value.toggle(event);
};

// Close menu when clicking outside the movie item
const handleDocumentClick = (e) => {
  if (openMenuMovieID.value != null) {
    const inside = e.target.closest('.movie-item') != null;
    if (!inside) {
      openMenuMovieID.value = null;
      if (movieMenu.value && typeof movieMenu.value.hide === 'function') {
        movieMenu.value.hide();
      }
    }
  }
};

const showMoveDialog = (movie) => {
  selectedMovie.value = movie;
  moveDialogVisible.value = true;
};

const showMultipleChoice = (movie) => {
  selectedMovie.value = movie;
  multipleChoiceVisible.value = true;
};

const selectMultipleChoice = async (result) => {
  try {
    await moviesService.rescanMovie(selectedMovie.value.id, result.ID);
    fetchMovies(0);
    multipleChoiceVisible.value = false;
  } catch (err) {
    console.error('Error selecting movie:', err);
    alert('Failed to update movie metadata.');
  }
};

const rescanMovies = async () => {
  console.log('Rescan triggered');
  try {
    await moviesService.rescan();
    fetchMovies(0);
  } catch (err) {
    console.error('Error rescanning:', err);
  }
};

const rescanMovieMeta = async (movie) => {
  try {
    const metaId = movie.meta?.ID || movie.TMDBMovieID;
    if (!metaId) {
      alert('No TMDB metadata found for this movie');
      return;
    }
    await moviesService.rescanMovie(movie.id, metaId);
    fetchMovies(0);
  } catch (err) {
    console.error('Error rescanning movie metadata:', err);
  }
};

const handleScroll = () => {
  const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
  const windowHeight = window.innerHeight;
  const documentHeight = document.documentElement.scrollHeight;
  if (scrollTop + windowHeight >= documentHeight - 100) {
    loadMore();
  }
};

const handleImageError = (event) => {
  event.target.style.display = 'none';
  event.target.nextElementSibling.style.display = 'block';
};

onMounted(() => {
  fetchMovies();
  fetchTargets();
  window.addEventListener('scroll', handleScroll);

  document.addEventListener('click', handleDocumentClick);

  // Restore scroll position if coming back from movie detail
  const savedScrollY = sessionStorage.getItem('movieDetailScrollY');
  if (savedScrollY) {
    nextTick(() => {
      window.scrollTo(0, parseInt(savedScrollY));
      sessionStorage.removeItem('movieDetailScrollY');
    });
  }
});

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll);
  document.removeEventListener('click', handleDocumentClick);
});
</script>

<style scoped>
.movies-container {
  padding: 20px 60px;
  width: 100%;
  box-sizing: border-box;
}

.nav-bar {
  margin-bottom: 20px;
}

.nav-link {
  color: #007bff;
  text-decoration: none;
  font-size: 18px;
  font-weight: bold;
}

.nav-link:hover {
  text-decoration: underline;
}

.current-search {
  margin-bottom: 20px;
  font-size: 18px;
  color: #333;
  font-weight: bold;
}

.toolbar {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.toolbar input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 16px;
}

.toolbar button {
  padding: 8px 16px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
}

.toolbar button:hover {
  background-color: #0056b3;
}

.loading, .error, .loading-more {
  text-align: center;
  padding: 20px;
  font-size: 18px;
}

.error {
  color: red;
}

.loading-more {
  color: #666;
}

.movies-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 40px;
  width: 100%;
}

@media (max-width: 768px) {
  .movies-grid {
    gap: 10px;
  }
}

@media (min-width: 1200px) {
  .movies-grid {
    gap: 40px;
  }
}

.movie-item {
  position: relative;
  overflow: hidden;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  transition: transform 0.2s, box-shadow 0.2s;
  background: white;
}

.movie-item:hover {
  transform: translateY(-3px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
}

.movie-item:hover .movie-poster img {
  opacity: 0.7;
}

.movie-poster {
  height: 300px;
  overflow: hidden;
  position: relative;
  cursor: pointer;
}

.movie-poster img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.movie-title-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.8), transparent);
  color: white;
  padding: 10px;
  font-size: 14px;
  font-weight: bold;
  text-align: center;
  opacity: 0;
  transition: opacity 0.3s ease;
  pointer-events: none;
}

.movie-item:hover .movie-title-overlay {
  opacity: 1;
}

.no-poster {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  font-size: 14px;
  font-weight: bold;
  text-align: center;
  padding: 10px;
  box-sizing: border-box;
  word-wrap: break-word;
  line-height: 1.2;
}

.play-button-overlay {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  opacity: 0;
  transition: opacity 0.3s ease;
  background: transparent;
  border-radius: 50%;
  width: 100px;
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: black;
  font-size: 20px;
  cursor: pointer;
}

.movie-item:hover .play-button-overlay {
  opacity: 1;
}

/* star-icon overlay for per-movie watchlist */
.star-icon {
  position: absolute;
  top: 6px;
  left: 6px;
  color: gray;
  font-size: 26px;
  z-index: 0;
  cursor: pointer;
}

.watchlist-star {
  color: #ffd700;
}

/* Force PrimeVue popup menu above all card overlays */
.p-menu-overlay {
  z-index: 9999 !important;
}

.menu-button {
  position: absolute;
  top: 5px;
  right: 5px;
  background: transparent;
  color: black;
  border: none;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 0;
}

.menu-button:hover {
  background: rgba(0, 0, 0, 0.1);
}

.multiple-choice-button {
  position: absolute;
  bottom: 35px;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(0, 0, 0, 0.7);
  color: white;
  border: none;
  padding: 6px 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  cursor: pointer;
  z-index: 2;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
}

.multiple-choice-button:hover {
  background: rgba(0, 0, 0, 0.9);
}

.multiple-choice-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 16px;
}

.multiple-choice-item {
  cursor: pointer;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  transition: transform 0.2s, box-shadow 0.2s;
  background: white;
}

.multiple-choice-item:hover {
  transform: translateY(-3px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
}

.multiple-choice-item img {
  width: 100%;
  height: 210px;
  object-fit: cover;
}

.no-poster-small {
  width: 100%;
  height: 210px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-sizing: border-box;
}

.multiple-choice-title {
  padding: 8px;
  font-size: 13px;
  font-weight: 600;
  text-align: center;
}

.multiple-choice-date {
  padding: 0 8px 8px;
  font-size: 12px;
  color: #64748b;
  text-align: center;
}

.tmdb-link {
  display: block;
  padding: 6px 8px 8px;
  text-align: center;
  color: #3b82f6;
  font-size: 12px;
  text-decoration: none;
  border-top: 1px solid #e2e8f0;
  transition: background 0.2s;
}

.tmdb-link:hover {
  background: #f1f5f9;
  text-decoration: underline;
}

.no-movies {
  text-align: center;
  padding: 40px;
  color: #666;
  font-size: 18px;
}

.target-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.target-button {
  padding: 10px 15px;
  background-color: #f8f9fa;
  border: 1px solid #dee2e6;
  border-radius: 4px;
  cursor: pointer;
  text-align: left;
  font-size: 16px;
}

.target-button:hover {
  background-color: #e9ecef;
}


</style>
