import axios from 'axios'

const baseURL = import.meta.env.API_URL || 'http://localhost:8080/api/v1'
let accessToken = ''
let refreshRequest: Promise<void> | null = null

export const api = axios.create({
	baseURL,
	withCredentials: true,
})

api.interceptors.request.use((cfg) => {
	if (accessToken) cfg.headers.Authorization = `Bearer ${accessToken}`
	return cfg
})

api.interceptors.response.use(
	(response) => response,
	async (error) => {
		const request = error.config
		if (error.response?.status !== 401 || request.url === '/login' || request._retry) throw error

		request._retry = true
		await refreshAccessToken()
		request.headers.Authorization = `Bearer ${accessToken}`
		return api(request)
	},
)

export function setAccessToken(token: string) {
	accessToken = token
}

export function clearAccessToken() {
	accessToken = ''
}

export function hasAccessToken() {
	return accessToken !== ''
}

export async function refreshAccessToken() {
	if (!refreshRequest) {
		refreshRequest = axios
			.post(`${baseURL}/refresh`, undefined, { withCredentials: true })
			.then(({ data }) => setAccessToken(data.token))
			.finally(() => {
				refreshRequest = null
			})
	}
	return refreshRequest
}

export const errMsg = (e: any): string =>
	e?.response?.data?.error ?? 'Something went wrong'
