<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { DefaultApi } from '../api';
import type { TeamItem } from '../api';
import axiosInstance from '../api/axios';
import Card from 'primevue/card';
import Toolbar from 'primevue/toolbar';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';

const api = new DefaultApi(undefined, '/api', axiosInstance);

const teams = ref<TeamItem[]>([]);
const loading = ref(true);

// Pagination
const limit = ref(50);
const offset = ref(0);

// Create Team Dialog
const showCreateDialog = ref(false);
const newTeamName = ref('');
const createLoading = ref(false);

const fetchTeams = async () => {
    loading.value = true;
    try {
        const response = await api.listTeams({ limit: limit.value, offset: offset.value });
        teams.value = response.data.items || [];
    } catch (err) {
        console.error('Failed to fetch teams', err);
    } finally {
        loading.value = false;
    }
};

onMounted(fetchTeams);

const handleCreateTeam = async () => {
    if (!newTeamName.value) return;
    createLoading.value = true;
    try {
        await api.createTeam({ teamInputBody: { name: newTeamName.value } });
        showCreateDialog.value = false;
        newTeamName.value = '';
        await fetchTeams();
    } catch (e) {
        console.error("Create team failed", e);
    } finally {
        createLoading.value = false;
    }
};
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950 p-6">
    <Toolbar class="mb-8 p-4 rounded-xl shadow-sm">
      <template #start>
        <span class="text-2xl font-bold text-primary px-4">Teams</span>
      </template>
      <template #end>
        <div class="flex gap-2">
            <Button icon="pi pi-plus" label="New Team" severity="info" text @click="showCreateDialog = true" />
            <router-link to="/dashboard">
                <Button icon="pi pi-home" label="Dashboard" text />
            </router-link>
        </div>
      </template>
    </Toolbar>

    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 animate-pulse">
        <div v-for="i in 3" :key="i" class="h-32 bg-gray-200 dark:bg-gray-800 rounded-xl"></div>
    </div>

    <div v-else>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <Card v-for="team in teams" :key="team.id" class="shadow-sm hover:shadow-md transition-shadow cursor-pointer border-l-4 border-primary-500" @click="$router.push(`/teams/${team.id}`)">
                <template #title>{{ team.name }}</template>
                <template #footer>
                    <div class="flex justify-end">
                        <Button label="Manage Members" icon="pi pi-users" text size="small" />
                    </div>
                </template>
            </Card>
        </div>

        <div v-if="teams.length === 0" class="text-center text-gray-500 py-12 italic">
            No teams found. Create a team to collaborate with others!
        </div>

        <!-- Pagination Simple -->
        <div class="flex justify-center gap-2 mt-12">
            <Button icon="pi pi-chevron-left" :disabled="offset === 0" @click="offset -= limit; fetchTeams()" text />
            <Button icon="pi pi-chevron-right" :disabled="teams.length < limit" @click="offset += limit; fetchTeams()" text />
        </div>
    </div>

    <!-- Create Team Dialog -->
    <Dialog v-model:visible="showCreateDialog" header="Create New Team" :style="{ width: '400px' }" :breakpoints="{ '640px': '92vw' }" modal>
        <div class="flex flex-col gap-4">
            <div class="flex flex-col gap-2">
                <label for="teamName" class="font-bold">Team Name</label>
                <InputText id="teamName" v-model="newTeamName" placeholder="e.g. Engineering Team" />
            </div>
            <div class="flex justify-end gap-2 mt-4">
                <Button label="Cancel" severity="secondary" text @click="showCreateDialog = false" />
                <Button label="Create Team" :loading="createLoading" @click="handleCreateTeam" />
            </div>
        </div>
    </Dialog>
  </div>
</template>
