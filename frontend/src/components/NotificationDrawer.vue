<script setup lang="ts">
import { computed } from 'vue';
import { useNotificationStore } from '../stores/notification';
import Drawer from 'primevue/drawer';
import Button from 'primevue/button';

const visible = defineModel<boolean>('visible', { default: false });
const store = useNotificationStore();

const hasUnread = computed(() => store.unreadCount > 0);

const typeIcon = (type: string) => {
  switch (type) {
    case 'badge_earned': return 'pi-trophy';
    case 'level_up': return 'pi-arrow-up';
    case 'team_invite': return 'pi-users';
    case 'task_assigned': return 'pi-check-square';
    default: return 'pi-bell';
  }
};

const formatTime = (iso: string) => {
  const date = new Date(iso);
  const diffMs = Date.now() - date.getTime();
  const minutes = Math.floor(diffMs / 60000);
  if (minutes < 1) return 'just now';
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  return date.toLocaleDateString();
};

const onItemClick = async (id: string, readAt?: string) => {
  if (!readAt) {
    try { await store.markRead(id); } catch { /* toast via interceptor */ }
  }
};

const onMarkAllRead = async () => {
  try { await store.markAllRead(); } catch { /* toast via interceptor */ }
};
</script>

<template>
  <Drawer v-model:visible="visible" header="Notifications" position="right" class="w-full md:w-96">
    <div class="flex flex-col h-full">
      <div class="flex justify-end mb-3">
        <Button
          label="Mark all as read"
          icon="pi pi-check-circle"
          text
          size="small"
          :disabled="!hasUnread"
          @click="onMarkAllRead"
        />
      </div>

      <div v-if="store.loading" class="py-8 text-center text-gray-500">
        <i class="pi pi-spin pi-spinner mr-2"></i>Loading...
      </div>

      <div v-else-if="store.items.length === 0" class="flex flex-col items-center py-12 text-gray-400">
        <i class="pi pi-inbox text-4xl mb-3"></i>
        <span class="italic">No notifications yet.</span>
      </div>

      <div v-else class="space-y-2 overflow-y-auto">
        <div
          v-for="n in store.items"
          :key="n.id"
          class="p-3 rounded-lg border cursor-pointer transition-colors"
          :class="n.read_at
            ? 'bg-white dark:bg-gray-900 border-gray-100 dark:border-gray-800 opacity-70'
            : 'bg-primary-50 dark:bg-primary-900/20 border-primary-100 dark:border-primary-800'"
          @click="onItemClick(n.id, n.read_at)"
        >
          <div class="flex items-start gap-3">
            <i :class="['pi', typeIcon(n.type)]" class="mt-1 text-primary-500"></i>
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between gap-2">
                <span class="font-medium text-sm truncate">{{ n.title }}</span>
                <span v-if="!n.read_at" class="w-2 h-2 rounded-full bg-primary-500 shrink-0"></span>
              </div>
              <p v-if="n.message" class="text-xs text-gray-500 mt-1">{{ n.message }}</p>
              <span class="text-[10px] text-gray-400">{{ formatTime(n.created_at) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Drawer>
</template>
