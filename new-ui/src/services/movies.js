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
  }
};