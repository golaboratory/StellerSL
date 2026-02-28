import axios from 'axios';
import router from '../router';
import { useAuthStore } from '../stores/auth';
import { useToastStore } from '../stores/toast';

const axiosInstance = axios.create({
  baseURL: '/api'
});

axiosInstance.interceptors.request.use((config) => {
  const auth = useAuthStore();
  if (auth.token) {
    config.headers.Authorization = `Bearer ${auth.token}`;
  }
  return config;
});

axiosInstance.interceptors.response.use(
  (response) => response,
  (error) => {
    const toast = useToastStore();
    
    if (error.response) {
      const status = error.response.status;
      const data = error.response.data;

      if (status === 401) {
        const auth = useAuthStore();
        auth.clearAuth();
        router.push('/login');
      } else if (status === 400 && data.errors) {
        // Huma validation errors
        const detail = data.errors.map((e: any) => `${e.location}: ${e.message}`).join('\n');
        toast.showToast('error', 'Validation Error', detail || data.detail || 'Invalid request');
      } else {
        toast.showToast('error', `API Error (${status})`, data.detail || data.message || 'Something went wrong');
      }
    } else {
      toast.showToast('error', 'Network Error', 'Could not connect to the server');
    }
    
    return Promise.reject(error);
  }
);

export default axiosInstance;
