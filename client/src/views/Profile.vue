<template>
  <div class="bg-white shadow rounded-lg p-6">
    <h2 class="text-2xl font-semibold mb-4">Profile</h2>

    <div v-if="auth.loading" class="text-gray-600">
      Loading profile...
    </div>

    <div v-else-if="auth.user" class="space-y-2 text-sm">
      <p>
        <span class="font-medium">ID: </span>
        <span>{{ auth.user.id }}</span>
      </p>
      <p>
        <span class="font-medium">Login: </span>
        <span>{{ auth.user.login }}</span>
      </p>
    </div>

    <div v-else class="text-gray-600">
      No profile data loaded.
    </div>

    <p v-if="auth.error" class="text-sm text-red-600 mt-3">
      {{ auth.error }}
    </p>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useUserStore } from '../stores/auth'

const auth = useUserStore()

onMounted(async () => {
  if (auth.isAuthenticated && !auth.user) {
    try {
      await auth.fetchProfile()
    } catch {
      // auth.error already contains the message
    }
  }
})
</script>