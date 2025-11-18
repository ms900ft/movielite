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
              :src="`http://localhost:8001/images/w342${movie.meta.poster_path}`"
              :alt="movie.title"
              @error="handleImageError"
            />
            <div v-else class="no-poster">{{ movie.title }}</div>
            <div class="star-icon" :class="{ 'watchlist-star': movie.watchlist }" @click.stop="toggleWatchlist(movie)">★</div>
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
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { debounce } from 'lodash-es';
import { moviesService } from '../services/movies.js';

const router = useRouter();
const route = useRoute();

const movies = ref([]);
const loading = ref(true);
const error = ref(null);
const currentOffset = ref(0);
const hasMore = ref(true);
const loadingMore = ref(false);
const limit = 20;
const searchQuery = ref(route.query.q || '');
const currentSearch = ref('');
const currentPerson = ref(null);

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
    }
    const response = await moviesService.getMovies(params);
    const newMovies = response.data || [];
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

const goToMovieDetail = (movieId) => {
  router.push(`/movie/${movieId}`);
};

const playMovie = async (movieId) => {
  try {
    await moviesService.playMovie(movieId);
    alert('Movie playback started!');
  } catch (err) {
    console.error('Error playing movie:', err);
    alert('Failed to start movie playback.');
  }
};

const toggleWatchlist = async (movie) => {
  try {
    const updatedMovie = { ...movie, watchlist: !movie.watchlist };
    await moviesService.updateMovie(movie.id, updatedMovie);
    movie.watchlist = !movie.watchlist;
  } catch (err) {
    console.error('Error toggling watchlist:', err);
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
  window.addEventListener('scroll', handleScroll);

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
});
</script>

<style scoped>
.movies-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
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
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 15px;
  width: 1000px;
}

@media (max-width: 768px) {
  .movies-grid {
    gap: 10px;
  }
}

@media (min-width: 1200px) {
  .movies-grid {
    gap: 20px;
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

.movie-poster {
  height: 225px;
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

.star-icon {
  position: absolute;
  top: 0px;
  left: 5px;
  color: gray;
  font-size: 26px;
  z-index: 10;
}

.watchlist-star {
  color: #ffd700;
}

.no-movies {
  text-align: center;
  padding: 40px;
  color: #666;
  font-size: 18px;
}
</style>