<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue';
import { useAuthStore } from '../stores/auth';
import { DefaultApi, Configuration } from '../api';
import Gantt from 'frappe-gantt';
import Toolbar from 'primevue/toolbar';
import Button from 'primevue/button';

const auth = useAuthStore();
const api = new DefaultApi(new Configuration({ basePath: '/api', accessToken: auth.token || undefined }));

const tasks = ref<any[]>([]);
const ganttChart = ref<any>(null);

onMounted(async () => {
    try {
        const response = await api.listTasks();
        tasks.value = response.data.items || [];
        
        if (tasks.value.length > 0) {
            const ganttTasks = tasks.value.map(t => ({
                id: t.id,
                name: t.title,
                start: new Date().toISOString(), // Mock start
                end: new Date(Date.now() + 86400000 * 2).toISOString(), // Mock end
                progress: t.status === 'done' ? 100 : (t.status === 'doing' ? 50 : 0)
            }));

            await nextTick();
            ganttChart.value = new Gantt("#gantt-target", ganttTasks, {
                header_height: 50,
                column_width: 30,
                step: 24,
                view_modes: ['Day', 'Week', 'Month'],
                bar_height: 20,
                bar_corner_radius: 3,
                arrow_curve: 5,
                padding: 18,
                view_mode: 'Day',
                date_format: 'YYYY-MM-DD',
                custom_popup_html: null,
                on_date_change: async (task: any, start: Date, end: Date) => {
                    console.log("Task date changed", task.id, start, end);
                    // In a real app, we would call an API like api.updateTask(task.id, { due_date: end })
                    // For now, we'll just log it.
                }
            });
        }
    } catch (err) {
        console.error(err);
    }
});
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950">
    <Toolbar class="p-4 shadow-sm">
      <template #start>
        <span class="text-2xl font-bold text-primary px-4">Gantt Chart</span>
      </template>
      <template #end>
        <router-link to="/dashboard">
          <Button icon="pi pi-home" label="Dashboard" text />
        </router-link>
      </template>
    </Toolbar>

    <div class="p-6">
        <div id="gantt-target" class="bg-white dark:bg-gray-900 border rounded-xl overflow-hidden p-4"></div>
        <div v-if="tasks.length === 0" class="text-center py-20 text-gray-500">
            No tasks to display in Gantt.
        </div>
    </div>
  </div>
</template>

<style>
/* Frappe Gantt styles might need to be imported or manually adjusted */
.gantt-container {
    overflow: auto;
}
</style>
