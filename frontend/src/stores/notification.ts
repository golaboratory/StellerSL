import { defineStore } from 'pinia';
import { ref } from 'vue';
import { DefaultApi } from '../api';
import type { NotificationItem } from '../api';
import axiosInstance from '../api/axios';

const api = new DefaultApi(undefined, '/api', axiosInstance);
const POLL_INTERVAL_MS = 30_000;

export const useNotificationStore = defineStore('notification', () => {
  const items = ref<NotificationItem[]>([]);
  const unreadCount = ref(0);
  const loading = ref(false);

  let pollTimer: ReturnType<typeof setInterval> | null = null;

  async function fetchUnreadCount() {
    try {
      const res = await api.countUnreadNotifications();
      unreadCount.value = res.data.count;
    } catch {
      // Polling errors are silent: the axios interceptor already handles 401
    }
  }

  async function fetchList() {
    loading.value = true;
    try {
      const res = await api.listNotifications({ limit: 50 });
      items.value = res.data.items || [];
    } finally {
      loading.value = false;
    }
  }

  async function markRead(id: string) {
    await api.markNotificationRead({ id });
    const item = items.value.find((n) => n.id === id);
    if (item && !item.read_at) {
      item.read_at = new Date().toISOString();
      unreadCount.value = Math.max(0, unreadCount.value - 1);
    }
  }

  async function markAllRead() {
    await api.markAllNotificationsRead();
    const now = new Date().toISOString();
    items.value.forEach((n) => {
      if (!n.read_at) n.read_at = now;
    });
    unreadCount.value = 0;
  }

  function startPolling() {
    if (pollTimer) return;
    fetchUnreadCount();
    pollTimer = setInterval(fetchUnreadCount, POLL_INTERVAL_MS);
  }

  function stopPolling() {
    if (pollTimer) {
      clearInterval(pollTimer);
      pollTimer = null;
    }
  }

  return {
    items,
    unreadCount,
    loading,
    fetchUnreadCount,
    fetchList,
    markRead,
    markAllRead,
    startPolling,
    stopPolling,
  };
});
