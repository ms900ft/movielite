<template>
  <div class="movies-container">
    <h1>Movies</h1>
    <div class="toolbar">
      <input v-model="searchQuery" type="text" placeholder="Search movies..." @input="onSearchInput" />
      <button @click="searchMovies">Search</button>
    </div>
    <div v-if="loading" class="loading">Loading movies...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else>
      <div class="movies-grid">
        <div v-for="movie in movies" :key="movie.id" class="movie-item">
          <div class="movie-poster">
            <img
              v-if="movie.meta && movie.meta.poster_path"
              :src="`http://localhost:8001/images/w342${movie.meta.poster_path}`"
              :alt="movie.title"
              @error="handleImageError"
            />
            <div v-else class="no-poster">{{ movie.title }}</div>
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
import { ref, onMounted, onUnmounted } from 'vue';
import { debounce } from 'lodash-es';
import { moviesService } from '../services/movies.js';

const movies = ref([]);
const loading = ref(true);
const error = ref(null);
const currentOffset = ref(0);
const hasMore = ref(true);
const loadingMore = ref(false);
const limit = 20;
const searchQuery = ref('');

const fetchMovies = async (offset = 0) => {
  try {
    if (offset === 0) {
      loading.value = true;
      error.value = null;
    } else {
      loadingMore.value = true;
    }
    const params = { limit, offset };
    if (searchQuery.value) {
      params.title = searchQuery.value;
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
  fetchMovies(0);
};

const debouncedSearch = debounce(() => {
  if (searchQuery.value.length >= 3 || searchQuery.value.length === 0) {
    searchMovies();
  }
}, 300);

const onSearchInput = () => {
  debouncedSearch();
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


.no-movies {
  text-align: center;
  padding: 40px;
  color: #666;
  font-size: 18px;
}
</style>