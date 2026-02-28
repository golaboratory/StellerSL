<script setup lang="ts">
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';
import { watch } from 'vue';
import { useToastStore } from './stores/toast';

const toast = useToast();
const toastStore = useToastStore();

watch(() => toastStore.message, (msg) => {
  if (msg) {
    toast.add({
      severity: msg.severity as any,
      summary: msg.summary,
      detail: msg.detail,
      life: 3000
    });
    toastStore.clearToast();
  }
});
</script>

<template>
  <Toast />
  <router-view />
</template>

<style>
/* Global styles if needed */
</style>
