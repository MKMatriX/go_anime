// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from "../stores/auth";
import { animeRoutes } from "./anime";

const routes = [
		{
			path: '/login',
			name: 'Login',
			component: () => import('../views/user/Login.vue'),
			meta: { guestOnly: true },
		},
		{
			path: '/register',
			name: 'Register',
			component: () => import('../views/user/Register.vue'),
			meta: { guestOnly: true },
		},
		{
			path: '/profile',
			name: 'Profile',
			component: () => import('../views/Profile.vue'),
			meta: { requiresAuth: true },
		},
		{
			path: '/',
			redirect: '/login', // default route
		},
		...animeRoutes
]


const router = createRouter({
	history: createWebHistory(),  // ← важно для Vite
	routes
})

router.beforeEach((to, from, next) => {
	const auth = useUserStore()

	if (to.meta.requiresAuth && !auth.isAuthenticated) {
		return next({ name: 'Login' })
	}

	if (to.meta.guestOnly && auth.isAuthenticated) {
		return next({ name: 'Profile' })
	}

	next()
})

export default router