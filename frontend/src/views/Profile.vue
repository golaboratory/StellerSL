<script setup lang="ts">
import { ref } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useToastStore } from '../stores/toast';
import { DefaultApi } from '../api';
import axiosInstance from '../api/axios';
import Card from 'primevue/card';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import Toolbar from 'primevue/toolbar';
import Avatar from 'primevue/avatar';
import FileUpload from 'primevue/fileupload';
import Password from 'primevue/password';
import { passwordStrengthError } from '../lib/validation';

const auth = useAuthStore();
const toast = useToastStore();
const api = new DefaultApi(undefined, '/api', axiosInstance);

const name = ref(auth.user?.name || '');
const avatarUrl = ref(auth.user?.avatar_url || '');
const loading = ref(false);
const uploadLoading = ref(false);

// Password change form (7-FE-4)
const currentPassword = ref('');
const newPassword = ref('');
const confirmPassword = ref('');
const passwordLoading = ref(false);
const passwordError = ref('');

const handleChangePassword = async () => {
    passwordError.value = '';
    if (!currentPassword.value) {
        passwordError.value = 'Current password is required';
        return;
    }
    const strengthError = passwordStrengthError(newPassword.value);
    if (strengthError) {
        passwordError.value = strengthError;
        return;
    }
    if (newPassword.value !== confirmPassword.value) {
        passwordError.value = 'New password and confirmation do not match';
        return;
    }
    if (newPassword.value === currentPassword.value) {
        passwordError.value = 'New password must differ from the current password';
        return;
    }

    passwordLoading.value = true;
    try {
        await api.changePassword({
            changePasswordInputBody: {
                current_password: currentPassword.value,
                new_password: newPassword.value,
            }
        });
        currentPassword.value = '';
        newPassword.value = '';
        confirmPassword.value = '';
        toast.showToast('success', 'Success', 'Password changed successfully');
    } catch (err) {
        console.error('Failed to change password', err);
    } finally {
        passwordLoading.value = false;
    }
};

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
            auth.user.avatar_url = avatarUrl.value;
        }
        
        toast.showToast('success', 'Success', 'Profile updated successfully');
    } catch (err) {
        console.error('Failed to update profile', err);
    } finally {
        loading.value = false;
    }
};

const onUpload = async (event: any) => {
    const file = event.files[0];
    if (!file) return;

    uploadLoading.value = true;
    const formData = new FormData();
    formData.append('file', file);

    try {
        // Using axiosInstance directly for multipart upload
        const response = await axiosInstance.post('/auth/avatar', formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        });
        
        const newUrl = response.data.url;
        avatarUrl.value = newUrl;
        
        if (auth.user) {
            auth.user.avatar_url = newUrl;
        }
        
        toast.showToast('success', 'Success', 'Avatar uploaded successfully');
    } catch (err) {
        console.error('Upload failed', err);
        toast.showToast('error', 'Upload Error', 'Failed to upload image');
    } finally {
        uploadLoading.value = false;
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
            <div class="relative group cursor-pointer mb-4">
                <Avatar :image="avatarUrl" :label="!avatarUrl ? name.charAt(0).toUpperCase() : undefined" class="w-32 h-32 text-4xl" shape="circle" />
                <div class="absolute inset-0 flex items-center justify-center bg-black/40 rounded-full opacity-0 group-hover:opacity-100 transition-opacity">
                    <i class="pi pi-camera text-white text-2xl"></i>
                </div>
                <FileUpload mode="basic" name="file" accept="image/*" :maxFileSize="1000000" @select="onUpload" class="absolute inset-0 opacity-0 w-full h-full cursor-pointer" :auto="true" />
            </div>
            
            <span v-if="uploadLoading" class="text-sm text-primary mb-2"><i class="pi pi-spin pi-spinner mr-2"></i>Uploading...</span>
            <span class="text-xl font-bold">{{ name }}</span>
            <span class="text-gray-500 text-sm">{{ auth.user?.email }}</span>
          </div>

          <form @submit.prevent="handleSave" class="space-y-6">
            <div class="flex flex-col gap-2">
              <label for="name" class="font-medium text-gray-700 dark:text-gray-300">Display Name</label>
              <InputText id="name" v-model="name" class="w-full" placeholder="Your name" />
            </div>

            <div class="flex flex-col gap-2">
              <label for="avatar" class="font-medium text-gray-700 dark:text-gray-300">Avatar URL (Optional)</label>
              <InputText id="avatar" v-model="avatarUrl" class="w-full" placeholder="https://example.com/avatar.png" />
              <small class="text-gray-500">You can also click on the avatar above to upload a file.</small>
            </div>

            <div class="pt-4">
              <Button type="submit" label="Save Changes" icon="pi pi-save" class="w-full" :loading="loading" />
            </div>
          </form>
        </template>
      </Card>

      <Card class="mt-6 shadow-sm">
        <template #title>Change Password</template>
        <template #content>
          <form @submit.prevent="handleChangePassword" class="space-y-4">
            <div class="flex flex-col gap-2">
              <label for="currentPassword" class="font-medium text-gray-700 dark:text-gray-300">Current Password</label>
              <Password inputId="currentPassword" v-model="currentPassword" :feedback="false" toggleMask fluid />
            </div>
            <div class="flex flex-col gap-2">
              <label for="newPassword" class="font-medium text-gray-700 dark:text-gray-300">New Password</label>
              <Password inputId="newPassword" v-model="newPassword" toggleMask fluid promptLabel="Min 8 chars, letters and numbers" />
            </div>
            <div class="flex flex-col gap-2">
              <label for="confirmPassword" class="font-medium text-gray-700 dark:text-gray-300">Confirm New Password</label>
              <Password inputId="confirmPassword" v-model="confirmPassword" :feedback="false" toggleMask fluid
                        :invalid="!!confirmPassword && newPassword !== confirmPassword" />
            </div>
            <small v-if="passwordError" class="text-red-500 block">{{ passwordError }}</small>
            <Button type="submit" label="Change Password" icon="pi pi-lock" severity="secondary" class="w-full" :loading="passwordLoading" />
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
                    <span class="text-gray-500">Email:</span>
                    <span class="font-mono text-xs">{{ auth.user?.email }}</span>
                </div>
            </div>
        </template>
      </Card>
    </div>
  </div>
</template>
