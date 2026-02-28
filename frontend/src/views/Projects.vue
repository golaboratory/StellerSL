<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useAuthStore } from '../stores/auth';
import { DefaultApi } from '../api';
import axiosInstance from '../api/axios';
import Card from 'primevue/card';
import Button from 'primevue/button';
import Toolbar from 'primevue/toolbar';
import { useRouter } from 'vue-router';

const auth = useAuthStore();
const router = useRouter();
const api = new DefaultApi(undefined, '/api', axiosInstance);

const projects = ref<any[]>([]);
const loading = ref(true);

const fetchProjects = async () => {
    loading.value = true;
    try {
        const response = await api.listProjects();
        projects.value = response.data.items || [];
    } catch (err) {
        console.error('Failed to fetch projects', err);
    } finally {
        loading.value = false;
    }
};

onMounted(fetchProjects);
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950 p-6">
    <Toolbar class="mb-8 p-4 rounded-xl shadow-sm">
      <template #start>
        <span class="text-2xl font-bold text-primary px-4">Projects</span>
      </template>
      <template #end>
        <router-link to="/dashboard">
          <Button icon="pi pi-home" label="Dashboard" text />
        </router-link>
        <router-link to="/tasks">
          <Button icon="pi pi-check-square" label="Tasks" text />
        </router-link>
      </template>
    </Toolbar>

    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-3 gap-6 animate-pulse">
        <div v-for="i in 3" :key="i" class="h-32 bg-gray-200 dark:bg-gray-800 rounded-xl"></div>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <Card v-for="project in projects" :key="project.id" class="shadow-sm cursor-pointer hover:shadow-md transition-shadow" @click="router.push(`/projects/${project.id}`)">
        <template #title>{{ project.name }}</template>
        <template #content>
          <p class="text-gray-500 line-clamp-2">{{ project.description || 'No description' }}</p>
        </template>
      </Card>
      <div v-if="projects.length === 0" class="col-span-full text-center text-gray-500 py-8">
        No projects found.
      </div>
    </div>
  </div>
</template>
