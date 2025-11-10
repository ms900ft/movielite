import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8001/api', // All API endpoints are under /api
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add request interceptor to include auth token if available
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('authToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Create separate instance for login (not under /api)
const loginApi = axios.create({
  baseURL: 'http://localhost:8001',
  headers: {
    'Content-Type': 'application/json',
  },
});

export { loginApi };

export default api;