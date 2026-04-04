import api from './api.js';

export const moviesService = {
  async getMovies(params = {}) {
    try {
      const response = await api.get('/movie', { params });
      return response.data;
    } catch (error) {
      console.error('Error fetching movies:', error);
      throw error;
    }
  },

  async getMovie(id) {
    try {
      const response = await api.get(`/movie/${id}`);
      return response.data;
    } catch (error) {
      console.error('Error fetching movie:', error);
      throw error;
    }
  },

  async updateMovie(id, movieData) {
    try {
      const response = await api.put(`/movie/${id}`, movieData);
      return response.data;
    } catch (error) {
      console.error('Error updating movie:', error);
      throw error;
    }
  },

  async deleteMovie(id) {
    try {
      const response = await api.delete(`/movie/${id}`);
      return response.data;
    } catch (error) {
      console.error('Error deleting movie:', error);
      throw error;
    }
  },

  async playMovie(id) {
    try {
      const response = await api.put(`/movie/${id}/play`);
      return response.data;
    } catch (error) {
      console.error('Error playing movie:', error);
      throw error;
    }
  },

  async showMovie(id) {
    try {
      const response = await api.put(`/movie/${id}/play`, null, { params: { showdir: 1 } });
      return response.data;
    } catch (error) {
      console.error('Error showing movie:', error);
      throw error;
    }
  },

  async getGenres() {
    try {
      const response = await api.get('/genre');
      return response.data;
    } catch (error) {
      console.error('Error fetching genres:', error);
      throw error;
    }
  },

  async getCountries() {
    try {
      const response = await api.get('/country');
      return response.data;
    } catch (error) {
      console.error('Error fetching countries:', error);
      throw error;
    }
  },

  async getPerson(id) {
    try {
      const response = await api.get(`/person/${id}`);
      return response.data;
    } catch (error) {
      console.error('Error fetching person:', error);
      throw error;
    }
  },

  async getTargets() {
    try {
      const response = await api.get('/targets');
      return response.data;
    } catch (error) {
      console.error('Error fetching targets:', error);
      throw error;
    }
  },

  async moveFile(id, dir) {
    try {
      const response = await api.put(`/file/${id}/move/${encodeURIComponent(dir)}`);
      return response.data;
    } catch (error) {
      console.error('Error moving file:', error);
      throw error;
    }
  },

  async rescan() {
    try {
      const response = await api.post('/file/rescan');
      return response.data;
    } catch (error) {
      console.error('Error rescanning:', error);
      throw error;
    }
  },

  async rescanMovie(id, metaId) {
    try {
      const response = await api.put(`/movie/${id}/addMeta/${metaId}`);
      return response.data;
    } catch (error) {
      console.error('Error rescanning movie:', error);
      throw error;
    }
  },

  async getUsers() {
    try {
      const response = await api.get('/user');
      return response.data;
    } catch (error) {
      console.error('Error fetching users:', error);
      throw error;
    }
  },

  async createUser(userData) {
    try {
      const response = await api.post('/user', userData);
      return response.data;
    } catch (error) {
      console.error('Error creating user:', error);
      throw error;
    }
  },

  async updateUser(id, userData) {
    try {
      const response = await api.put(`/user/${id}`, userData);
      return response.data;
    } catch (error) {
      console.error('Error updating user:', error);
      throw error;
    }
  },

  async deleteUser(id) {
    try {
      const response = await api.delete(`/user/${id}`);
      return response.data;
    } catch (error) {
      console.error('Error deleting user:', error);
      throw error;
    }
  }
};