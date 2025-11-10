<template>
  <div class="movie-detail">
    <div v-if="loading" class="loading">Loading movie details...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else-if="movie" class="movie-content">
      <div class="movie-header">
        <button @click="goBack" class="back-button">← Back</button>
        <h1>{{ movie.title }}</h1>
      </div>

      <div class="movie-main">
        <div class="movie-poster">
          <img
            v-if="movie.meta && movie.meta.poster_path"
            :src="`http://localhost:8001/images/w500${movie.meta.poster_path}`"
            :alt="movie.title"
          />
          <div v-else class="no-poster">{{ movie.title }}</div>
        </div>

        <div class="movie-info">
          <div class="movie-meta">
            <p v-if="movie.meta && movie.meta.release_date">
              <strong>Release Date:</strong> {{ formatDate(movie.meta.release_date) }}
            </p>
            <p v-if="movie.meta && movie.meta.Runtime">
              <strong>Runtime:</strong> {{ movie.meta.Runtime }} minutes
            </p>
            <p v-if="movie.meta && movie.meta.Genres && movie.meta.Genres.length">
              <strong>Genres:</strong> {{ movie.meta.Genres.map(g => g.Name).join(', ') }}
            </p>
            <p v-if="movie.meta && movie.meta.vote_average">
              <strong>Rating:</strong> {{ movie.meta.vote_average }}/10 ({{ movie.meta.vote_count }} votes)
            </p>
            <p v-if="movie.meta && movie.meta.spoken_languages && movie.meta.spoken_languages.length">
              <strong>Languages:</strong> {{ movie.meta.spoken_languages.map(l => l.name).join(', ') }}
            </p>
            <p v-if="movie.meta && movie.meta.production_countries && movie.meta.production_countries.length">
              <strong>Countries:</strong> {{ movie.meta.production_countries.map(c => c.Name).join(', ') }}
            </p>
            <p v-if="movie.meta && movie.meta.budget">
              <strong>Budget:</strong> ${{ movie.meta.budget.toLocaleString() }}
            </p>
            <p v-if="movie.meta && movie.meta.revenue">
              <strong>Revenue:</strong> ${{ movie.meta.revenue.toLocaleString() }}
            </p>
            <p v-if="movie.rating !== undefined">
              <strong>User Rating:</strong> {{ movie.rating }}/10
            </p>
            <p>
              <strong>Watchlist:</strong> {{ movie.watchlist ? 'Yes' : 'No' }}
            </p>
          </div>

          <div v-if="movie.meta && movie.meta.Tagline" class="movie-tagline">
            <h3>Tagline</h3>
            <p><em>"{{ movie.meta.Tagline }}"</em></p>
          </div>

          <div v-if="movie.meta && movie.meta.Overview" class="movie-overview">
            <h3>Overview</h3>
            <p>{{ movie.meta.Overview }}</p>
          </div>

          <div v-if="movie.meta && movie.meta.Credits && movie.meta.Credits.Cast && movie.meta.Credits.Cast.length" class="movie-cast">
            <h3>Cast</h3>
            <div class="cast-list">
              <div v-for="actor in movie.meta.Credits.Cast.slice(0, 10)" :key="actor.ID" class="cast-member">
                <strong @click="searchPersonMovies(actor.ID)" class="person-link">{{ actor.Name }}</strong> as {{ actor.Character }}
              </div>
            </div>
          </div>

          <div v-if="movie.meta && movie.meta.Credits && movie.meta.Credits.Crew && movie.meta.Credits.Crew.length" class="movie-crew">
            <h3>Crew</h3>
            <div class="crew-list">
              <div v-for="crew in movie.meta.Credits.Crew.slice(0, 10)" :key="crew.ID" class="crew-member">
                <strong @click="searchPersonMovies(crew.ID)" class="person-link">{{ crew.Name }}</strong> - {{ crew.Job }}
              </div>
            </div>
          </div>

          <div v-if="movie.meta && movie.meta.production_companies && movie.meta.production_companies.length" class="movie-production">
            <h3>Production Companies</h3>
            <div class="production-list">
              <div v-for="company in movie.meta.production_companies" :key="company.id" class="production-company">
                {{ company.Name }}
              </div>
            </div>
          </div>

          <div v-if="movie.File" class="movie-file">
            <h3>File Information</h3>
            <p><strong>Filename:</strong> {{ movie.File.FileName }}</p>
            <p><strong>Path:</strong> {{ movie.File.FullPath }}</p>
            <p><strong>Size:</strong> {{ formatFileSize(movie.File.Size) }}</p>
          </div>

          <div class="movie-actions">
            <button @click="playMovie" class="play-button">Play Movie</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { moviesService } from '../services/movies.js';

const router = useRouter();

const route = useRoute();
const movie = ref(null);
const loading = ref(true);
const error = ref(null);

const fetchMovie = async () => {
  try {
    loading.value = true;
    error.value = null;
    const response = await moviesService.getMovie(route.params.id);
    movie.value = response;
  } catch (err) {
    error.value = 'Failed to load movie details. Please try again.';
    console.error('Error fetching movie:', err);
  } finally {
    loading.value = false;
  }
};

const playMovie = async () => {
  try {
    await moviesService.playMovie(movie.value.id);
    alert('Movie playback started!');
  } catch (err) {
    console.error('Error playing movie:', err);
    alert('Failed to start movie playback.');
  }
};

const formatDate = (dateString) => {
  if (!dateString) return '';
  return new Date(dateString).toLocaleDateString();
};

const formatFileSize = (bytes) => {
  if (!bytes) return 'Unknown';
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  if (bytes === 0) return '0 Bytes';
  const i = parseInt(Math.floor(Math.log(bytes) / Math.log(1024)));
  return Math.round(bytes / Math.pow(1024, i) * 100) / 100 + ' ' + sizes[i];
};

const searchPersonMovies = (personId) => {
  router.push(`/movies?person=${personId}`);
};

const goBack = () => {
  router.go(-1);
};

onMounted(() => {
  fetchMovie();
});
</script>

<style scoped>
.movie-detail {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.loading, .error {
  text-align: center;
  padding: 40px;
  font-size: 18px;
}

.error {
  color: red;
}

.movie-header {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 30px;
}

.back-button {
  padding: 8px 16px;
  background-color: #6c757d;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
}

.back-button:hover {
  background-color: #5a6268;
}

.movie-main {
  display: flex;
  gap: 30px;
}

.movie-poster {
  flex-shrink: 0;
  width: 300px;
}

.movie-poster img {
  width: 100%;
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
}

.no-poster {
  width: 100%;
  height: 450px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  font-size: 18px;
  font-weight: bold;
  text-align: center;
  padding: 20px;
  box-sizing: border-box;
  border-radius: 8px;
  word-wrap: break-word;
  line-height: 1.2;
}

.movie-info {
  flex: 1;
}

.movie-meta {
  margin-bottom: 20px;
}

.movie-meta p {
  margin: 8px 0;
  font-size: 16px;
}

.movie-overview {
  margin-bottom: 30px;
}

.movie-overview h3 {
  margin-bottom: 10px;
  font-size: 20px;
}

.movie-overview p {
  line-height: 1.6;
  font-size: 16px;
}

.movie-tagline {
  margin-bottom: 30px;
}

.movie-tagline h3 {
  margin-bottom: 10px;
  font-size: 18px;
}

.movie-tagline p {
  font-size: 16px;
  font-style: italic;
  color: #666;
}

.movie-cast, .movie-crew, .movie-production, .movie-file {
  margin-bottom: 30px;
}

.movie-cast h3, .movie-crew h3, .movie-production h3, .movie-file h3 {
  margin-bottom: 15px;
  font-size: 18px;
}

.cast-list, .crew-list, .production-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 10px;
}

.cast-member, .crew-member, .production-company {
  padding: 8px;
  background-color: #f8f9fa;
  border-radius: 4px;
  font-size: 14px;
}

.cast-member strong, .crew-member strong {
  color: #007bff;
  cursor: pointer;
}

.cast-member strong:hover, .crew-member strong:hover {
  text-decoration: underline;
}

.movie-file p {
  margin: 5px 0;
  font-size: 14px;
}

.play-button {
  padding: 12px 24px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 18px;
  font-weight: bold;
}

.play-button:hover {
  background-color: #0056b3;
}

@media (max-width: 768px) {
  .movie-main {
    flex-direction: column;
  }

  .movie-poster {
    width: 100%;
    max-width: 300px;
    align-self: center;
  }
}
</style>