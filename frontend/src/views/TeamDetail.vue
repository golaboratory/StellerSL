<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { useToastStore } from '../stores/toast';
import { DefaultApi, TeamMemberUser } from '../api';
import axiosInstance from '../api/axios';
import Card from 'primevue/card';
import Button from 'primevue/button';
import Toolbar from 'primevue/toolbar';
import Tag from 'primevue/tag';
import InputText from 'primevue/inputtext';
import ConfirmDialog from 'primevue/confirmdialog';
import { useConfirm } from "primevue/useconfirm";

const route = useRoute();
const router = useRouter();
const auth = useAuthStore();
const toast = useToastStore();
const confirm = useConfirm();
const api = new DefaultApi(undefined, '/api', axiosInstance);

const teamID = route.params.id as string;
const team = ref<any>(null);
const members = ref<TeamMemberUser[]>([]);
const loading = ref(true);

const newMemberID = ref('');
const addLoading = ref(false);

const fetchTeamDetails = async () => {
    loading.value = true;
    try {
        const [teamsRes, membersRes] = await Promise.all([
            api.listTeams(),
            api.listTeamMembers(teamID)
        ]);
        team.value = teamsRes.data.items?.find((t: any) => t.id === teamID);
        members.value = membersRes.data.body.items || [];
    } catch (err) {
        console.error('Failed to fetch team details', err);
        toast.showToast('error', 'Error', 'Failed to load team details');
    } finally {
        loading.value = false;
    }
};

onMounted(fetchTeamDetails);

const handleAddMember = async () => {
    if (!newMemberID.value) return;
    addLoading.value = true;
    try {
        await api.addTeamMember({ id: teamID, teamMemberInputBody: { user_id: newMemberID.value, role: 'member' } });
        newMemberID.value = '';
        toast.showToast('success', 'Success', 'Member added to team');
        await fetchTeamDetails();
    } catch (e) {
        console.error("Failed to add member", e);
        toast.showToast('error', 'Error', 'Failed to add member');
    } finally {
        addLoading.value = false;
    }
};

const removeMember = (userID: string, userName: string) => {
    confirm.require({
        message: `Are you sure you want to remove ${userName} from the team?`,
        header: 'Confirm Removal',
        icon: 'pi pi-exclamation-triangle',
        acceptProps: { label: 'Remove', severity: 'danger' },
        accept: async () => {
            try {
                await api.removeTeamMember(teamID, userID);
                toast.showToast('success', 'Success', 'Member removed');
                await fetchTeamDetails();
            } catch (err) {
                toast.showToast('error', 'Error', 'Failed to remove member');
            }
        }
    });
};

const copyInviteLink = () => {
    const url = `${window.location.origin}/signup?team=${teamID}`;
    navigator.clipboard.writeText(url).then(() => {
        toast.showToast('success', 'Copied', 'Invitation link copied to clipboard');
    });
};
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950 p-6">
    <ConfirmDialog />
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
        <template #title>Members ({{ members.length }})</template>
        <template #content>
            <div class="space-y-4">
                <div v-for="member in members" :key="member.id" class="p-4 border rounded-lg flex items-center justify-between bg-white dark:bg-gray-900 border-gray-100 dark:border-gray-800">
                    <div class="flex items-center gap-3">
                        <i class="pi pi-user text-xl text-gray-400"></i>
                        <div>
                            <div class="font-medium">{{ member.name }}</div>
                            <div class="text-xs text-gray-500">{{ member.email }}</div>
                        </div>
                    </div>
                    <div class="flex items-center gap-2">
                        <Tag value="MEMBER" severity="info" />
                        <Button icon="pi pi-user-minus" severity="danger" text rounded @click="removeMember(member.id, member.name)" v-if="member.id !== auth.user?.id" />
                    </div>
                </div>
                <div v-if="members.length === 0" class="text-center py-8 text-gray-500 italic">
                    No members found.
                </div>
            </div>
        </template>
      </Card>

      <div class="md:col-span-1 space-y-6">
        <Card class="shadow-sm">
            <template #title>Team Invitation</template>
            <template #content>
                <div class="flex flex-col gap-4">
                    <p class="text-sm text-gray-500">Share this link with others to invite them to your team.</p>
                    <Button icon="pi pi-copy" label="Copy Invite Link" severity="secondary" @click="copyInviteLink" />
                </div>
            </template>
        </Card>

        <Card class="shadow-sm">
            <template #title>Add Member</template>
            <template #content>
                <div class="flex flex-col gap-4">
                    <div class="flex flex-col gap-2">
                        <label for="memberID" class="text-sm font-medium">User ID</label>
                        <InputText id="memberID" v-model="newMemberID" placeholder="UUID of the user" />
                    </div>
                    <Button label="Add to Team" :loading="addLoading" @click="handleAddMember" />
                </div>
            </template>
        </Card>
      </div>
    </div>
  </div>
</template>
