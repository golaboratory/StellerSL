<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import axiosInstance from '../api/axios';
import { DefaultApi } from '../api';
import { useToastStore } from '../stores/toast';
import { Calendar } from 'v-calendar';
import 'v-calendar/dist/style.css';
import Toolbar from 'primevue/toolbar';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Textarea from 'primevue/textarea';
import SelectButton from 'primevue/selectbutton';
import Select from 'primevue/select';
import Tag from 'primevue/tag';

const api = new DefaultApi(undefined, '/api', axiosInstance);
const toast = useToastStore();

const tasks = ref<any[]>([]);
const projects = ref<any[]>([]);
const selectedDate = ref(new Date());

// Create dialog state
const showCreateDialog = ref(false);
const newTask = ref({ title: '', description: '', priority: 0, project_id: '' });
const createLoading = ref(false);
const createError = ref('');

// Edit dialog state
const showEditDialog = ref(false);
const editTask = ref<any>(null);
const editLoading = ref(false);

const loadTasks = async () => {
    try {
        const response = await api.listTasks();
        tasks.value = response.data.items || [];
    } catch (err) {
        console.error(err);
    }
};

onMounted(async () => {
    await loadTasks();
    try {
        const res = await api.listProjects({ limit: 100, offset: 0 });
        projects.value = res.data.items || [];
    } catch (err) {
        console.error(err);
    }
});

const attributes = computed(() => tasks.value.map(t => ({
    key: t.id,
    dot: t.status === 'done' ? 'green' : (t.status === 'doing' ? 'orange' : 'blue'),
    customData: t,
    dates: [t.due_date ? new Date(t.due_date) : new Date()],
    popover: {
        label: t.title,
        visibility: 'hover' as const
    }
})));

const selectedTasks = computed(() => tasks.value.filter(t => {
    const d = t.due_date ? new Date(t.due_date) : new Date();
    return d.toDateString() === selectedDate.value.toDateString();
}));

const onDayClick = (day: any) => {
    selectedDate.value = day.date;
};

// Due date is the selected day at 17:00 local (end of work day)
const selectedDateAsDue = () => {
    const d = new Date(selectedDate.value);
    d.setHours(17, 0, 0, 0);
    return d.toISOString();
};

const openCreateDialog = () => {
    newTask.value = { title: '', description: '', priority: 0, project_id: '' };
    createError.value = '';
    showCreateDialog.value = true;
};

const handleCreateTask = async () => {
    if (!newTask.value.title.trim()) {
        createError.value = 'Title is required';
        return;
    }
    createLoading.value = true;
    try {
        await api.createTask({
            taskInputBody: {
                title: newTask.value.title.trim(),
                description: newTask.value.description,
                status: 'todo',
                priority: newTask.value.priority,
                project_id: newTask.value.project_id || undefined,
                due_date: selectedDateAsDue(),
            }
        });
        showCreateDialog.value = false;
        toast.showToast('success', 'Created', 'Task scheduled for ' + selectedDate.value.toLocaleDateString());
        await loadTasks();
    } catch (e) {
        console.error('Create task failed', e);
    } finally {
        createLoading.value = false;
    }
};

const openEditDialog = (task: any) => {
    editTask.value = { ...task };
    showEditDialog.value = true;
};

const handleUpdateTask = async () => {
    if (!editTask.value?.title?.trim()) return;
    editLoading.value = true;
    try {
        await api.updateTask({
            id: editTask.value.id,
            updateTaskRequest: {
                title: editTask.value.title.trim(),
                description: editTask.value.description || '',
                status: editTask.value.status,
                priority: editTask.value.priority,
                project_id: editTask.value.project_id || undefined,
                assigned_to: editTask.value.assigned_to || undefined,
                due_date: editTask.value.due_date || undefined,
            }
        });
        showEditDialog.value = false;
        toast.showToast('success', 'Updated', 'Task updated');
        await loadTasks();
    } catch (e) {
        console.error('Update task failed', e);
    } finally {
        editLoading.value = false;
    }
};

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
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950">
    <Toolbar class="p-4 shadow-sm">
      <template #start>
        <span class="text-2xl font-bold text-primary px-4">Calendar</span>
      </template>
      <template #end>
        <router-link to="/dashboard">
          <Button icon="pi pi-home" label="Dashboard" text />
        </router-link>
      </template>
    </Toolbar>

    <div class="p-6 flex flex-col items-center gap-6">
        <Calendar
            :attributes="attributes"
            expanded
            is-expanded
            class="max-w-4xl w-full border-none shadow-xl rounded-2xl p-4"
            :is-dark="true"
            @dayclick="onDayClick"
        />

        <div class="max-w-4xl w-full bg-white dark:bg-gray-900 rounded-2xl shadow-lg p-6">
            <div class="flex items-center justify-between mb-4">
                <h3 class="text-xl font-bold">Tasks for {{ selectedDate.toLocaleDateString() }}</h3>
                <Button icon="pi pi-plus" label="Add Task" size="small" @click="openCreateDialog" />
            </div>
            <div class="space-y-3">
                <div v-for="task in selectedTasks" :key="task.id"
                     class="p-3 border rounded-lg flex justify-between items-center dark:border-gray-800 cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
                     @click="openEditDialog(task)">
                    <div class="flex items-center gap-2">
                        <i class="pi pi-pencil text-xs text-gray-400"></i>
                        <span>{{ task.title }}</span>
                    </div>
                    <Tag :value="task.status.toUpperCase()" :severity="getStatusSeverity(task.status)" />
                </div>
                <div v-if="selectedTasks.length === 0" class="text-center py-4 text-gray-500 italic">
                    No tasks for this day. Click "Add Task" to schedule one.
                </div>
            </div>
        </div>
    </div>

    <!-- Create Task Dialog -->
    <Dialog v-model:visible="showCreateDialog" :header="'New Task — ' + selectedDate.toLocaleDateString()" :style="{ width: '450px' }" :breakpoints="{ '640px': '92vw' }" modal>
        <div class="flex flex-col gap-4">
            <div class="flex flex-col gap-2">
                <label for="newTitle" class="font-bold">Title</label>
                <InputText id="newTitle" v-model="newTask.title" :invalid="!!createError" placeholder="What needs to be done?" />
                <small v-if="createError" class="text-red-500">{{ createError }}</small>
            </div>
            <div class="flex flex-col gap-2">
                <label for="newProject" class="font-bold">Project (Optional)</label>
                <Select id="newProject" v-model="newTask.project_id" :options="projects" optionLabel="name" optionValue="id" placeholder="No project" showClear />
            </div>
            <div class="flex flex-col gap-2">
                <label class="font-bold">Priority</label>
                <SelectButton v-model="newTask.priority" :options="[{ label: 'Normal', value: 0 }, { label: 'High', value: 1 }]" optionLabel="label" optionValue="value" />
            </div>
            <div class="flex flex-col gap-2">
                <label for="newDesc" class="font-bold">Description</label>
                <Textarea id="newDesc" v-model="newTask.description" rows="3" />
            </div>
            <div class="flex justify-end gap-2 mt-4">
                <Button label="Cancel" severity="secondary" text @click="showCreateDialog = false" />
                <Button label="Create Task" :loading="createLoading" @click="handleCreateTask" />
            </div>
        </div>
    </Dialog>

    <!-- Edit Task Dialog -->
    <Dialog v-model:visible="showEditDialog" header="Edit Task" :style="{ width: '450px' }" :breakpoints="{ '640px': '92vw' }" modal>
        <div v-if="editTask" class="flex flex-col gap-4">
            <div class="flex flex-col gap-2">
                <label for="editTitle" class="font-bold">Title</label>
                <InputText id="editTitle" v-model="editTask.title" :invalid="!editTask.title?.trim()" />
                <small v-if="!editTask.title?.trim()" class="text-red-500">Title is required</small>
            </div>
            <div class="flex flex-col gap-2">
                <label class="font-bold">Status</label>
                <SelectButton v-model="editTask.status" :options="['todo', 'doing', 'done']" class="uppercase" />
            </div>
            <div class="flex flex-col gap-2">
                <label class="font-bold">Priority</label>
                <SelectButton v-model="editTask.priority" :options="[{ label: 'Normal', value: 0 }, { label: 'High', value: 1 }]" optionLabel="label" optionValue="value" />
            </div>
            <div class="flex flex-col gap-2">
                <label for="editDesc" class="font-bold">Description</label>
                <Textarea id="editDesc" v-model="editTask.description" rows="3" />
            </div>
            <div class="flex justify-end gap-2 mt-4">
                <Button label="Cancel" severity="secondary" text @click="showEditDialog = false" />
                <Button label="Save Changes" :loading="editLoading" :disabled="!editTask.title?.trim()" @click="handleUpdateTask" />
            </div>
        </div>
    </Dialog>
  </div>
</template>
