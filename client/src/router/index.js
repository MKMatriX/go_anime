// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from "../stores/auth";

// Lazy-load pages
const Login = () => import('../views/Login.vue')
const Register = () => import('../views/Register.vue')
const Profile = () => import('../views/Profile.vue')
// Пока можно оставить пустой массив маршрутов
const routes = [
		{
			path: '/login',
			name: 'Login',
			component: Login,
			meta: { guestOnly: true },
		},
		{
			path: '/register',
			name: 'Register',
			component: Register,
			meta: { guestOnly: true },
		},
		{
			path: '/profile',
			name: 'Profile',
			component: Profile,
			meta: { requiresAuth: true },
		},
		{
			path: '/',
			redirect: '/login', // default route
		},
]


const router = createRouter({
	history: createWebHistory(),  // ← важно для Vite
	routes
})

router.beforeEach((to, from, next) => {
	const auth = useAuthStore()

	if (to.meta.requiresAuth && !auth.isAuthenticated) {
		return next({ name: 'Login' })
	}

	if (to.meta.guestOnly && auth.isAuthenticated) {
		return next({ name: 'Profile' })
	}

	next()
})

export default router