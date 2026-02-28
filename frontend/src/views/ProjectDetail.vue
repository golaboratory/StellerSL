<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { DefaultApi, TaskItem, ProjectOutputBody } from '../api';
import axiosInstance from '../api/axios';
import Card from 'primevue/card';
import Button from 'primevue/button';
import Toolbar from 'primevue/toolbar';
import Tag from 'primevue/tag';

const route = useRoute();
const router = useRouter();
const auth = useAuthStore();
const api = new DefaultApi(undefined, '/api', axiosInstance);

const projectID = route.params.id as string;
const project = ref<any>(null);
const tasks = ref<TaskItem[]>([]);
const loading = ref(true);

const fetchProjectDetails = async () => {
    loading.value = true;
    try {
        // Since we don't have a specific GetProject API that returns tasks, 
        // we'll fetch the project list and find ours, then fetch all tasks and filter.
        // In a real app, the API should handle this filtering.
        const [projRes, tasksRes] = await Promise.all([
            api.listProjects(),
            api.listTasks()
        ]);
        
        project.value = projRes.data.items?.find((p: any) => p.id === projectID);
        tasks.value = (tasksRes.data.items || []).filter((t: any) => t.project_id === projectID);
    } catch (err) {
        console.error('Failed to fetch project details', err);
    } finally {
        loading.value = false;
    }
};

onMounted(fetchProjectDetails);

const getStatusSeverity = (status: string) => {
    switch (status) {
        case 'todo': return 'secondary';
        case 'doing': return 'warn';
        case 'done': return 'success';
        default: return 'info';
    }
};
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950 p-6">
    <Toolbar class="mb-8 p-4 rounded-xl shadow-sm">
      <template #start>
        <div class="flex items-center gap-4 px-4">
            <Button icon="pi pi-arrow-left" text @click="router.push('/projects')" />
            <span class="text-2xl font-bold text-primary">{{ project?.name || 'Project Details' }}</span>
        </div>
      </template>
      <template #end>
        <router-link to="/dashboard">
          <Button icon="pi pi-home" label="Dashboard" text />
        </router-link>
      </template>
    </Toolbar>

    <div v-if="loading" class="space-y-6 animate-pulse">
        <div class="h-32 bg-gray-200 dark:bg-gray-800 rounded-xl"></div>
        <div class="h-64 bg-gray-200 dark:bg-gray-800 rounded-xl"></div>
    </div>

    <div v-else class="grid grid-cols-1 gap-6">
      <Card class="shadow-sm border-t-4 border-primary-500">
        <template #title>Description</template>
        <template #content>
          <p class="text-gray-600 dark:text-gray-400">{{ project?.description || 'No description provided.' }}</p>
        </template>
      </Card>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card class="shadow-sm">
            <template #title>Project Members</template>
            <template #content>
                <div class="space-y-4">
                    <p class="text-sm text-gray-500 italic">User assignment management UI coming soon...</p>
                    <Button label="Assign New Member" icon="pi pi-user-plus" severity="secondary" outlined class="w-full" />
                </div>
            </template>
        </Card>

        <Card class="shadow-sm">
            <template #title>Project Stats</template>
            <template #content>
                <div class="flex justify-around text-center">
                    <div>
                        <div class="text-2xl font-bold">{{ tasks.length }}</div>
                        <div class="text-xs text-gray-500 uppercase">Tasks</div>
                    </div>
                    <div>
                        <div class="text-2xl font-bold text-green-500">{{ tasks.filter(t => t.status === 'done').length }}</div>
                        <div class="text-xs text-gray-500 uppercase">Done</div>
                    </div>
                </div>
            </template>
        </Card>
      </div>

      <Card class="shadow-sm">
        <template #title>Project Tasks ({{ tasks.length }})</template>
        <template #content>
            <div class="space-y-4">
                <div v-for="task in tasks" :key="task.id" class="p-4 border rounded-lg flex items-center justify-between bg-white dark:bg-gray-900 border-gray-100 dark:border-gray-800">
                    <span class="font-medium">{{ task.title }}</span>
                    <Tag :value="task.status.toUpperCase()" :severity="getStatusSeverity(task.status)" />
                </div>
                <div v-if="tasks.length === 0" class="text-center py-8 text-gray-500 italic">
                    No tasks assigned to this project yet.
                </div>
            </div>
        </template>
      </Card>
    </div>
  </div>
</template>
