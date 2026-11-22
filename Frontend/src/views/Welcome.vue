<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { api, clearAccessToken } from '../api'

const router = useRouter()
const username = ref('')

async function logout() {
  await api.post('/logout').catch(() => undefined)
  clearAccessToken()
  router.push('/login')
}

onMounted(async () => {
  try {
    const { data } = await api.get('/me')
    username.value = data.username
  } catch {
    await logout()
  }
})
</script>

<template>
  <div class="min-h-screen bg-gray-100 flex items-center justify-center p-4">
    <div class="w-full max-w-lg bg-white rounded-xl shadow overflow-hidden">
      <div class="bg-green-500 text-center py-2 font-medium">IT 02-3</div>
      <div class="p-10 text-center space-y-6">
        <p class="text-xl">Welcome User : {{ username }}</p>
        <button @click="logout" class="text-sm text-gray-500 underline">ออกจากระบบ</button>
      </div>
    </div>
  </div>
</template>
