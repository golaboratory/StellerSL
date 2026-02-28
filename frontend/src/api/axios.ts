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
      if (error.response.status === 401) {
        const auth = useAuthStore();
        auth.clearAuth();
        router.push('/login');
      } else {
        toast.showToast('error', 'API Error', error.response.data.message || 'Something went wrong');
      }
    } else {
      toast.showToast('error', 'Network Error', 'Could not connect to the server');
    }
    
    return Promise.reject(error);
  }
);

export default axiosInstance;
