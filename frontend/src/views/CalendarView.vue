<script setup lang="ts">
import { ref, onMounted } from 'vue';
import axiosInstance from '../api/axios';
import { DefaultApi } from '../api';
import { Calendar } from 'v-calendar';
import 'v-calendar/dist/style.css';
import Toolbar from 'primevue/toolbar';
import Button from 'primevue/button';

const api = new DefaultApi(undefined, '/api', axiosInstance);

const tasks = ref<any[]>([]);
const attributes = ref<any[]>([]);
const selectedDate = ref(new Date());
const selectedTasks = ref<any[]>([]);

onMounted(async () => {
    try {
        const response = await api.listTasks();
        tasks.value = response.data.items || [];
        
        attributes.value = tasks.value.map(t => ({
            key: t.id,
            dot: t.status === 'done' ? 'green' : (t.status === 'doing' ? 'orange' : 'blue'),
            customData: t,
            dates: t.due_date ? new Date(t.due_date) : new Date(), 
            popover: {
                label: t.title,
                visibility: 'hover'
            }
        }));
        
        updateSelectedTasks();
    } catch (err) {
        console.error(err);
    }
});

const updateSelectedTasks = () => {
    selectedTasks.value = tasks.value.filter(t => {
        const d = t.due_date ? new Date(t.due_date) : new Date();
        return d.toDateString() === selectedDate.value.toDateString();
    });
};

const onDayClick = (day: any) => {
    selectedDate.value = day.date;
    updateSelectedTasks();
};
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950">
    <Toolbar class="p-4 shadow-sm">
      <template #start>
        <span class="text-2xl font-bold text-primary px-4">Calendar</span>
      </template>
      <template #end>
        <router-link to="/dashboard">
          <Button icon="pi pi-home" label="Dashboard" text />
        </router-link>
      </template>
    </Toolbar>

    <div class="p-6 flex flex-col items-center gap-6">
        <Calendar 
            :attributes="attributes"
            expanded
            is-expanded
            class="max-w-4xl w-full border-none shadow-xl rounded-2xl p-4"
            :is-dark="true"
            @dayclick="onDayClick"
        />

        <div class="max-w-4xl w-full bg-white dark:bg-gray-900 rounded-2xl shadow-lg p-6">
            <h3 class="text-xl font-bold mb-4">Tasks for {{ selectedDate.toLocaleDateString() }}</h3>
            <div class="space-y-3">
                <div v-for="task in selectedTasks" :key="task.id" class="p-3 border rounded-lg flex justify-between items-center dark:border-gray-800">
                    <span>{{ task.title }}</span>
                    <span class="text-xs uppercase px-2 py-1 bg-gray-100 dark:bg-gray-800 rounded">{{ task.status }}</span>
                </div>
                <div v-if="selectedTasks.length === 0" class="text-center py-4 text-gray-500 italic">
                    No tasks for this day.
                </div>
            </div>
        </div>
    </div>
  </div>
</template>
