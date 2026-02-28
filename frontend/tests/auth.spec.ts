import { describe, it, expect, beforeEach } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useAuthStore } from '../src/stores/auth';

describe('Auth Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    localStorage.clear();
  });

  it('initializes with null token and user', () => {
    const auth = useAuthStore();
    expect(auth.token).toBeNull();
    expect(auth.user).toBeNull();
    expect(auth.isAuthenticated).toBe(false);
  });

  it('sets authentication state correctly', () => {
    const auth = useAuthStore();
    const mockUser = { id: '1', name: 'Test', email: 'test@example.com' };
    
    auth.setAuth('mock-token', mockUser);
    
    expect(auth.token).toBe('mock-token');
    expect(auth.user).toEqual(mockUser);
    expect(auth.isAuthenticated).toBe(true);
    expect(localStorage.getItem('token')).toBe('mock-token');
  });

  it('clears authentication state correctly', () => {
    const auth = useAuthStore();
    auth.setAuth('mock-token', { id: '1', name: 'Test', email: 'test@example.com' });
    
    auth.clearAuth();
    
    expect(auth.token).toBeNull();
    expect(auth.user).toBeNull();
    expect(auth.isAuthenticated).toBe(false);
    expect(localStorage.getItem('token')).toBeNull();
  });
});
