<template>
  <div class="movie-detail-overlay-content">
    <div v-if="loading" class="loading">Loading movie details...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else-if="movie" class="movie-content">
      <div class="movie-main">
        <div class="movie-poster">
          <img
            v-if="movie.meta && movie.meta.poster_path"
            :src="`/images/w500${movie.meta.poster_path}`"
            :alt="movie.title"
            class="poster-clickable"
            @click="openModal(`/images/original${movie.meta.poster_path}`)"
          />
          <div v-else class="no-poster">{{ movie.title }}</div>
        </div>

        <div class="movie-info">
          <template v-if="editingTitle">
            <div class="title-edit-row">
              <input v-model="editTitle" @keyup.enter="saveTitle" @keyup.escape="cancelEditTitle" class="title-input" />
              <button @click="saveTitle" class="save-button">Save</button>
              <button @click="cancelEditTitle" class="cancel-button">Cancel</button>
            </div>
          </template>
          <template v-else>
            <h1>{{ movie.title }} <button @click="startEditTitle" class="edit-title-button" title="Edit title">✏️</button></h1>
          </template>

          <div class="movie-meta">
            <div class="meta-row">
              <span class="meta-label">TMDB ID:</span>
              <template v-if="editingTmdbId">
                <input v-model="editTmdbId" @keyup.enter="saveTmdbId" @keyup.escape="cancelEditTmdbId" class="tmdb-input" />
                <button @click="saveTmdbId" class="save-button">Save</button>
                <button @click="cancelEditTmdbId" class="cancel-button">Cancel</button>
              </template>
              <template v-else>
                <span class="meta-value">{{ movie.TMDBMovieID || 'N/A' }}</span>
                <button @click="startEditTmdbId" class="edit-tmdb-button" title="Edit TMDB ID">✏️</button>
              </template>
            </div>
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

          <div class="movie-credits">
            <div v-if="movie.meta && movie.meta.Credits && movie.meta.Credits.Cast && movie.meta.Credits.Cast.length" class="movie-cast">
              <h3>Cast</h3>
              <div class="cast-list">
                <div v-for="actor in movie.meta.Credits.Cast" :key="actor.ID" class="cast-member">
                  <img v-if="actor.profile_path" :src="`/images/w185${actor.profile_path}`" :alt="actor.Name" class="person-image" @click="openModal(`/images/w500${actor.profile_path}`)" />
                  <div v-else class="person-image-placeholder"></div>
                  <strong @click="searchPersonMovies(actor.ID)" class="person-link">{{ actor.Name }}</strong> as {{ actor.Character }}
                </div>
              </div>
            </div>

            <div v-if="movie.meta && movie.meta.Credits && movie.meta.Credits.Crew && movie.meta.Credits.Crew.length" class="movie-crew">
              <h3>Crew</h3>
              <div class="crew-list">
                <div v-for="crew in movie.meta.Credits.Crew" :key="crew.ID" class="crew-member">
                  <img v-if="crew.profile_path" :src="`/images/w185${crew.profile_path}`" :alt="crew.Name" class="person-image" @click="openModal(`/images/w500${crew.profile_path}`)" />
                  <div v-else class="person-image-placeholder"></div>
                  <strong @click="searchPersonMovies(crew.ID)" class="person-link">{{ crew.Name }}</strong> - {{ crew.Job }}
                </div>
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
        </div>
      </div>
    </div>

    <!-- Image Modal -->
    <div v-if="modalVisible" class="image-modal" @click="closeModal">
      <div class="modal-content" @click.stop>
        <img :src="modalImage" :alt="modalImage" class="modal-image" @click="closeModal" />
      </div>
    </div>

    <!-- Move Movie Dialog -->
    <Dialog v-model:visible="moveDialogVisible" modal header="Move Movie" :style="{ width: '50rem' }">
      <p>Select a directory to move the movie to:</p>
      <div class="target-list">
        <button v-for="target in targets" :key="target.name" @click="moveMovie(target.name)" class="target-button">
          {{ target.name }}
        </button>
      </div>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { moviesService } from '../services/movies.js';
import Dialog from 'primevue/dialog';

const props = defineProps({
  movieId: {
    type: [String, Number],
    required: true
  }
});

const emit = defineEmits(['close', 'searchPerson']);

const movie = ref(null);
const loading = ref(true);
const error = ref(null);
const modalVisible = ref(false);
const modalImage = ref('');
const targets = ref([]);
const moveDialogVisible = ref(false);
const editingTmdbId = ref(false);
const editTmdbId = ref('');
const editingTitle = ref(false);
const editTitle = ref('');

const fetchMovie = async () => {
  try {
    loading.value = true;
    error.value = null;
    const response = await moviesService.getMovie(props.movieId);
    movie.value = response;
  } catch (err) {
    error.value = 'Failed to load movie details. Please try again.';
    console.error('Error fetching movie:', err);
  } finally {
    loading.value = false;
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

const playMovie = async () => {
  try {
    await moviesService.playMovie(movie.value.id);
  } catch (err) {
    console.error('Error playing movie:', err);
  }
};

const moveMovie = async (targetDir) => {
  const fileId = movie.value.File?.ID;
  if (!fileId) {
    alert('File information not available for this movie.');
    return;
  }
  try {
    await moviesService.moveFile(fileId, targetDir);
    alert(`Movie moved to ${targetDir}`);
    emit('close');
    moveDialogVisible.value = false;
  } catch (err) {
    console.error('Error moving movie:', err);
    alert('Failed to move movie.');
  }
};

const showMoveDialog = () => {
  moveDialogVisible.value = true;
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
  emit('searchPerson', personId);
};

const openModal = (imageSrc) => {
  modalImage.value = imageSrc;
  modalVisible.value = true;
};

const closeModal = () => {
  modalVisible.value = false;
};

const startEditTitle = () => {
  editingTitle.value = true;
  editTitle.value = movie.value.title;
};

const saveTitle = async () => {
  if (!editTitle.value.trim()) return;
  try {
    await moviesService.updateMovie(movie.value.id, { title: editTitle.value.trim() });
    movie.value.title = editTitle.value.trim();
    editingTitle.value = false;
  } catch (err) {
    console.error('Error updating title:', err);
    alert('Failed to update title.');
  }
};

const cancelEditTitle = () => {
  editingTitle.value = false;
  editTitle.value = '';
};

const startEditTmdbId = () => {
  editingTmdbId.value = true;
  editTmdbId.value = movie.value.TMDBMovieID || '';
};

const saveTmdbId = async () => {
  const metaId = parseInt(editTmdbId.value.trim(), 10);
  if (!metaId) {
    movie.value.TMDBMovieID = 0;
    movie.value.meta = null;
    editingTmdbId.value = false;
    try {
      await moviesService.updateMovie(movie.value.id, { title: movie.value.title });
    } catch (err) {
      console.error('Error updating movie:', err);
    }
    return;
  }
  try {
    const updatedMovie = await moviesService.rescanMovie(movie.value.id, metaId);
    movie.value = updatedMovie;
    editingTmdbId.value = false;
  } catch (err) {
    console.error('Error updating TMDB ID:', err);
    alert('Failed to update TMDB metadata.');
  }
};

const cancelEditTmdbId = () => {
  editingTmdbId.value = false;
  editTmdbId.value = '';
};

onMounted(() => {
  fetchMovie();
  fetchTargets();
});
</script>

<style scoped>
.movie-detail-overlay-content {
  height: 100%;
  overflow-y: auto;
  padding: 20px;
  background: white;
  border-radius: 12px;
}

.loading, .error {
  text-align: center;
  padding: 40px;
  font-size: 18px;
}

.error {
  color: red;
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

.movie-info h1 {
  margin: 0 0 20px 0;
  font-size: 2rem;
}

.movie-info h1 .edit-title-button {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 18px;
  padding: 4px;
  opacity: 0.6;
  transition: opacity 0.2s;
  filter: grayscale(100%);
}

.movie-info h1 .edit-title-button:hover {
  opacity: 1;
}

.title-edit-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 20px;
}

.title-edit-row .title-input {
  font-size: 1.5rem;
  padding: 6px 12px;
  border: 2px solid #007bff;
  border-radius: 4px;
  flex: 1;
}

.movie-meta {
  margin-bottom: 20px;
}

.movie-meta p {
  margin: 8px 0;
  font-size: 16px;
}

.meta-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 8px 0;
  font-size: 16px;
}

.meta-label {
  font-weight: bold;
  min-width: 100px;
}

.meta-value {
  color: #333;
}

.edit-tmdb-button {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 14px;
  padding: 2px;
  opacity: 0.6;
  transition: opacity 0.2s;
  filter: grayscale(100%);
}

.edit-tmdb-button:hover {
  opacity: 1;
}

.tmdb-input {
  font-size: 14px;
  padding: 4px 8px;
  border: 2px solid #007bff;
  border-radius: 4px;
  width: 120px;
}

.save-button {
  padding: 6px 16px;
  background-color: #28a745;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.save-button:hover {
  background-color: #218838;
}

.cancel-button {
  padding: 6px 16px;
  background-color: #6c757d;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.cancel-button:hover {
  background-color: #5a6268;
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

.movie-credits {
  display: flex;
  gap: 30px;
  margin-bottom: 30px;
}

.movie-cast, .movie-crew {
  flex: 1;
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

.person-image {
  width: 50px;
  height: 75px;
  object-fit: cover;
  margin-right: 10px;
  border-radius: 4px;
  float: left;
  cursor: pointer;
}

.person-image-placeholder {
  width: 50px;
  height: 75px;
  margin-right: 10px;
  border-radius: 4px;
  float: left;
  background-color: #e2e8f0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.person-image-placeholder::after {
  content: '👤';
  font-size: 24px;
  opacity: 0.5;
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

.move-button {
  padding: 12px 24px;
  background-color: #28a745;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 18px;
  font-weight: bold;
  margin-left: 10px;
}

.move-button:hover {
  background-color: #218838;
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

.image-modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.9);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 10000;
}

.modal-content {
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-image {
  max-width: 90vw;
  max-height: 90vh;
  object-fit: contain;
  cursor: pointer;
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
