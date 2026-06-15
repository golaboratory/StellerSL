<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { useToastStore } from '../stores/toast';
import { DefaultApi } from '../api';
import type { TeamMemberUser, AccountUser } from '../api';
import axiosInstance from '../api/axios';
import Card from 'primevue/card';
import Button from 'primevue/button';
import Toolbar from 'primevue/toolbar';
import Tag from 'primevue/tag';
import InputText from 'primevue/inputtext';
import Dialog from 'primevue/dialog';
import Select from 'primevue/select';
import ConfirmDialog from 'primevue/confirmdialog';
import { useConfirm } from "primevue/useconfirm";
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';

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

// Member search dialog (6-6 / 7-FE-1)
const addDialog = ref(false);
const userSearchQuery = ref('');
const searchResults = ref<AccountUser[]>([]);
const searchLoading = ref(false);
const addLoading = ref(false);
const newMemberRole = ref<'member' | 'admin'>('member');

// Team rename dialog (7-FE-2)
const renameDialog = ref(false);
const renameValue = ref('');
const renameLoading = ref(false);

const roleOptions = [
    { label: 'Admin', value: 'admin' },
    { label: 'Member', value: 'member' },
];

const myRole = computed(() => members.value.find(m => m.id === auth.user?.id)?.role || '');
const canManage = computed(() => myRole.value === 'owner' || myRole.value === 'admin');
const isOwner = computed(() => myRole.value === 'owner');

const roleSeverity = (role: string) => {
    switch (role) {
        case 'owner': return 'warn';
        case 'admin': return 'info';
        default: return 'secondary';
    }
};

const fetchTeamDetails = async () => {
    loading.value = true;
    try {
        const [teamRes, membersRes] = await Promise.all([
            api.getTeam({ id: teamID }),
            api.listTeamMembers({ id: teamID })
        ]);
        team.value = teamRes.data;
        members.value = membersRes.data.items || [];
    } catch (err) {
        console.error('Failed to fetch team details', err);
        toast.showToast('error', 'Error', 'Failed to load team details');
    } finally {
        loading.value = false;
    }
};

onMounted(fetchTeamDetails);

const handleSearchUsers = async () => {
    if (userSearchQuery.value.length < 2) {
        searchResults.value = [];
        return;
    }
    searchLoading.value = true;
    try {
        const res = await api.searchUsers({ q: userSearchQuery.value });
        searchResults.value = res.data.items || [];
    } catch (e) {
        console.error("Search failed", e);
    } finally {
        searchLoading.value = false;
    }
};

const handleAddMember = async (userID: string) => {
    addLoading.value = true;
    try {
        await api.addTeamMember({
            id: teamID,
            teamMemberInputBody: { user_id: userID, role: newMemberRole.value }
        });
        toast.showToast('success', 'Success', 'Member added to team');
        await fetchTeamDetails();
    } catch (e) {
        console.error("Failed to add member", e);
    } finally {
        addLoading.value = false;
    }
};

// AddTeamMember upserts, so changing a role reuses the same endpoint
const handleRoleChange = async (member: TeamMemberUser, newRole: string) => {
    try {
        await api.addTeamMember({
            id: teamID,
            teamMemberInputBody: { user_id: member.id, role: newRole as any }
        });
        toast.showToast('success', 'Success', `${member.name} is now ${newRole}`);
        await fetchTeamDetails();
    } catch (e) {
        console.error("Failed to change role", e);
        await fetchTeamDetails();
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
                await api.removeTeamMember({ id: teamID, userId: userID });
                toast.showToast('success', 'Success', 'Member removed');
                await fetchTeamDetails();
            } catch (err) {
                toast.showToast('error', 'Error', 'Failed to remove member');
            }
        }
    });
};

const openRenameDialog = () => {
    renameValue.value = team.value?.name || '';
    renameDialog.value = true;
};

const handleRename = async () => {
    if (!renameValue.value.trim()) return;
    renameLoading.value = true;
    try {
        await api.updateTeam({ id: teamID, updateTeamRequest: { name: renameValue.value.trim() } });
        renameDialog.value = false;
        toast.showToast('success', 'Success', 'Team renamed');
        await fetchTeamDetails();
    } catch (e) {
        console.error("Rename failed", e);
    } finally {
        renameLoading.value = false;
    }
};

const deleteTeam = () => {
    confirm.require({
        message: 'Are you sure you want to delete this team? This cannot be undone.',
        header: 'Confirm Deletion',
        icon: 'pi pi-exclamation-triangle',
        acceptProps: { label: 'Delete', severity: 'danger' },
        accept: async () => {
            try {
                await api.deleteTeam({ id: teamID });
                toast.showToast('success', 'Success', 'Team deleted');
                router.push('/teams');
            } catch (err) {
                console.error('Delete failed', err);
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

    <!-- Add Member Dialog (user search) -->
    <Dialog v-model:visible="addDialog" header="Add Team Member" :modal="true" class="w-full max-w-md">
        <div class="flex flex-col gap-4">
            <div class="flex flex-col gap-2">
                <label for="searchUser" class="font-medium">Search User by Name or Email</label>
                <IconField iconPosition="left">
                    <InputIcon class="pi pi-search" />
                    <InputText id="searchUser" v-model="userSearchQuery" placeholder="Type to search..." class="w-full" @input="handleSearchUsers" />
                </IconField>
            </div>

            <div class="flex items-center gap-2">
                <label class="text-sm font-medium">Add as:</label>
                <Select v-model="newMemberRole" :options="roleOptions" optionLabel="label" optionValue="value" class="w-32" size="small" />
            </div>

            <div class="max-h-64 overflow-y-auto border rounded-lg border-gray-100 dark:border-gray-800">
                <div v-if="searchLoading" class="p-4 text-center"><i class="pi pi-spin pi-spinner mr-2"></i>Searching...</div>
                <div v-else-if="searchResults.length === 0 && userSearchQuery.length >= 2" class="p-4 text-center text-gray-500">No users found.</div>
                <div v-else-if="userSearchQuery.length < 2" class="p-4 text-center text-gray-400 text-sm italic">Enter at least 2 characters to search.</div>

                <div v-for="user in searchResults" :key="user.id" class="flex items-center justify-between p-3 hover:bg-gray-50 dark:hover:bg-gray-900 border-b last:border-0 border-gray-100 dark:border-gray-800">
                    <div class="flex flex-col">
                        <span class="text-sm font-bold">{{ user.name }}</span>
                        <span class="text-xs text-gray-500">{{ user.email }}</span>
                    </div>
                    <Button icon="pi pi-plus" size="small" rounded text :loading="addLoading"
                            :disabled="members.some(m => m.id === user.id)"
                            @click="handleAddMember(user.id)" />
                </div>
            </div>
        </div>
    </Dialog>

    <!-- Rename Team Dialog -->
    <Dialog v-model:visible="renameDialog" header="Rename Team" :style="{ width: '400px' }" :breakpoints="{ '640px': '92vw' }" modal>
        <div class="flex flex-col gap-4">
            <div class="flex flex-col gap-2">
                <label for="teamName" class="font-bold">Team Name</label>
                <InputText id="teamName" v-model="renameValue" :invalid="!renameValue.trim()" />
                <small v-if="!renameValue.trim()" class="text-red-500">Team name is required</small>
            </div>
            <div class="flex justify-end gap-2 mt-2">
                <Button label="Cancel" severity="secondary" text @click="renameDialog = false" />
                <Button label="Save" :loading="renameLoading" :disabled="!renameValue.trim()" @click="handleRename" />
            </div>
        </div>
    </Dialog>

    <Toolbar class="mb-8 p-4 rounded-xl shadow-sm">
      <template #start>
        <div class="flex items-center gap-4 px-4">
            <Button icon="pi pi-arrow-left" text @click="router.push('/teams')" />
            <span class="text-2xl font-bold text-primary">{{ team?.name || 'Team Details' }}</span>
            <Tag v-if="myRole" :value="myRole.toUpperCase()" :severity="roleSeverity(myRole)" />
        </div>
      </template>
      <template #end>
        <div class="flex gap-2">
            <Button v-if="canManage" icon="pi pi-pencil" label="Rename" text @click="openRenameDialog" />
            <Button v-if="isOwner" icon="pi pi-trash" label="Delete Team" severity="danger" text @click="deleteTeam" />
            <router-link to="/dashboard">
              <Button icon="pi pi-home" label="Dashboard" text />
            </router-link>
        </div>
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
                <div v-for="member in members" :key="member.id" class="p-4 border rounded-lg flex items-center justify-between flex-wrap gap-2 bg-white dark:bg-gray-900 border-gray-100 dark:border-gray-800">
                    <div class="flex items-center gap-3 min-w-0">
                        <i class="pi pi-user text-xl text-gray-400"></i>
                        <div>
                            <div class="font-medium">{{ member.name }} <span v-if="member.id === auth.user?.id" class="text-xs text-gray-400">(you)</span></div>
                            <div class="text-xs text-gray-500">{{ member.email }}</div>
                        </div>
                    </div>
                    <div class="flex items-center gap-2">
                        <!-- Owner role is fixed; admins/owners can change other members' roles -->
                        <Select v-if="canManage && member.role !== 'owner'"
                                :modelValue="member.role"
                                :options="roleOptions" optionLabel="label" optionValue="value"
                                size="small" class="w-32"
                                @update:modelValue="(v: string) => handleRoleChange(member, v)" />
                        <Tag v-else :value="(member.role || 'member').toUpperCase()" :severity="roleSeverity(member.role || 'member')" />
                        <Button icon="pi pi-user-minus" severity="danger" text rounded
                                v-if="canManage && member.id !== auth.user?.id && member.role !== 'owner'"
                                @click="removeMember(member.id, member.name)" />
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

        <Card class="shadow-sm" v-if="canManage">
            <template #title>Add Member</template>
            <template #content>
                <div class="flex flex-col gap-4">
                    <p class="text-sm text-gray-500">Search registered users by name or email and add them to this team.</p>
                    <Button icon="pi pi-user-plus" label="Search & Add Member" @click="addDialog = true" />
                </div>
            </template>
        </Card>
      </div>
    </div>
  </div>
</template>
