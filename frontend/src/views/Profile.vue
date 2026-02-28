<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useToastStore } from '../stores/toast';
import { DefaultApi } from '../api';
import axiosInstance from '../api/axios';
import Card from 'primevue/card';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import Toolbar from 'primevue/toolbar';
import Avatar from 'primevue/avatar';

const auth = useAuthStore();
const toast = useToastStore();
const api = new DefaultApi(undefined, '/api', axiosInstance);

const name = ref(auth.user?.name || '');
const avatarUrl = ref(''); // We don't have it in the user object yet, need to fetch or add to auth store
const loading = ref(false);

const handleSave = async () => {
    if (!name.value) return;
    loading.value = true;
    try {
        await api.updateProfile({ 
            updateProfileInputBody: { name: name.value, avatar_url: avatarUrl.value } 
        });
        
        // Update local store
        if (auth.user) {
            auth.user.name = name.value;
            // If we had avatarUrl in store, update it too
        }
        
        toast.showToast('success', 'Success', 'Profile updated successfully');
    } catch (err) {
        toast.showToast('error', 'Error', 'Failed to update profile');
    } finally {
        loading.value = false;
    }
};
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950 p-6">
    <Toolbar class="mb-8 p-4 rounded-xl shadow-sm">
      <template #start>
        <span class="text-2xl font-bold text-primary px-4">User Profile</span>
      </template>
      <template #end>
        <router-link to="/dashboard">
          <Button icon="pi pi-home" label="Dashboard" text />
        </router-link>
      </template>
    </Toolbar>

    <div class="max-w-2xl mx-auto">
      <Card class="shadow-sm">
        <template #content>
          <div class="flex flex-col items-center mb-8">
            <Avatar :image="avatarUrl || 'https://www.gravatar.com/avatar/00000000000000000000000000000000?d=mp&f=y'" class="w-32 h-32 text-4xl mb-4" shape="circle" />
            <span class="text-xl font-bold">{{ name }}</span>
            <span class="text-gray-500 text-sm">{{ auth.user?.email }}</span>
          </div>

          <form @submit.prevent="handleSave" class="space-y-6">
            <div class="flex flex-col gap-2">
              <label for="name" class="font-medium text-gray-700 dark:text-gray-300">Display Name</label>
              <InputText id="name" v-model="name" class="w-full" placeholder="Your name" />
            </div>

            <div class="flex flex-col gap-2">
              <label for="avatar" class="font-medium text-gray-700 dark:text-gray-300">Avatar URL</label>
              <InputText id="avatar" v-model="avatarUrl" class="w-full" placeholder="https://example.com/avatar.png" />
            </div>

            <div class="pt-4">
              <Button type="submit" label="Save Changes" icon="pi pi-save" class="w-full" :loading="loading" />
            </div>
          </form>
        </template>
      </Card>

      <Card class="mt-6 shadow-sm border-l-4 border-blue-500">
        <template #title>Account Info</template>
        <template #content>
            <div class="space-y-2">
                <div class="flex justify-between">
                    <span class="text-gray-500">User ID:</span>
                    <span class="font-mono text-xs">{{ auth.user?.id }}</span>
                </div>
                <div class="flex justify-between">
                    <span class="text-gray-500">Tenant ID:</span>
                    <span class="font-mono text-xs">{{ auth.tenantID }}</span>
                </div>
            </div>
        </template>
      </Card>
    </div>
  </div>
</template>
