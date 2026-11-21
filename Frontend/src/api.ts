import axios from 'axios'

export const api = axios.create({
    baseURL: import.meta.env.BASE_URL ?? 'http://localhost:8080/api',
})

api.interceptors.request.use((cfg) => {
    const token = localStorage.getItem('token')
    if (token) cfg.headers.Authorization = `Bearer ${token}`
    return cfg
})

export const errMsg = (e: any): string =>
    e?.response?.data?.error ?? 'Something went wrong'
