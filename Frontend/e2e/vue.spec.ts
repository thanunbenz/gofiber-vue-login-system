import { test, expect } from '@playwright/test'

test('logs in successfully', async ({ page }) => {
  await page.route('**/api/v1/login', async (route) => {
    await route.fulfill({ json: { token: 'access-token' } })
  })
  await page.route('**/api/v1/me', async (route) => {
    await route.fulfill({ json: { username: 'thanunbenz' } })
  })

  await page.goto('/login')
  await page.getByPlaceholder('User').fill('thanunbenz')
  await page.getByPlaceholder('Password').fill('password')
  await page.getByRole('button').click()

  await expect(page).toHaveURL(/\/welcome$/)
  await expect(page.getByText('Welcome User : thanunbenz')).toBeVisible()
})

test('shows an error when login fails', async ({ page }) => {
  await page.route('**/api/v1/login', async (route) => {
    await route.fulfill({ status: 401, json: { error: 'invalid username or password' } })
  })

  await page.goto('/login')
  await page.getByPlaceholder('User').fill('thanunbenz')
  await page.getByPlaceholder('Password').fill('wrong-password')
  await page.getByRole('button').click()

  await expect(page.getByText('invalid username or password')).toBeVisible()
})
