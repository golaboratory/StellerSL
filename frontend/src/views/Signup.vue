<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import axiosInstance from '../api/axios';
import { DefaultApi } from '../api';
import InputText from 'primevue/inputtext';
import Card from 'primevue/card';
import Password from 'primevue/password';
import Message from 'primevue/message';
import Button from 'primevue/button';

const router = useRouter();
const route = useRoute();
const api = new DefaultApi(undefined, '/api', axiosInstance);

const name = ref('');
const email = ref('');
const password = ref('');
const inviteTeamID = ref('');
const error = ref('');
const loading = ref(false);

onMounted(() => {
    if (route.query.team) {
        inviteTeamID.value = route.query.team as string;
    }
});

const handleSignup = async () => {
  loading.value = true;
  error.value = '';
  
  try {
    await api.register({
      registerInputBody: {
        tenant_id: "00000000-0000-0000-0000-000000000001", // Default tenant for POC
        email: email.value,
        password: password.value,
        name: name.value,
        invite_team_id: inviteTeamID.value || undefined,
      },
    });
    router.push('/login');
  } catch (err: any) {
    error.value = 'Failed to register. Please try again.';
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <div class="flex items-center justify-center min-h-screen bg-gray-100 dark:bg-gray-900 px-4">
    <Card class="w-full max-w-md shadow-lg">
      <template #title>
        <div class="text-center text-3xl font-bold text-primary">StellerSL</div>
      </template>
      <template #subtitle>
        <div class="text-center text-gray-500">Create an account</div>
      </template>
      <template #content>
        <form @submit.prevent="handleSignup" class="flex flex-col gap-6 mt-4">
          <div class="flex flex-col gap-2">
            <label for="name" class="font-semibold">Full Name</label>
            <InputText id="name" v-model="name" placeholder="John Doe" required />
          </div>

          <div class="flex flex-col gap-2">
            <label for="email" class="font-semibold">Email</label>
            <InputText id="email" v-model="email" type="email" placeholder="email@example.com" required />
          </div>
          
          <div class="flex flex-col gap-2">
            <label for="password" class="font-semibold">Password</label>
            <Password id="password" v-model="password" :feedback="false" toggleMask placeholder="Min 8 characters" required />
          </div>

          <Message v-if="error" severity="error">{{ error }}</Message>

          <Button type="submit" label="Sign Up" :loading="loading" class="w-full mt-2" />
          
          <div class="text-center text-sm mt-2">
            Already have an account? 
            <router-link to="/login" class="text-primary font-bold hover:underline">Login</router-link>
          </div>
        </form>
      </template>
    </Card>
  </div>
</template>
