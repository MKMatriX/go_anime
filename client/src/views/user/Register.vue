<template>
  <div class="bg-white shadow rounded-lg p-6">
    <h2 class="text-2xl font-semibold mb-4">Register</h2>

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

      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">
          Confirm Password
        </label>
        <input
          v-model="form.confirmPassword"
          type="password"
          class="w-full border rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <p v-if="errors.confirmPassword" class="text-xs text-red-600 mt-1">
          {{ errors.confirmPassword }}
        </p>
      </div>

      <p v-if="auth.error" class="text-sm text-red-600">
        {{ auth.error }}
      </p>

      <button
        type="submit"
        :disabled="auth.loading"
        class="w-full bg-green-600 hover:bg-green-700 text-white text-sm font-medium py-2 px-4 rounded-md disabled:opacity-60"
      >
        <span v-if="auth.loading">Registering...</span>
        <span v-else>Register</span>
      </button>
    </form>

    <p class="mt-4 text-sm text-gray-600">
      Already have an account?
      <RouterLink to="/login" class="text-blue-600 hover:underline">
        Login
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
  confirmPassword: '',
})

const errors = reactive({
  login: '',
  password: '',
  confirmPassword: '',
})

const validate = () => {
  errors.login = ''
  errors.password = ''
  errors.confirmPassword = ''

  if (!form.login.trim()) {
    errors.login = 'Login is required'
  }
  if (!form.password.trim()) {
    errors.password = 'Password is required'
  } else if (form.password.length < 6) {
    errors.password = 'Password must be at least 6 characters'
  }
  if (!form.confirmPassword.trim()) {
    errors.confirmPassword = 'Please confirm your password'
  } else if (form.confirmPassword !== form.password) {
    errors.confirmPassword = 'Passwords do not match'
  }

  return !errors.login && !errors.password && !errors.confirmPassword
}

const onSubmit = async () => {
  if (!validate()) return
  try {
    await auth.register({
      login: form.login,
      password: form.password,
      confirmPassword: form.confirmPassword,
    })
    // After successful register, redirect to login
    router.push({ name: 'Login' })
  } catch {
    // error already set in auth.error
  }
}
</script>