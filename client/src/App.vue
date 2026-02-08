<template>
  <div class="min-h-screen bg-gray-100">
    <header class="bg-white shadow">
      <div class="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
        <h1 class="text-xl font-semibold text-gray-800">
          <a href="/" class="hover:shadow">
            MKMatriX Pet Project
          </a>
        </h1>
        <nav class="flex items-center gap-4">
          <RouterLink
            v-if="auth.isAuthenticated"
            class="text-sm text-blue-600 hover:underline"
            to="/anime"
          >
            Anime
          </RouterLink>
          <RouterLink
            v-if="inAnime"
            class="text-sm text-blue-600 hover:underline"
            to="/anime/create"
          >
            Добавить
          </RouterLink>
        </nav>
        <nav class="flex items-center gap-4">
          <RouterLink
            v-if="auth.isAuthenticated"
            class="text-sm text-blue-600 hover:underline"
            to="/profile"
          >
            Profile
          </RouterLink>
          <RouterLink
            v-if="!auth.isAuthenticated"
            class="text-sm text-blue-600 hover:underline"
            to="/login"
          >
            Login
          </RouterLink>
          <RouterLink
            v-if="!auth.isAuthenticated"
            class="text-sm text-blue-600 hover:underline"
            to="/register"
          >
            Register
          </RouterLink>
          <button
            v-if="auth.isAuthenticated"
            @click="handleLogout"
            class="text-sm text-red-600 hover:underline"
          >
            Logout
          </button>
        </nav>
      </div>
    </header>

    <main class="max-w-7xl mx-auto px-4 py-10">
      <RouterView />
    </main>
  </div>
</template>

<script setup>
import { useUserStore } from './stores/auth'
import { useRouter, RouterLink, RouterView } from 'vue-router'
import { onMounted, ref } from "vue";

const auth = useUserStore()
const route = useRouter()
const inAnime = ref(false)

onMounted(() => {
  let matched = route.currentRoute.value.matched[0] || {}
  inAnime.value = matched.path === "/anime"
})

route.afterEach((to, from) => {
  let matched = route.currentRoute.value.matched[0] || {}
  inAnime.value = matched.path === "/anime"
});


const handleLogout = () => {
  auth.logout()
  route.push({ name: 'Login' })
}
</script>