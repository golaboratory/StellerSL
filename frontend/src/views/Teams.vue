<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useAuthStore } from '../stores/auth';
import { DefaultApi, TeamItem } from '../api';
import axiosInstance from '../api/axios';
import Card from 'primevue/card';
import Button from 'primevue/button';
import Toolbar from 'primevue/toolbar';
import InputText from 'primevue/inputtext';
import Dialog from 'primevue/dialog';
import { useRouter } from 'vue-router';

const auth = useAuthStore();
const router = useRouter();
const api = new DefaultApi(undefined, '/api', axiosInstance);

const teams = ref<TeamItem[]>([]);
const loading = ref(true);
const showCreateDialog = ref(false);
const newTeamName = ref('');

const fetchTeams = async () => {
    loading.value = true;
    try {
        const response = await api.listTeams();
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
    try {
        await api.createTeam({ teamInputBody: { name: newTeamName.value } });
        showCreateDialog.value = false;
        newTeamName.value = '';
        await fetchTeams();
    } catch (e) {
        console.error("Failed to create team", e);
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

    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-3 gap-6 animate-pulse">
        <div v-for="i in 3" :key="i" class="h-32 bg-gray-200 dark:bg-gray-800 rounded-xl"></div>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <Card v-for="team in teams" :key="team.id" class="shadow-sm cursor-pointer hover:shadow-md transition-shadow" @click="router.push(`/teams/${team.id}`)">
        <template #title>{{ team.name }}</template>
        <template #content>
          <div class="flex items-center gap-2 text-gray-500">
            <i class="pi pi-users"></i>
            <span>Manage members</span>
          </div>
        </template>
      </Card>
      
      <div v-if="teams.length === 0" class="col-span-full text-center text-gray-500 py-8">
        No teams found.
      </div>
    </div>

    <Dialog v-model:visible="showCreateDialog" header="Create New Team" :style="{ width: '400px' }" modal>
        <div class="flex flex-col gap-4">
            <div class="flex flex-col gap-2">
                <label for="teamName">Team Name</label>
                <InputText id="teamName" v-model="newTeamName" placeholder="Engineering, Marketing, etc." />
            </div>
            <div class="flex justify-end gap-2 mt-2">
                <Button label="Cancel" severity="secondary" text @click="showCreateDialog = false" />
                <Button label="Create" @click="handleCreateTeam" />
            </div>
        </div>
    </Dialog>
  </div>
</template>
