<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useAuthStore } from '../stores/auth';
import { DefaultApi, Configuration, TaskItem } from '../api';
import axiosInstance from '../api/axios';
import Card from 'primevue/card';
import Toolbar from 'primevue/toolbar';
import Button from 'primevue/button';
import Tag from 'primevue/tag';
import Checkbox from 'primevue/checkbox';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Textarea from 'primevue/textarea';
import SelectButton from 'primevue/selectbutton';
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';

const auth = useAuthStore();
const api = new DefaultApi(undefined, '/api', axiosInstance);

const tasks = ref<TaskItem[]>([]);
const loading = ref(true);
const selectedTasks = ref<string[]>([]);
const searchQuery = ref('');

// Pagination
const limit = ref(50);
const offset = ref(0);

// Bulk Create Dialog
const showBulkDialog = ref(false);
const bulkText = ref('');
const bulkLoading = ref(false);

// Edit Task Dialog
const showEditDialog = ref(false);
const editTask = ref<any>(null);
const editLoading = ref(false);

const filteredTasks = computed(() => {
    if (!searchQuery.value) return tasks.value;
    const query = searchQuery.value.toLowerCase();
    return tasks.value.filter(t => 
        t.title.toLowerCase().includes(query) || 
        (t as any).description?.toLowerCase().includes(query)
    );
});

const openEditDialog = (task: any) => {
    editTask.value = { ...task };
    showEditDialog.value = true;
};

const handleUpdateTask = async () => {
    if (!editTask.value) return;
    editLoading.value = true;
    try {
        await api.updateTask({ 
            id: editTask.value.id, 
            taskInput: {
                title: editTask.value.title,
                status: editTask.value.status,
                priority: editTask.value.priority,
                description: editTask.value.description || "",
                project_id: editTask.value.project_id
            }
        });
        showEditDialog.value = false;
        await fetchTasks();
    } catch (e) {
        console.error("Update failed", e);
    } finally {
        editLoading.value = false;
    }
};

const fetchTasks = async () => {
    loading.value = true;
    try {
        const response = await api.listTasks({ limit: limit.value, offset: offset.value });
        tasks.value = response.data.items || [];
        selectedTasks.value = [];
    } catch (err) {
        console.error('Failed to fetch tasks', err);
    } finally {
        loading.value = false;
    }
};

onMounted(fetchTasks);

const getStatusSeverity = (status: string) => {
    switch (status) {
        case 'todo': return 'secondary';
        case 'doing': return 'warn';
        case 'done': return 'success';
        default: return 'info';
    }
};

const updateStatus = async (task: any, newStatus: string) => {
    try {
        await api.updateTaskStatus({ 
            id: task.id, 
            taskStatusUpdateInputBody: { status: newStatus } 
        });
        task.status = newStatus;
        if (newStatus === 'done') {
            await fetchTasks(); // Refresh to update EXP/Badges
        }
    } catch(e) {
        console.error("Failed to update status", e);
    }
};

const handleBulkStatus = async (newStatus: string) => {
    if (selectedTasks.value.length === 0) return;
    try {
        await api.bulkUpdateTasksStatus({ 
            bulkTaskUpdateInputBody: { ids: selectedTasks.value, status: newStatus } 
        });
        await fetchTasks();
    } catch (e) {
        console.error("Bulk update failed", e);
    }
};

const handleBulkDelete = async () => {
    if (selectedTasks.value.length === 0 || !confirm(`Delete ${selectedTasks.value.length} tasks?`)) return;
    try {
        await api.bulkDeleteTasks({ 
            bulkTaskDeleteInputBody: { ids: selectedTasks.value } 
        });
        await fetchTasks();
    } catch (e) {
        console.error("Bulk delete failed", e);
    }
};

const handleBulkCreate = async () => {
    const lines = bulkText.value.split('\n').map(l => l.trim()).filter(l => l.length > 0);
    if (lines.length === 0) return;

    bulkLoading.value = true;
    try {
        const tasksToCreate = lines.map(title => ({
            title,
            project_id: "", 
            description: ""
        }));
        await api.bulkCreateTasks({ 
            bulkTaskCreateInputBody: { tasks: tasksToCreate } 
        });
        showBulkDialog.value = false;
        bulkText.value = '';
        await fetchTasks();
    } catch (e) {
        console.error("Bulk create failed", e);
    } finally {
        bulkLoading.value = false;
    }
};
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950 p-6">
    <Toolbar class="mb-8 p-4 rounded-xl shadow-sm">
      <template #start>
        <div class="flex items-center gap-4">
            <span class="text-2xl font-bold text-primary px-4 hidden md:block">Tasks (GTD)</span>
            <IconField iconPosition="left">
                <InputIcon class="pi pi-search" />
                <InputText v-model="searchQuery" placeholder="Search tasks..." class="w-full md:w-80" />
            </IconField>
        </div>
      </template>
      <template #end>
        <div class="flex gap-2">
            <Button icon="pi pi-plus" label="Bulk Add" severity="info" text @click="showBulkDialog = true" />
            <router-link to="/dashboard">
                <Button icon="pi pi-home" label="Dashboard" text />
            </router-link>
        </div>
      </template>
    </Toolbar>

    <!-- Bulk Actions Bar -->
    <div v-if="selectedTasks.length > 0" class="mb-4 p-4 bg-primary-50 dark:bg-primary-900/20 rounded-lg flex items-center justify-between border border-primary-100 dark:border-primary-800 animate-fadein">
        <span class="font-bold">{{ selectedTasks.length }} tasks selected</span>
        <div class="flex gap-2">
            <Button label="Set Todo" severity="secondary" size="small" @click="handleBulkStatus('todo')" />
            <Button label="Set Doing" severity="warn" size="small" @click="handleBulkStatus('doing')" />
            <Button label="Set Done" severity="success" size="small" @click="handleBulkStatus('done')" />
            <Button icon="pi pi-trash" severity="danger" size="small" @click="handleBulkDelete" />
        </div>
    </div>

    <div v-if="loading" class="space-y-4 animate-pulse">
        <div v-for="i in 3" :key="i" class="h-24 bg-gray-200 dark:bg-gray-800 rounded-xl"></div>
    </div>

    <div v-else class="space-y-4">
      <Card v-for="task in filteredTasks" :key="task.id" class="shadow-sm transition-all" :class="{'border-primary-500 bg-primary-50/50': selectedTasks.includes(task.id)}">
        <template #content>
            <div class="flex items-center gap-4">
                <Checkbox v-model="selectedTasks" :value="task.id" />
                <div class="flex-1 flex items-center justify-between">
                    <div>
                        <div class="flex items-center gap-2">
                            <h3 class="text-lg font-bold">{{ task.title }}</h3>
                            <i v-if="task.priority > 0" class="pi pi-exclamation-circle text-orange-500" title="High Priority"></i>
                        </div>
                        <div class="flex items-center gap-2 mt-2">
                            <Tag :value="task.status.toUpperCase()" :severity="getStatusSeverity(task.status)" />
                            <span v-if="task.due_date && task.due_date !== '0001-01-01 00:00:00 +0000 UTC'" class="text-xs text-gray-500">
                                <i class="pi pi-calendar mr-1"></i>{{ new Date(task.due_date).toLocaleDateString() }}
                            </span>
                        </div>
                    </div>
                    <div class="flex gap-2">
                        <Button icon="pi pi-pencil" severity="secondary" text @click="openEditDialog(task)" />
                        <Button v-if="task.status !== 'todo'" icon="pi pi-step-backward" text @click="updateStatus(task, task.status === 'done' ? 'doing' : 'todo')" />
                        <Button v-if="task.status !== 'done'" icon="pi pi-check" severity="success" text @click="updateStatus(task, task.status === 'todo' ? 'doing' : 'done')" />
                    </div>
                </div>
            </div>
        </template>
      </Card>
      
      <div v-if="filteredTasks.length === 0" class="text-center text-gray-500 py-8">
        {{ searchQuery ? 'No tasks match your search.' : 'No tasks found.' }}
      </div>

      <!-- Pagination Simple -->
      <div class="flex justify-center gap-2 mt-8">
          <Button icon="pi pi-chevron-left" :disabled="offset === 0" @click="offset -= limit; fetchTasks()" text />
          <Button icon="pi pi-chevron-right" :disabled="tasks.length < limit" @click="offset += limit; fetchTasks()" text />
      </div>
    </div>

    <!-- Edit Task Dialog -->
    <Dialog v-model:visible="showEditDialog" header="Edit Task" :style="{ width: '450px' }" modal>
        <div v-if="editTask" class="flex flex-col gap-4">
            <div class="flex flex-col gap-2">
                <label for="title" class="font-bold">Title</label>
                <InputText id="title" v-model="editTask.title" />
            </div>
            <div class="flex flex-col gap-2">
                <label for="status" class="font-bold">Status</label>
                <SelectButton v-model="editTask.status" :options="['todo', 'doing', 'done']" class="uppercase" />
            </div>
            <div class="flex flex-col gap-2">
                <label for="priority" class="font-bold">Priority</label>
                <SelectButton v-model="editTask.priority" :options="[0, 1]" :optionLabel="(opt) => opt === 1 ? 'High' : 'Normal'" />
            </div>
            <div class="flex flex-col gap-2">
                <label for="description" class="font-bold">Description</label>
                <Textarea id="description" v-model="editTask.description" rows="3" />
            </div>
            <div class="flex justify-end gap-2 mt-4">
                <Button label="Cancel" severity="secondary" text @click="showEditDialog = false" />
                <Button label="Save Changes" :loading="editLoading" @click="handleUpdateTask" />
            </div>
        </div>
    </Dialog>

    <!-- Bulk Create Dialog -->
    <Dialog v-model:visible="showBulkDialog" header="Bulk Add Tasks" :style="{ width: '500px' }" modal>
        <div class="flex flex-col gap-4">
            <label>Enter one task title per line:</label>
            <Textarea v-model="bulkText" rows="10" class="w-full font-mono" placeholder="Task 1&#10;Task 2&#10;Task 3" />
            <div class="flex justify-end gap-2">
                <Button label="Cancel" severity="secondary" text @click="showBulkDialog = false" />
                <Button label="Create Tasks" :loading="bulkLoading" @click="handleBulkCreate" />
            </div>
        </div>
    </Dialog>
  </div>
</template>
