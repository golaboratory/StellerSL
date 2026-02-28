import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useToastStore = defineStore('toast', () => {
  const message = ref<{severity: string, summary: string, detail: string} | null>(null);

  function showToast(severity: string, summary: string, detail: string) {
    message.value = { severity, summary, detail };
  }

  function clearToast() {
    message.value = null;
  }

  return { message, showToast, clearToast };
});
