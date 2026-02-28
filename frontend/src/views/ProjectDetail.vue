<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { useToastStore } from '../stores/toast';
import { DefaultApi, TaskItem, AccountUser } from '../api';
import axiosInstance from '../api/axios';
import Card from 'primevue/card';
import Button from 'primevue/button';
import Toolbar from 'primevue/toolbar';
import Tag from 'primevue/tag';
import InputText from 'primevue/inputtext';
import Dialog from 'primevue/dialog';
import ConfirmDialog from 'primevue/confirmdialog';
import { useConfirm } from "primevue/useconfirm";

const route = useRoute();
const router = useRouter();
const auth = useAuthStore();
const toast = useToastStore();
const confirm = useConfirm();
const api = new DefaultApi(undefined, '/api', axiosInstance);

const projectID = route.params.id as string;
const project = ref<any>(null);
const tasks = ref<TaskItem[]>([]);
const members = ref<AccountUser[]>([]);
const loading = ref(true);

const assignDialog = ref(false);
const newMemberID = ref('');
const assignLoading = ref(false);

const fetchProjectDetails = async () => {
    loading.value = true;
    try {
        const [projRes, tasksRes, membersRes] = await Promise.all([
            api.getProject({ id: projectID }),
            api.listProjectTasks({ id: projectID }),
            api.listProjectUsers({ id: projectID })
        ]);
        
        project.value = projRes.data;
        tasks.value = tasksRes.data.items || [];
        members.value = membersRes.data.items || [];
    } catch (err) {
        console.error('Failed to fetch project details', err);
        toast.showToast('error', 'Error', 'Failed to load project details');
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

const handleAssignMember = async () => {
    if (!newMemberID.value) return;
    assignLoading.value = true;
    try {
        await api.assignProjectUser({ 
            id: projectID, 
            projectUserAssignmentInputBody: { user_id: newMemberID.value } 
        });
        newMemberID.value = '';
        assignDialog.value = false;
        toast.showToast('success', 'Success', 'User assigned to project');
        await fetchProjectDetails();
    } catch (e) {
        toast.showToast('error', 'Error', 'Failed to assign user');
    } finally {
        assignLoading.value = false;
    }
};

const unassignMember = (userID: string, userName: string) => {
    confirm.require({
        message: `Remove ${userName} from this project?`,
        header: 'Confirm Unassignment',
        icon: 'pi pi-exclamation-triangle',
        acceptProps: { label: 'Remove', severity: 'danger' },
        accept: async () => {
            try {
                await api.unassignProjectUser({ id: projectID, userId: userID });
                toast.showToast('success', 'Success', 'User removed from project');
                await fetchProjectDetails();
            } catch (err) {
                toast.showToast('error', 'Error', 'Failed to remove user');
            }
        }
    });
};

const deleteProject = () => {
    confirm.require({
        message: 'Are you sure you want to delete this project? All associated tasks will be unlinked.',
        header: 'Confirm Deletion',
        icon: 'pi pi-exclamation-triangle',
        acceptProps: { label: 'Delete', severity: 'danger' },
        accept: async () => {
            try {
                await api.deleteProject({ id: projectID });
                toast.showToast('success', 'Success', 'Project deleted successfully');
                router.push('/projects');
            } catch (err) {
                toast.showToast('error', 'Error', 'Failed to delete project');
            }
        }
    });
};
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950 p-6">
    <ConfirmDialog />
    
    <Dialog v-model:visible="assignDialog" header="Assign Member" :modal="true" class="w-full max-w-md">
        <div class="flex flex-col gap-4">
            <div class="flex flex-col gap-2">
                <label for="assignUserID" class="font-medium">User ID (UUID)</label>
                <InputText id="assignUserID" v-model="newMemberID" placeholder="Enter User ID" />
            </div>
            <Button label="Assign" :loading="assignLoading" @click="handleAssignMember" />
        </div>
    </Dialog>

    <Toolbar class="mb-8 p-4 rounded-xl shadow-sm">
      <template #start>
        <div class="flex items-center gap-4 px-4">
            <Button icon="pi pi-arrow-left" text @click="router.push('/projects')" />
            <span class="text-2xl font-bold text-primary">{{ project?.name || 'Project Details' }}</span>
        </div>
      </template>
      <template #end>
        <div class="flex gap-2">
            <Button icon="pi pi-trash" label="Delete Project" severity="danger" text @click="deleteProject" />
            <router-link to="/dashboard">
                <Button icon="pi pi-home" label="Dashboard" text />
            </router-link>
        </div>
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
            <template #title>Project Members ({{ members.length }})</template>
            <template #content>
                <div class="space-y-3">
                    <div v-for="user in members" :key="user.id" class="flex items-center justify-between p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors">
                        <div class="flex items-center gap-3">
                            <i class="pi pi-user text-gray-400"></i>
                            <div>
                                <div class="text-sm font-medium">{{ user.name }}</div>
                                <div class="text-xs text-gray-500">{{ user.email }}</div>
                            </div>
                        </div>
                        <Button icon="pi pi-user-minus" severity="danger" text rounded size="small" @click="unassignMember(user.id, user.name)" />
                    </div>
                    <div v-if="members.length === 0" class="text-center py-4 text-gray-500 italic text-sm">
                        No members assigned.
                    </div>
                    <Button label="Assign Member" icon="pi pi-user-plus" severity="secondary" outlined class="w-full mt-2" @click="assignDialog = true" />
                </div>
            </template>
        </Card>

        <Card class="shadow-sm">
            <template #title>Project Stats</template>
            <template #content>
                <div class="flex justify-around text-center h-full items-center py-4">
                    <div>
                        <div class="text-3xl font-bold">{{ tasks.length }}</div>
                        <div class="text-xs text-gray-500 uppercase tracking-wider">Tasks</div>
                    </div>
                    <div class="w-px h-12 bg-gray-200"></div>
                    <div>
                        <div class="text-3xl font-bold text-green-500">{{ tasks.filter(t => t.status === 'done').length }}</div>
                        <div class="text-xs text-gray-500 uppercase tracking-wider">Completed</div>
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
