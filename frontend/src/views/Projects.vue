<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { DefaultApi } from '../api';
import type { ProjectItem } from '../api';
import axiosInstance from '../api/axios';
import Card from 'primevue/card';
import Toolbar from 'primevue/toolbar';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Textarea from 'primevue/textarea';

const api = new DefaultApi(undefined, '/api', axiosInstance);

const projects = ref<ProjectItem[]>([]);
const loading = ref(true);

// Pagination
const limit = ref(50);
const offset = ref(0);

// Create Project Dialog
const showCreateDialog = ref(false);
const newProject = ref({ name: '', description: '' });
const createLoading = ref(false);

const fetchProjects = async () => {
    loading.value = true;
    try {
        const response = await api.listProjects({ limit: limit.value, offset: offset.value });
        projects.value = response.data.items || [];
    } catch (err) {
        console.error('Failed to fetch projects', err);
    } finally {
        loading.value = false;
    }
};

onMounted(fetchProjects);

const nameError = ref('');

const handleCreateProject = async () => {
    if (!newProject.value.name.trim()) {
        nameError.value = 'Project name is required';
        return;
    }
    nameError.value = '';
    createLoading.value = true;
    try {
        await api.createProject({ projectInputBody: { name: newProject.value.name, description: newProject.value.description } });
        showCreateDialog.value = false;
        newProject.value = { name: '', description: '' };
        await fetchProjects();
    } catch (e) {
        console.error("Create project failed", e);
    } finally {
        createLoading.value = false;
    }
};
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950 p-6">
    <Toolbar class="mb-8 p-4 rounded-xl shadow-sm">
      <template #start>
        <span class="text-2xl font-bold text-primary px-4">Projects</span>
      </template>
      <template #end>
        <div class="flex gap-2">
            <Button icon="pi pi-plus" label="New Project" severity="info" text @click="showCreateDialog = true" />
            <router-link to="/dashboard">
                <Button icon="pi pi-home" label="Dashboard" text />
            </router-link>
        </div>
      </template>
    </Toolbar>

    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 animate-pulse">
        <div v-for="i in 3" :key="i" class="h-48 bg-gray-200 dark:bg-gray-800 rounded-xl"></div>
    </div>

    <div v-else>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <Card v-for="project in projects" :key="project.id" class="shadow-sm hover:shadow-md transition-shadow cursor-pointer border-t-4 border-primary-500" @click="$router.push(`/projects/${project.id}`)">
                <template #title>{{ project.name }}</template>
                <template #content>
                    <p class="text-gray-600 dark:text-gray-400 line-clamp-3 text-sm h-12">{{ project.description || 'No description provided.' }}</p>
                </template>
                <template #footer>
                    <div class="flex justify-end">
                        <Button label="View Details" icon="pi pi-external-link" text size="small" />
                    </div>
                </template>
            </Card>
        </div>

        <div v-if="projects.length === 0" class="text-center text-gray-500 py-12 italic">
            No projects found. Create your first project to get started!
        </div>

        <!-- Pagination Simple -->
        <div class="flex justify-center gap-2 mt-12">
            <Button icon="pi pi-chevron-left" :disabled="offset === 0" @click="offset -= limit; fetchProjects()" text />
            <Button icon="pi pi-chevron-right" :disabled="projects.length < limit" @click="offset += limit; fetchProjects()" text />
        </div>
    </div>

    <!-- Create Project Dialog -->
    <Dialog v-model:visible="showCreateDialog" header="Create New Project" :style="{ width: '450px' }" :breakpoints="{ '640px': '92vw' }" modal>
        <div class="flex flex-col gap-4">
            <div class="flex flex-col gap-2">
                <label for="name" class="font-bold">Project Name</label>
                <InputText id="name" v-model="newProject.name" placeholder="e.g. Website Redesign" :invalid="!!nameError" />
                <small v-if="nameError" class="text-red-500">{{ nameError }}</small>
            </div>
            <div class="flex flex-col gap-2">
                <label for="desc" class="font-bold">Description</label>
                <Textarea id="desc" v-model="newProject.description" rows="3" placeholder="What is this project about?" />
            </div>
            <div class="flex justify-end gap-2 mt-4">
                <Button label="Cancel" severity="secondary" text @click="showCreateDialog = false" />
                <Button label="Create Project" :loading="createLoading" @click="handleCreateProject" />
            </div>
        </div>
    </Dialog>
  </div>
</template>
