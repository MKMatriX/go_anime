// src/features/anime/router/anime.routes.ts
import { RouteRecordRaw } from 'vue-router'

export const animeRoutes: RouteRecordRaw[] = [
	{
		path: '/anime',
		meta: { requiresAuth: true },           // можно группу мета-данных
		children: [
			{
				path: '',
				name: 'anime.list',
				component: () => import('../views/anime/List.vue'),
			},
			{
				path: 'create',
				name: 'anime.create',
				component: () => import('../views/anime/Create.vue'),
			},
			{
				path: ':id(\\d+)/edit',
				name: 'anime.edit',
				component: () => import('../views/anime/Edit.vue'),
				props: true,
			},
			{
				path: ':id',
				name: 'anime.detail',
				component: () => import('../views/anime/View.detail.vue'),
				props: true,
			}
		]
	}
]