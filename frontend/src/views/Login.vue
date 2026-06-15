<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import axiosInstance from '../api/axios';
import { DefaultApi } from '../api';
import { isValidEmail } from '../lib/validation';
import InputText from 'primevue/inputtext';
import Password from 'primevue/password';
import Button from 'primevue/button';
import Message from 'primevue/message';
import Card from 'primevue/card';

const router = useRouter();
const auth = useAuthStore();
const api = new DefaultApi(undefined, '/api', axiosInstance);

const email = ref('test@example.com');
const password = ref('password123');
const error = ref('');
const loading = ref(false);
const fieldErrors = ref<{ email?: string; password?: string }>({});

const validate = () => {
  fieldErrors.value = {};
  if (!email.value) {
    fieldErrors.value.email = 'Email is required';
  } else if (!isValidEmail(email.value)) {
    fieldErrors.value.email = 'Enter a valid email address';
  }
  if (!password.value) {
    fieldErrors.value.password = 'Password is required';
  }
  return Object.keys(fieldErrors.value).length === 0;
};

const handleLogin = async () => {
  error.value = '';
  if (!validate()) return;
  loading.value = true;

  try {
    const { data } = await api.login({
      loginInputBody: { email: email.value, password: password.value },
    });
    auth.setAuth(data.token, data.user);
    router.push('/dashboard');
  } catch (err: any) {
    error.value = 'Invalid credentials';
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
            <InputText id="email" v-model="email" type="email" placeholder="email@example.com" :invalid="!!fieldErrors.email" required />
            <small v-if="fieldErrors.email" class="text-red-500">{{ fieldErrors.email }}</small>
          </div>

          <div class="flex flex-col gap-2">
            <label for="password" class="font-semibold">Password</label>
            <Password id="password" v-model="password" :feedback="false" toggleMask placeholder="Your password" :invalid="!!fieldErrors.password" required />
            <small v-if="fieldErrors.password" class="text-red-500">{{ fieldErrors.password }}</small>
          </div>

          <Message v-if="error" severity="error">{{ error }}</Message>

          <Button type="submit" label="Login" :loading="loading" class="w-full mt-2" />
        </form>
      </template>
    </Card>
  </div>
</template>
