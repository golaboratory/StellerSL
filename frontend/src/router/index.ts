import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import Login from '../views/Login.vue';
import Signup from '../views/Signup.vue';
import Dashboard from '../views/Dashboard.vue';
import Projects from '../views/Projects.vue';
import ProjectDetail from '../views/ProjectDetail.vue';
import Tasks from '../views/Tasks.vue';
import Teams from '../views/Teams.vue';
import TeamDetail from '../views/TeamDetail.vue';
import Gantt from '../views/Gantt.vue';
import CalendarView from '../views/CalendarView.vue';

const routes = [
  { path: '/', redirect: '/dashboard' },
  { path: '/login', name: 'Login', component: Login },
  { path: '/signup', name: 'Signup', component: Signup },
  { 
    path: '/dashboard', 
    name: 'Dashboard', 
    component: Dashboard,
    meta: { requiresAuth: true }
  },
  { 
    path: '/projects', 
    name: 'Projects', 
    component: Projects,
    meta: { requiresAuth: true }
  },
  { 
    path: '/projects/:id', 
    name: 'ProjectDetail', 
    component: ProjectDetail,
    meta: { requiresAuth: true }
  },
  { 
    path: '/tasks', 
    name: 'Tasks', 
    component: Tasks,
    meta: { requiresAuth: true }
  },
  { 
    path: '/teams', 
    name: 'Teams', 
    component: Teams,
    meta: { requiresAuth: true }
  },
  { 
    path: '/teams/:id', 
    name: 'TeamDetail', 
    component: TeamDetail,
    meta: { requiresAuth: true }
  },
  { 
    path: '/gantt', 
    name: 'Gantt', 
    component: Gantt,
    meta: { requiresAuth: true }
  },
  { 
    path: '/calendar', 
    name: 'Calendar', 
    component: CalendarView,
    meta: { requiresAuth: true }
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

router.beforeEach((to, from, next) => {
  const auth = useAuthStore();
  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    next('/login');
  } else if (to.name === 'Login' && auth.isAuthenticated) {
    next('/dashboard');
  } else {
    next();
  }
});

export default router;
