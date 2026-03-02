<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import InputText from 'primevue/inputtext';
import Password from 'primevue/password';
import Button from 'primevue/button';
import Message from 'primevue/message';
import Card from 'primevue/card';

const router = useRouter();
const auth = useAuthStore();

const email = ref('test@example.com');
const password = ref('password123');
const error = ref('');
const loading = ref(false);

const handleLogin = async () => {
  loading.value = true;
  error.value = '';
  
  try {
    const response = await fetch('/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ body: { email: email.value, password: password.value } }),
    });

    if (!response.ok) {
      throw new Error('Invalid credentials');
    }

    const data = await response.json();
    auth.setAuth(data.token, data.user);
    router.push('/dashboard');
  } catch (err: any) {
    error.value = err.message;
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
        <div class="text-center text-gray-500">Login to manage your tasks</div>
      </template>
      <template #content>
        <form @submit.prevent="handleLogin" class="flex flex-col gap-6 mt-4">
          <div class="flex flex-col gap-2">
            <label for="email" class="font-semibold">Email</label>
            <InputText id="email" v-model="email" type="email" placeholder="email@example.com" required />
          </div>
          
          <div class="flex flex-col gap-2">
            <label for="password" class="font-semibold">Password</label>
            <Password id="password" v-model="password" :feedback="false" toggleMask placeholder="Your password" required />
          </div>

          <Message v-if="error" severity="error">{{ error }}</Message>

          <Button type="submit" label="Sign In" :loading="loading" class="w-full mt-2" />
        </form>
      </template>
    </Card>
  </div>
</template>
