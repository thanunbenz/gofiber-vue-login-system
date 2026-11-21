<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { api, errMsg } from '../api'

const router = useRouter()
const username = ref('')
const password = ref('')
const confirmPassword = ref('')
const error = ref('')

async function submit() {
  error.value = ''
  if (password.value !== confirmPassword.value) {
    error.value = 'Passwords do not match'
    return
  }
  try {
    await api.post('/register', {
      username: username.value,
      password: password.value,
      confirmPassword: confirmPassword.value,
    })
    router.push('/login')
  } catch (e) {
    error.value = errMsg(e)
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-100 flex items-center justify-center p-4">
    <div class="w-full max-w-sm bg-white rounded-xl shadow overflow-hidden">
      <div class="bg-green-500 text-center py-2 font-medium">IT 02-2</div>
      <form @submit.prevent="submit" class="p-6 space-y-4">
        <input
          v-model="username"
          placeholder="User"
          required
          class="w-full border rounded px-3 py-2"
        />
        <input
          v-model="password"
          type="password"
          placeholder="Password"
          required
          class="w-full border rounded px-3 py-2"
        />
        <input
          v-model="confirmPassword"
          type="password"
          placeholder="Confirm Password"
          required
          class="w-full border rounded px-3 py-2"
        />
        <p v-if="error" class="text-red-600 text-sm">{{ error }}</p>
        <button class="w-full bg-green-500 hover:bg-green-600 text-white rounded py-2">
          สมัครสมาชิก
        </button>
      </form>
    </div>
  </div>
</template>
