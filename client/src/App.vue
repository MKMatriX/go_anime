<template>
  <div class="min-h-screen bg-gray-100">
    <header class="bg-white shadow">
      <div class="max-w-4xl mx-auto px-4 py-4 flex items-center justify-between">
        <h1 class="text-xl font-semibold text-gray-800">
          Vue Auth Demo
        </h1>
        <nav class="flex items-center gap-4">
          <RouterLink
            v-if="auth.isAuthenticated"
            class="text-sm text-blue-600 hover:underline"
            to="/anime"
          >
            Anime
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

    <main class="mx-auto px-4 py-10">
      <RouterView />
    </main>
  </div>
</template>

<script setup>
import { useUserStore } from './stores/auth'
import { useRouter, RouterLink, RouterView } from 'vue-router'

const auth = useUserStore()
const router = useRouter()

const handleLogout = () => {
  auth.logout()
  router.push({ name: 'Login' })
}
</script>