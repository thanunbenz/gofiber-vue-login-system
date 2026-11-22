<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { api, errMsg, setAccessToken } from '../api'

const router = useRouter()
const username = ref('')
const password = ref('')
const error = ref('')

async function submit() {
  error.value = ''
  try {
    const { data } = await api.post('/login', {
      username: username.value,
      password: password.value,
    })
    setAccessToken(data.token)
    router.push('/welcome')
  } catch (e) {
    error.value = errMsg(e)
  }
}
</script>

<template>
  <div class="min-h-screen bg-gray-100 flex items-center justify-center p-4">
    <div class="w-full max-w-sm bg-white rounded-xl shadow overflow-hidden">
      <div class="bg-green-500 text-center py-2 font-medium">IT 02-1</div>
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
        <p v-if="error" class="text-red-600 text-sm">{{ error }}</p>
        <button class="w-full bg-green-500 hover:bg-green-600 text-white rounded py-2">
          ลงชื่อเข้าใช้งาน
        </button>
        <router-link to="/register" class="block text-center text-blue-600 underline text-sm">
          สมัครสมาชิก
        </router-link>
      </form>
    </div>
  </div>
</template>
