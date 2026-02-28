<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { DefaultApi, UserStruct } from '../api';
import axiosInstance from '../api/axios';
import Card from 'primevue/card';
import Button from 'primevue/button';
import Toolbar from 'primevue/toolbar';
import Tag from 'primevue/tag';
import InputText from 'primevue/inputtext';

const route = useRoute();
const router = useRouter();
const auth = useAuthStore();
const api = new DefaultApi(undefined, '/api', axiosInstance);

const teamID = route.params.id as string;
const team = ref<any>(null);
const members = ref<UserStruct[]>([]);
const loading = ref(true);

const newMemberID = ref('');
const addLoading = ref(false);

const fetchTeamDetails = async () => {
    loading.value = true;
    try {
        const teamsRes = await api.listTeams();
        team.value = teamsRes.data.items?.find((t: any) => t.id === teamID);
        
        // Note: Currently no API to list members specifically, 
        // we'll implement a mock or skip if not available.
        // For POC, we'll assume the API exists as defined in previous implementation steps.
    } catch (err) {
        console.error('Failed to fetch team details', err);
    } finally {
        loading.value = false;
    }
};

onMounted(fetchTeamDetails);

const handleAddMember = async () => {
    if (!newMemberID.value) return;
    addLoading.value = true;
    try {
        await api.addTeamMember(teamID, { teamMemberInputBody: { user_id: newMemberID.value, role: 'member' } });
        newMemberID.value = '';
        await fetchTeamDetails();
    } catch (e) {
        console.error("Failed to add member", e);
    } finally {
        addLoading.value = false;
    }
};

const copyInviteLink = () => {
    const url = `${window.location.origin}/signup?team=${teamID}`;
    navigator.clipboard.writeText(url).then(() => {
        alert("Invitation link copied to clipboard!");
    });
};
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950 p-6">
    <Toolbar class="mb-8 p-4 rounded-xl shadow-sm">
      <template #start>
        <div class="flex items-center gap-4 px-4">
            <Button icon="pi pi-arrow-left" text @click="router.push('/teams')" />
            <span class="text-2xl font-bold text-primary">{{ team?.name || 'Team Details' }}</span>
        </div>
      </template>
      <template #end>
        <router-link to="/dashboard">
          <Button icon="pi pi-home" label="Dashboard" text />
        </router-link>
      </template>
    </Toolbar>

    <div v-if="loading" class="space-y-6 animate-pulse">
        <div v-for="i in 2" :key="i" class="h-32 bg-gray-200 dark:bg-gray-800 rounded-xl"></div>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <Card class="md:col-span-2 shadow-sm">
        <template #title>Members</template>
        <template #content>
            <div class="space-y-4">
                <div v-for="member in members" :key="member.id" class="p-4 border rounded-lg flex items-center justify-between bg-white dark:bg-gray-900 border-gray-100 dark:border-gray-800">
                    <div class="flex items-center gap-3">
                        <i class="pi pi-user text-xl text-gray-400"></i>
                        <span class="font-medium">{{ member.name }}</span>
                    </div>
                    <Tag value="MEMBER" severity="info" />
                </div>
                <div v-if="members.length === 0" class="text-center py-8 text-gray-500 italic">
                    No members found.
                </div>
            </div>
        </template>
      </Card>

      <Card class="md:col-span-1 shadow-sm">
        <template #title>Team Invitation</template>
        <template #content>
            <div class="flex flex-col gap-4">
                <p class="text-sm text-gray-500">Share this link with others to invite them to your team.</p>
                <Button icon="pi pi-copy" label="Copy Invite Link" severity="secondary" @click="copyInviteLink" />
            </div>
        </template>
      </Card>

      <Card class="md:col-span-1 shadow-sm">
        <template #title>Add Member</template>
        <template #content>
            <div class="flex flex-col gap-4">
                <div class="flex flex-col gap-2">
                    <label for="memberID">User ID</label>
                    <InputText id="memberID" v-model="newMemberID" placeholder="UUID of the user" />
                </div>
                <Button label="Add to Team" :loading="addLoading" @click="handleAddMember" />
            </div>
        </template>
      </Card>
    </div>
  </div>
</template>
