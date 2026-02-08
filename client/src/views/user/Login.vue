<template>
  <div class="bg-white shadow rounded-lg p-6">
    <h2 class="text-2xl font-semibold mb-4">Login</h2>

    <form @submit.prevent="onSubmit" class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">
          Login
        </label>
        <input
          v-model="form.login"
          type="text"
          class="w-full border rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <p v-if="errors.login" class="text-xs text-red-600 mt-1">
          {{ errors.login }}
        </p>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">
          Password
        </label>
        <input
          v-model="form.password"
          type="password"
          class="w-full border rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <p v-if="errors.password" class="text-xs text-red-600 mt-1">
          {{ errors.password }}
        </p>
      </div>

      <p v-if="auth.error" class="text-sm text-red-600">
        {{ auth.error }}
      </p>

      <button
        type="submit"
        :disabled="auth.loading"
        class="w-full bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium py-2 px-4 rounded-md disabled:opacity-60"
      >
        <span v-if="auth.loading">Logging in...</span>
        <span v-else>Login</span>
      </button>
    </form>

    <p class="mt-4 text-sm text-gray-600">
      Don’t have an account?
      <RouterLink to="/register" class="text-blue-600 hover:underline">
        Register
      </RouterLink>
    </p>
  </div>
</template>

<script setup>
import { reactive } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useUserStore } from '../../stores/auth'

const auth = useUserStore()
const router = useRouter()

const form = reactive({
  login: '',
  password: '',
})

const errors = reactive({
  login: '',
  password: '',
})

const validate = () => {
  errors.login = ''
  errors.password = ''

  if (!form.login.trim()) {
    errors.login = 'Login is required'
  }
  if (!form.password.trim()) {
    errors.password = 'Password is required'
  }

  return !errors.login && !errors.password
}

const onSubmit = async () => {
  if (!validate()) return
  try {
    await auth.login({
      login: form.login,
      password: form.password,
    })
    router.push({ name: 'Profile' })
  } catch {
    // error already stored in auth.error
  }
}
</script>