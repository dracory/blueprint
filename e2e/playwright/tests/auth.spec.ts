import { test, expect, type Page } from '@playwright/test';

// Credentials of the users seeded by e2e/browser (seed.go).
const ADMIN_EMAIL = 'e2e-admin@blueprint.local';
const USER_EMAIL = 'e2e-user@blueprint.local';
const PASSWORD = 'password123';

async function login(page: Page, email: string) {
  await page.goto('/auth/login');
  await page.fill('input[name="email"]', email);
  await page.fill('input[name="password"]', PASSWORD);
  await page.locator('button.ButtonLogin, button[type="submit"]').first().click();
}

test.describe('authentication', () => {
  test('login page loads', async ({ page }) => {
    await page.goto('/auth/login');
    await expect(page.locator('input[name="email"]')).toBeVisible();
    await expect(page.locator('input[name="password"]')).toBeVisible();
  });

  test('admin lands on /admin after login', async ({ page }) => {
    await login(page, ADMIN_EMAIL);
    await expect(page).toHaveURL(/\/admin/, { timeout: 10_000 });
  });

  test('regular user lands on /user after login', async ({ page }) => {
    await login(page, USER_EMAIL);
    await expect(page).toHaveURL(/\/user/, { timeout: 10_000 });
  });
});

test.describe('authorization', () => {
  test('regular user cannot access /admin', async ({ page }) => {
    await login(page, USER_EMAIL);
    await expect(page).toHaveURL(/\/user/, { timeout: 10_000 });

    // AdminMiddleware redirects non-admins to "/" with a flash error.
    await page.goto('/admin');
    await expect(page).not.toHaveURL(/\/admin/, { timeout: 10_000 });
  });
});
