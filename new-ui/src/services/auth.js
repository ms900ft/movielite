import { loginApi } from './api.js';

export const authService = {
  async login(credentials) {
    try {
      const response = await loginApi.post('/login', credentials);
      const { token, user_name, is_admin } = response.data;

      // Store token in localStorage
      localStorage.setItem('authToken', token);
      localStorage.setItem('userName', user_name);
      localStorage.setItem('isAdmin', is_admin);

      return { token, user_name, is_admin };
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  },

  logout() {
    localStorage.removeItem('authToken');
    localStorage.removeItem('userName');
    localStorage.removeItem('isAdmin');
  },

  isAuthenticated() {
    return !!localStorage.getItem('authToken');
  },

  getUser() {
    return {
      name: localStorage.getItem('userName'),
      isAdmin: localStorage.getItem('isAdmin') === 'true'
    };
  }
};