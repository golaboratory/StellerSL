<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue';
import axiosInstance from '../api/axios';
import { DefaultApi } from '../api';
import { useToastStore } from '../stores/toast';
import Gantt from 'frappe-gantt';
import Toolbar from 'primevue/toolbar';
import Button from 'primevue/button';

const api = new DefaultApi(undefined, '/api', axiosInstance);
const toast = useToastStore();

const tasks = ref<any[]>([]);

const DAY_MS = 86400000;

// Tasks without a due date are shown as a 2-day bar starting today;
// dragging them assigns a real due date via the API.
const toGanttTask = (t: any) => {
    const end = t.due_date ? new Date(t.due_date) : new Date(Date.now() + DAY_MS * 2);
    let start = t.created_at ? new Date(t.created_at) : new Date(end.getTime() - DAY_MS * 2);
    if (start.getTime() >= end.getTime()) {
        start = new Date(end.getTime() - DAY_MS);
    }
    return {
        id: t.id,
        name: t.title,
        start: start.toISOString(),
        end: end.toISOString(),
        progress: t.status === 'done' ? 100 : (t.status === 'doing' ? 50 : 0)
    };
};

const onDateChange = async (ganttTask: any, _start: Date, end: Date) => {
    const task = tasks.value.find(t => t.id === ganttTask.id);
    if (!task) return;
    try {
        // PUT /tasks/{id} overwrites all fields, so resend the current values
        await api.updateTask({
            id: task.id,
            updateTaskRequest: {
                title: task.title,
                description: task.description || '',
                status: task.status,
                priority: task.priority,
                project_id: task.project_id || undefined,
                assigned_to: task.assigned_to || undefined,
                due_date: end.toISOString(),
            }
        });
        task.due_date = end.toISOString();
        toast.showToast('success', 'Updated', `Due date of "${task.title}" set to ${end.toLocaleDateString()}`);
    } catch (err) {
        console.error('Failed to update due date', err);
    }
};

onMounted(async () => {
    try {
        const response = await api.listTasks();
        tasks.value = response.data.items || [];

        if (tasks.value.length > 0) {
            const ganttTasks = tasks.value.map(toGanttTask);

            await nextTick();
            // Note: do not store the instance in a reactive ref — Vue's proxy
            // would wrap the library's internal DOM references
            new Gantt("#gantt-target", ganttTasks, {
                column_width: 30,
                view_modes: ['Day', 'Week', 'Month'],
                bar_height: 20,
                bar_corner_radius: 3,
                arrow_curve: 5,
                padding: 18,
                view_mode: 'Day',
                view_mode_select: true,
                date_format: 'YYYY-MM-DD',
                on_date_change: onDateChange
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
        <p class="text-xs text-gray-500 mb-2 px-1">
            <i class="pi pi-info-circle mr-1"></i>Drag a bar to change its schedule — the due date is saved automatically.
        </p>
        <div id="gantt-target" class="bg-white dark:bg-gray-900 border rounded-xl overflow-hidden p-4"></div>
        <div v-if="tasks.length === 0" class="text-center py-20 text-gray-500">
            No tasks to display in Gantt.
        </div>
    </div>
  </div>
</template>

<style>
/* The package "exports" map blocks deep JS imports of the stylesheet and the
   "style" condition is not applied by postcss-import, so reference the file
   directly through node_modules */
@import "../../node_modules/frappe-gantt/dist/frappe-gantt.css";

.gantt-container {
    overflow: auto;
}
</style>
