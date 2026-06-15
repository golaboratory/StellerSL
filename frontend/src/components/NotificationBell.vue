<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { useNotificationStore } from '../stores/notification';
import NotificationDrawer from './NotificationDrawer.vue';
import Button from 'primevue/button';
import OverlayBadge from 'primevue/overlaybadge';

const store = useNotificationStore();
const drawerVisible = ref(false);

onMounted(() => store.startPolling());
onUnmounted(() => store.stopPolling());

const openDrawer = async () => {
  drawerVisible.value = true;
  try { await store.fetchList(); } catch { /* toast via interceptor */ }
};
</script>

<template>
  <div>
    <OverlayBadge v-if="store.unreadCount > 0" :value="store.unreadCount > 99 ? '99+' : String(store.unreadCount)" severity="danger">
      <Button icon="pi pi-bell" text aria-label="Notifications" @click="openDrawer" />
    </OverlayBadge>
    <Button v-else icon="pi pi-bell" text aria-label="Notifications" @click="openDrawer" />

    <NotificationDrawer v-model:visible="drawerVisible" />
  </div>
</template>
