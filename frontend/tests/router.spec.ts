import { describe, it, expect, beforeEach, vi } from 'vitest';
import { createPinia, setActivePinia } from 'pinia';
import router from '../src/router';
import { useAuthStore } from '../src/stores/auth';

describe('Router Navigation Guards', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    localStorage.clear();
  });

  it('redirects unauthenticated users to login', async () => {
    const auth = useAuthStore();
    auth.clearAuth();
    
    await router.push('/dashboard');
    expect(router.currentRoute.value.path).toBe('/login');
  });

  it('allows authenticated users to access dashboard', async () => {
    const auth = useAuthStore();
    auth.setAuth('mock-token', { id: '1', name: 'Test', email: 'test@example.com' });
    
    await router.push('/dashboard');
    expect(router.currentRoute.value.path).toBe('/dashboard');
  });

  it('redirects authenticated users away from login to dashboard', async () => {
    const auth = useAuthStore();
    auth.setAuth('mock-token', { id: '1', name: 'Test', email: 'test@example.com' });
    
    await router.push('/login');
    expect(router.currentRoute.value.path).toBe('/dashboard');
  });
});
