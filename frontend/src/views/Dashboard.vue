<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { DefaultApi } from '../api';
import type { DashboardOutputBody } from '../api';
import axiosInstance from '../api/axios';
import Tag from 'primevue/tag';
import { useToast } from 'primevue/usetoast';
import Card from 'primevue/card';
import Toolbar from 'primevue/toolbar';
import Button from 'primevue/button';
import ProgressBar from 'primevue/progressbar';
import NotificationBell from '../components/NotificationBell.vue';
import { badgeIconUrl } from '../lib/badgeIcons';

const auth = useAuthStore();
const router = useRouter();
const toast = useToast();
const api = new DefaultApi(undefined, '/api', axiosInstance);


const stats = ref<DashboardOutputBody>({ total_tasks: 0, pending_tasks: 0, completed_tasks: 0, daily_activity: [], recent_activity: [] });
const growth = ref<any>({ level: 1, exp: 0, character_type: 'default' });
const badges = ref<any[]>([]);
const loading = ref(true);

onMounted(async () => {
  if (!auth.isAuthenticated) {
    router.push('/login');
    return;
  }

  try {
    const [statsRes, growthRes, badgesRes] = await Promise.all([
        api.getDashboardStats(),
        api.getUserGrowth(),
        api.listBadges()
    ]);
    stats.value = statsRes.data;
    growth.value = growthRes.data;

    const newBadges = badgesRes.data.items || [];
    const oldBadgeCount = parseInt(localStorage.getItem('badge_count') || '0');

    if (newBadges.length > oldBadgeCount) {
        toast.add({
            severity: 'success',
            summary: 'Achievement Unlocked!',
            detail: `You've earned ${newBadges.length - oldBadgeCount} new badge(s)!`,
            life: 5000
        });
    }

    badges.value = newBadges;
    localStorage.setItem('badge_count', newBadges.length.toString());
  } catch (err) {
    console.error('Failed to fetch dashboard data', err);
  } finally {
    loading.value = false;
  }
});

const handleLogout = () => {
  auth.clearAuth();
  router.push('/login');
};

const toggleDarkMode = () => {
    document.documentElement.classList.toggle('p-dark');
};

const getExpProgress = () => {
    return (growth.value.exp % 100);
};

// Character evolution stages (バッジの仕様.md):
// Lv.1-4 egg / Lv.5-9 chick / Lv.10-19 chicken / Lv.20+ phoenix
type CharacterStyle = { icon: string; label: string; color: string; bg: string };
const characterConfig: { egg: CharacterStyle; chick: CharacterStyle; chicken: CharacterStyle; phoenix: CharacterStyle } = {
    egg:     { icon: 'pi-circle',  label: 'たまご (Egg)',           color: 'text-primary-500', bg: 'bg-primary-100' },
    chick:   { icon: 'pi-twitter', label: 'ひよこ (Chick)',         color: 'text-orange-500',  bg: 'bg-orange-100' },
    chicken: { icon: 'pi-heart',   label: 'にわとり (Chicken)',     color: 'text-red-500',     bg: 'bg-red-100' },
    phoenix: { icon: 'pi-bolt',    label: 'フェニックス (Phoenix)', color: 'text-yellow-500',  bg: 'bg-yellow-100' },
};

const getAvatarConfig = (): CharacterStyle => {
    const byType = (characterConfig as Record<string, CharacterStyle | undefined>)[growth.value.character_type];
    if (byType) return byType;
    // Fallback for accounts whose character_type predates auto-evolution
    const level = growth.value.level;
    if (level >= 20) return characterConfig.phoenix;
    if (level >= 10) return characterConfig.chicken;
    if (level >= 5) return characterConfig.chick;
    return characterConfig.egg;
};

const formatAction = (action: string) => {
    return action.replace('_', ' ').toUpperCase();
};
</script>


<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950 p-6">

    <Toolbar class="mb-8 p-4 rounded-xl shadow-sm">
      <template #start>
        <span class="text-2xl font-bold text-primary px-4">StellerSL Dashboard</span>
      </template>
      <template #end>
        <div class="nav-compact flex items-center justify-end flex-wrap gap-1 md:gap-3">
          <router-link to="/teams">
            <Button icon="pi pi-users" label="Teams" text />
          </router-link>
          <router-link to="/tasks">
            <Button icon="pi pi-check-square" label="Tasks" text />
          </router-link>
          <router-link to="/gantt">
            <Button icon="pi pi-chart-bar" label="Gantt" text />
          </router-link>
          <router-link to="/calendar">
            <Button icon="pi pi-calendar" label="Calendar" text />
          </router-link>
          <router-link to="/profile">
            <Button icon="pi pi-user" label="Profile" text />
          </router-link>
          <NotificationBell />
          <Button icon="pi pi-moon" text @click="toggleDarkMode" />
          <span class="font-medium hidden md:inline">Welcome, {{ auth.user?.name }}</span>
          <Button icon="pi pi-sign-out" label="Logout" severity="secondary" text @click="handleLogout" />
        </div>
      </template>
    </Toolbar>

    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-3 gap-6 animate-pulse">
        <div v-for="i in 3" :key="i" class="h-32 bg-gray-200 dark:bg-gray-800 rounded-xl"></div>
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-4 gap-6 mb-8">
      <!-- Growth Card -->
      <Card class="lg:col-span-1 border-t-4 border-primary-500 shadow-sm">
        <template #title>Your Character</template>
        <template #content>
          <div class="flex flex-col items-center gap-4">
            <div :class="['w-24 h-24 rounded-full flex items-center justify-center border-4', getAvatarConfig().bg, getAvatarConfig().color.replace('text', 'border')]">
                <i :class="['pi text-4xl', getAvatarConfig().icon, getAvatarConfig().color]"></i>
            </div>
            <div class="text-center">
                <div class="text-xl font-bold uppercase tracking-wider text-primary">Level {{ growth.level }}</div>
                <div class="text-sm text-gray-500 mt-1">Character: {{ getAvatarConfig().label }}</div>
            </div>
            <div class="w-full">
                <div class="flex justify-between text-xs mb-1">
                    <span>EXP: {{ growth.exp }}</span>
                    <span>Next: {{ (Math.floor(growth.exp / 100) + 1) * 100 }}</span>
                </div>
                <ProgressBar :value="getExpProgress()" :show-value="false" style="height: 8px" />
            </div>
          </div>
        </template>
      </Card>

      <!-- Task Stats -->
      <div class="lg:col-span-3 grid grid-cols-1 md:grid-cols-3 gap-6">
        <Card class="border-l-4 border-blue-500 shadow-sm">
            <template #title>Total Tasks</template>
            <template #content>
            <div class="text-4xl font-bold">{{ stats.total_tasks }}</div>
            </template>
        </Card>
        
        <Card class="border-l-4 border-orange-500 shadow-sm">
            <template #title>Pending</template>
            <template #content>
            <div class="text-4xl font-bold text-orange-500">{{ stats.pending_tasks }}</div>
            </template>
        </Card>

        <Card class="border-l-4 border-green-500 shadow-sm">
            <template #title>Completed</template>
            <template #content>
            <div class="text-4xl font-bold text-green-500">{{ stats.completed_tasks }}</div>
            </template>
        </Card>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        <!-- Daily Activity -->
        <Card class="shadow-sm">
            <template #title>Daily Activity</template>
            <template #content>
                <div class="h-64 flex items-end gap-2 border-b-2 border-gray-100 dark:border-gray-800 pb-2">
                    <div v-for="day in stats.daily_activity" :key="day.date" class="flex-1 flex flex-col gap-1 items-center group">
                        <div class="w-full bg-blue-400 rounded-t-sm" :style="{ height: Math.min((day.created ?? 0) * 10, 150) + 'px' }"></div>
                        <div class="w-full bg-green-400 rounded-t-sm" :style="{ height: Math.min((day.completed ?? 0) * 10, 150) + 'px' }"></div>
                        <span class="text-xs text-gray-400 mt-2">{{ (day.date ?? '').split('-').slice(1).join('/') }}</span>
                    </div>
                </div>
                <div class="flex gap-4 mt-6 text-sm text-gray-500 justify-center">
                    <div class="flex items-center gap-2"><div class="w-3 h-3 bg-blue-400"></div> Created</div>
                    <div class="flex items-center gap-2"><div class="w-3 h-3 bg-green-400"></div> Completed</div>
                </div>
            </template>
        </Card>

        <!-- Recent Activity -->
        <Card class="shadow-sm">
            <template #title>Recent Activity</template>
            <template #content>
                <div class="space-y-4">
                    <div v-for="log in stats.recent_activity" :key="log.id" class="flex items-center justify-between p-3 border-b border-gray-100 dark:border-gray-800 last:border-0">
                        <div class="flex flex-col">
                            <span class="text-sm font-medium">{{ log.task_title }}</span>
                            <span class="text-xs text-gray-500">{{ log.date }}</span>
                        </div>
                        <Tag :value="formatAction(log.action)" :severity="log.action.includes('completed') ? 'success' : 'info'" />
                    </div>
                    <div v-if="!stats.recent_activity || stats.recent_activity.length === 0" class="text-center py-8 text-gray-400 italic">
                        No recent activity.
                    </div>
                </div>
            </template>
        </Card>
    </div>

    <div class="grid grid-cols-1 gap-6">
        <!-- Badges -->
        <Card class="shadow-sm">
            <template #title>Achievements</template>
            <template #content>
                <div v-if="badges.length === 0" class="flex flex-col items-center justify-center h-32 text-gray-400 italic">
                    <i class="pi pi-lock text-4xl mb-4"></i>
                    <span>Complete tasks to earn badges!</span>
                </div>
                <div v-else class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
                    <div v-for="badge in badges" :key="badge.id" class="flex flex-col items-center gap-2 p-4 border rounded-xl bg-gray-50 dark:bg-gray-900 border-gray-200 dark:border-gray-800 transition-transform hover:scale-105">
                        <img v-if="badgeIconUrl(badge.requirement_type)" :src="badgeIconUrl(badge.requirement_type)!" :alt="badge.name" class="badge-icon w-14 h-14" />
                        <i v-else :class="'pi ' + badge.icon_slug" class="text-3xl text-yellow-500"></i>
                        <span class="text-xs font-bold text-center">{{ badge.name }}</span>
                        <span class="text-[10px] text-gray-500 text-center">{{ badge.description }}</span>
                    </div>
                </div>
            </template>
        </Card>
    </div>
  </div>
</template>
